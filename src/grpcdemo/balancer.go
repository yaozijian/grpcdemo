package main

import (
	"fmt"
	"sync"
	"time"

	"container/heap"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type (
	roundRobinWithCookie struct {
		available serverList // 在线服务器
		offlines  serverList // 离线服务器

		addrch  chan []grpc.Address // 对外发布服务器地址列表
		newsrv  chan int
		waitsrv bool

		// 用于退出处理
		stop   context.Context
		cancel func()
		wait   sync.WaitGroup
		sync.Mutex
	}

	serverItem struct {
		addr  grpc.Address  // 服务器地址
		total time.Duration // 调用总计耗时
		avg   time.Duration // 单次调用平均耗时
		cnt   int64         // 调用次数
		calls int64         // 未决调用数
		sync.Mutex
	}

	serverList []*serverItem
)

const (
	key_server_addr = "serveraddr"
)

func newBalancer(dstlist ...string) grpc.Balancer {

	stop, cancel := context.WithCancel(context.Background())

	rbwc := &roundRobinWithCookie{
		available: make(serverList, 0, 1000),
		offlines:  make(serverList, 0, 1000),
		addrch:    make(chan []grpc.Address, 1),
		newsrv:    make(chan int, 1),
		stop:      stop,
		cancel:    cancel,
	}

	for _, server := range dstlist {
		rbwc.addServer(grpc.Address{Addr: server}, &rbwc.offlines)
	}

	rbwc.wait.Add(1)
	go rbwc.sendAddrsNotify()

	return rbwc
}

func (rbwc *roundRobinWithCookie) Start(target string, config grpc.BalancerConfig) error {
	fmt.Printf("启动负载均衡器: target=%v\n", target)
	return nil
}

// 加入已经建立连接的地址到可用服务器列表中
func (rbwc *roundRobinWithCookie) Up(addr grpc.Address) (down func(error)) {

	fmt.Printf("已经建立到服务器的连接: %v\n", addr.Addr)

	rbwc.Lock()
	rbwc.delServer(addr, &rbwc.offlines)
	rbwc.addServer(addr, &rbwc.available)
	if rbwc.waitsrv {
		select {
		case rbwc.newsrv <- 1:
		default:
		}
	}
	rbwc.Unlock()

	down = func(err error) { rbwc.onDown(addr, err) }

	return
}

// 从可用服务器列表中删除已经连接断开的
func (rbwc *roundRobinWithCookie) onDown(addr grpc.Address, err error) {
	rbwc.Lock()
	rbwc.delServer(addr, &rbwc.available)
	rbwc.addServer(addr, &rbwc.offlines)
	rbwc.Unlock()
	fmt.Printf("到服务器的连接已经断开: addr=%v err=%v\n", addr.Addr, err)
}

func (rbwc *roundRobinWithCookie) Get(ctx context.Context, opts grpc.BalancerGetOptions) (addr grpc.Address, put func(), err error) {

	rbwc.Lock()
	defer rbwc.Unlock()

begin_get_server:

	// 默认选择第一个服务器,即接受调用次数最少的那个服务器
	selected_idx := 0
	err = nil

	if len(rbwc.available) == 0 {
		err = fmt.Errorf("没有可用的服务器")
	} else if md, ok := metadata.FromContext(ctx); ok && len(md) > 0 {
		// Cookie 记录的服务器仍然可用
		if x, _ := md[key_server_addr]; len(x) > 0 {
			for idx, item := range rbwc.available {
				if item.addr.Addr == x[0] {
					selected_idx = idx
					break
				}
			}
		}
	}

	if err == nil {

		item := rbwc.available[selected_idx]
		item.cnt++
		item.calls++
		heap.Fix(&rbwc.available, selected_idx)

		addr = item.addr
		addrx := item.addr
		begin := time.Now()
		put = func() { rbwc.onCallDone(addrx, begin) }
	} else if opts.BlockingWait {

		var timeout <-chan time.Time
		var timer *time.Timer

		// 有最后期限吗?
		if t, ok := ctx.Deadline(); ok {
			if dura := time.Now().Sub(t); int64(dura) > int64(time.Millisecond) {
				timer = time.NewTimer(dura)
				timeout = timer.C
			} else {
				err = fmt.Errorf("没有可用的服务器: 已经到达最后期限")
				return
			}
		}

		rbwc.waitsrv = true

		rbwc.Unlock()

		fmt.Println("开始等待有可用的服务器")

		// 等待直到超时或者有新的服务器上线
		select {
		case <-timeout:
			timer.Stop()
			err = fmt.Errorf("没有可用的服务器: 等待超时")
			rbwc.Lock()
			rbwc.waitsrv = false
		case <-rbwc.newsrv:
			if timer != nil {
				timer.Stop()
			}
			rbwc.Lock()
			rbwc.waitsrv = false
			goto begin_get_server
		}
	} else {
		err = fmt.Errorf("没有可用的服务器")
	}

	return
}

func (rbwc *roundRobinWithCookie) Notify() <-chan []grpc.Address {
	return rbwc.addrch
}

func (rbwc *roundRobinWithCookie) Close() error {
	rbwc.cancel()
	rbwc.wait.Wait()
	return nil
}

// 添加一个地址到服务器列表中
func (rbwc *roundRobinWithCookie) addServer(addr grpc.Address, plist *serverList) {

	addit := true

	for _, item := range *plist {
		if item.addr == addr {
			addit = false
			break
		}
	}

	if addit {
		item := &serverItem{addr: addr}
		heap.Push(plist, item)
	}
}

// 从服务器列表中删除一个地址
func (rbwc *roundRobinWithCookie) delServer(addr grpc.Address, from *serverList) {
	for idx, item := range *from {
		if item.addr == addr {
			heap.Remove(from, idx)
			break
		}
	}
}

// RPC调用完成,统计调用时间
func (rbwc *roundRobinWithCookie) onCallDone(addr grpc.Address, begin time.Time) {

	rbwc.Lock()
	defer rbwc.Unlock()

	for _, item := range rbwc.available {
		if item.addr == addr {
			item.calls--
			item.total += time.Now().Sub(begin)
			item.avg = time.Duration(int64(item.total) / item.cnt)
			//fmt.Printf("调用完成: addr=%v cnt=%v avg=%v\n", item.addr, item.cnt, item.avg.String())
			break
		}
	}
}

// 向外部发送可用的服务器地址列表
func (rbwc *roundRobinWithCookie) sendAddrsNotify() {

	defer rbwc.wait.Done()

	for {
		select {
		case rbwc.addrch <- rbwc.addressList():
		case <-rbwc.stop.Done():
			return
		}
	}
}

func (rbwc *roundRobinWithCookie) addressList() []grpc.Address {

	rbwc.Lock()
	defer rbwc.Unlock()

	addrs := make([]grpc.Address, len(rbwc.offlines)+len(rbwc.available))

	for idx := range rbwc.offlines {
		addrs[idx] = rbwc.offlines[idx].addr
	}

	off := len(rbwc.offlines)

	for idx := range rbwc.available {
		addrs[idx+off] = rbwc.available[idx].addr
	}

	return addrs
}

func (obj *serverList) Len() int           { return len(*obj) }
func (obj *serverList) Less(x, y int) bool { return (*obj)[x].cnt < (*obj)[y].cnt }
func (obj *serverList) Swap(x, y int)      { (*obj)[x], (*obj)[y] = (*obj)[y], (*obj)[x] }
func (obj *serverList) Push(v interface{}) { *obj = append(*obj, v.(*serverItem)) }
func (obj *serverList) Pop() interface{} {

	tmplist := *obj

	cnt := len(tmplist)
	last := tmplist[cnt-1]
	tmplist[cnt-1] = nil
	tmplist = tmplist[:cnt-1]

	*obj = tmplist

	return last
}
