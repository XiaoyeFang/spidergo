package main

import (
	"flag"
	"net"

	"crawler-apkmirror/cmd"
	"crawler-apkmirror/collect"
	"crawler-apkmirror/config"
	pb "crawler-apkmirror/protos"

	_ "net/http/pprof"

	"github.com/golang/glog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	flag.Parse()
	cmd.Execute()
	grpcStart()
	glog.Flush()
}

func grpcStart() {
	lis, err := net.Listen("tcp", config.AppConf.GrpcListen)
	if err != nil {
		panic(err)
	}
	//glog.Errorf("am-grpc-port %v", config.AppConf.GrpcListen)
	server := grpc.NewServer()
	srv := collect.CollectServer{}
	pb.RegisterCrawlerAMServiceServer(server, &srv)
	reflection.Register(server)
	if err := server.Serve(lis); err != nil {
		panic(err)
	}
}