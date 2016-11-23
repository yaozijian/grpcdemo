package main

import (
	"fmt"
	"io"
	"math/rand"
	"net"
	"os"
	"strings"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type (
	mathServer struct {
		header  metadata.MD
		trailer metadata.MD
	}
)

func main() {

	srv := &mathServer{}

	if nic, _ := net.InterfaceByName(os.Args[1]); nic != nil {
		if addrs, _ := nic.Addrs(); len(addrs) > 0 {
			addr := addrs[0].String()
			if idx := strings.Index(addr, "/"); idx > 0 {
				addr = addr[:idx]
			}
			srv.header = metadata.Pairs("serveraddr", addr+":9999")
			srv.trailer = metadata.Pairs("End", "End")
			fmt.Printf("在地址%v上监听\n", addr+":9999")
		}
	}

	server := grpc.NewServer()
	RegisterMathServer(server, srv)
	listener, _ := net.Listen("tcp4", ":9999")
	server.Serve(listener)
}

func (server *mathServer) setHeaderTrailer(ctx context.Context, stream grpc.ServerStream) {
	if len(server.header) > 0 {
		if ctx != nil {
			grpc.SetHeader(ctx, server.header)
			grpc.SetTrailer(ctx, server.trailer)
		} else if stream != nil {
			stream.SetHeader(server.header)
			stream.SetTrailer(server.trailer)
		}
	}
}

func (server *mathServer) Add(ctx context.Context, op *OpRequest) (reply *OpReply, err error) {

	server.setHeaderTrailer(ctx, nil)

	reply = &OpReply{
		A: op.A,
		B: op.B,
		R: op.A + op.B,
	}
	return
}

func (server *mathServer) Addlist(stream Math_AddlistServer) error {

	result := &Operator{}

	server.setHeaderTrailer(nil, stream)

	for {
		op, err := stream.Recv()
		if err == nil {
			result.A += op.A
		} else if err == io.EOF {
			return stream.SendAndClose(result)
		} else {
			return err
		}
	}
}

func (server *mathServer) Random(op *Operator, stream Math_RandomServer) error {

	server.setHeaderTrailer(nil, stream)

	cnt := int32(0)
	for op.A <= 0 || cnt < op.A {
		cnt++
		if err := stream.Send(&Operator{A: int32(rand.Int())}); err != nil {
			return err
		}
	}
	return nil
}

func (server *mathServer) Addstream(stream Math_AddstreamServer) error {

	server.setHeaderTrailer(nil, stream)

	for {
		if op, err := stream.Recv(); err == nil {
			res := &OpReply{
				A: op.A,
				B: op.B,
				R: op.A + op.B,
			}
			if err = stream.Send(res); err != nil {
				return err
			}
		} else if err == io.EOF {
			return nil
		} else {
			return err
		}
	}
}
