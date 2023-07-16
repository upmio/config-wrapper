package main

import (
	"context"
	"fmt"
	"github.com/upmio/config-wrapper/app/config"
	"log"

	"google.golang.org/grpc"
)

var (
	url       = "192.168.26.21:22105"
	namespace = "default"
	configmap = "test"
)

func main() {
	// grpc.Dial负责和gRPC服务建立链接
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// NewHelloServiceClient函数基于已经建立的链接构造HelloServiceClient对象,
	// 返回的client其实是一个HelloServiceClient接口对象
	client := config.NewSyncConfigServiceClient(conn)

	// 通过接口定义的方法就可以调用服务端对应的gRPC服务提供的方法
	req := &config.SyncConfigRequest{
		Namespace:     namespace,
		ConfigmapName: configmap,
	}
	resp, err := client.SyncConfig(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp.GetMessage())
}