package main

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type demoClient struct {
	client MathClient
	ctx    context.Context
}

func main() {

	balancer := newBalancer(os.Args[1:]...)
	defer balancer.Close()

	opts := []grpc.DialOption{grpc.WithInsecure(), grpc.WithBalancer(balancer)}
	//opts = append(opts, grpc.WithBlock()) // 阻塞直到连接完成?

	conn, _ := grpc.Dial("target", opts...)
	defer conn.Close()

	rand.Seed(time.Now().UnixNano())

	client := &demoClient{
		client: NewMathClient(conn),
		ctx:    context.Background(),
	}

	client.demo1()
	client.demo2()
	client.demo3()
	client.demo4()
}

func (obj *demoClient) demo1() {

	var header metadata.MD

	op := &OpRequest{A: int32(rand.Intn(100)), B: int32(rand.Intn(100))}

	// grpc.FailFast(false) 表示等待有可用的服务器
	reply, e := obj.client.Add(obj.ctx, op, grpc.Header(&header), grpc.FailFast(false))

	if e == nil {
		if len(header) > 0 {
			fmt.Printf("demo 1 From %v: %d + %d = %v\n\n", header["serveraddr"], op.A, op.B, reply.R)
			obj.ctx = metadata.NewContext(obj.ctx, header)
		} else {
			fmt.Printf("demo 1 From **Unknown**: %d + %d = %v\n\n", op.A, op.B, reply.R)
		}
	} else {
		fmt.Printf("demo1 error: %v\n", e)
	}
}

func (obj *demoClient) demo2() {

	stream, e := obj.client.Addlist(obj.ctx)

	idx := 0
	for e == nil && idx < 5 {
		a := &Operator{A: int32(rand.Intn(100))}
		fmt.Printf("%v + ", a.A)
		e = stream.Send(a)
		idx++
	}

	if e == nil {
		r, _ := stream.CloseAndRecv()
		fmt.Printf(" = %v\n\n", r.A)
	}

	if e == nil {
		if h, e := stream.Header(); h != nil {
			fmt.Printf("demo2 From %v: \n\n", h["serveraddr"])
		} else {
			fmt.Printf("获取流header失败: %v\n", e)
		}
	} else {
		fmt.Printf("demo2: e=%v\n\n", e)
	}
}

func (obj *demoClient) demo3() {

	stream, e := obj.client.Random(obj.ctx, &Operator{A: 5})

	if e == nil {
		if h, e := stream.Header(); h != nil {
			fmt.Printf("demo3: From %v: \n", h["serveraddr"])
		} else {
			fmt.Printf("获取流header失败: %v\n", e)
		}
	}

	for e == nil {
		x, e := stream.Recv()
		if e == nil {
			fmt.Printf("random: %v\n", x.A)
		} else if e == io.EOF {
			fmt.Printf("random: End\n\n")
			break
		} else {
			fmt.Printf("random: %v\n", e)
			break
		}
	}

	fmt.Printf("demo3: e=%v\n\n", e)
}

func (obj *demoClient) demo4() {

	stream, e := obj.client.Addstream(obj.ctx)

	var r *OpReply

	idx := 0
	for e == nil && idx < 4 {
		op := &OpRequest{A: int32(rand.Intn(100)), B: int32(rand.Intn(100))}
		if e = stream.Send(op); e == nil {
			if r, e = stream.Recv(); e == nil {
				fmt.Printf("%v + %v = %v\n", r.A, r.B, r.R)
				idx++
			} else {
				fmt.Printf("%v + %v = %v\n", op.A, op.B, e)
			}
		} else {
			fmt.Printf("%v + %v = %v\n", op.A, op.B, e)
		}
	}

	if e == nil {
		if h, e := stream.Header(); h != nil {
			fmt.Printf("demo 4 From %v: \n", h["serveraddr"])
		} else {
			fmt.Printf("获取流header失败: %v\n", e)
		}
	}

	fmt.Println("demo4", e)
}
