
学习[gRPC](http://www.grpc.io/)时写的示例程序。主要代码在src/grpcdemo目录中。

程序实现的远程服务很简单，参考gRPC给出的示例代码写成，主要在于尝试gRPC的用法。

* 不使用TLS时，客户端连接时需使用grpc.WithInsecure()选项。
* 客户端执行远程方法调用时，可使用grpc.FailFast(false)表示必要时等待有可用的服务器完成调用。
* 服务器端调用SetHeader和SetTrailer来设置头部和尾部元数据才有意义,客户端调用这两个方法没有意义。
* balancer.go实现了简单负载均衡。负载均衡器工作在客户端，在执行远程方法调用时，负载均衡器从多个远程服务器中选择一个来发起远程调用。这与nginx、haproxy之类的工作在服务器端的负载均衡器是不一样的。
