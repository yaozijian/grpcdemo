#!/bin/bash

cur=$(pwd)

vendor=$(pwd)/vendor

cd src/grpcdemo
rm -f vendor
ln -s $vendor vendor

GOPATH=$cur go build client.go balancer.go math.pb.go
GOPATH=$cur go build server.go math.pb.go

cd $cur
