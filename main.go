package main

import (
	"bookstore/pb"
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

//Bookstore

func main() {
	// 连接数据库
	db, err := NewDB("bookstore.db")
	if err != nil {
		fmt.Printf("connect to db failed,err:%v\n", err)
		return
	}
	// 创建server
	srv := server{
		bs: &bookstore{db: db},
	}

	//启动grpc服务
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Printf("falied to listen,err:%v\n", err)
		return
	}
	//创建grpc服务
	s := grpc.NewServer()
	//注册服务
	pb.RegisterBookstoreServer(s, &srv)
	//开启goroutine单独启动grpc server
	go func() {
		log.Println(s.Serve(l))
	}()

	//grpc Gateway
	conn, err := grpc.DialContext(
		context.Background(),
		"127.0.0.1:8080",
		grpc.WithBlock(), //阻塞，等待连接成功
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("failed to dial server ", err)
	}
	//创建HTTP server(核心)
	gwmux := runtime.NewServeMux()
	//注册服务
	err = pb.RegisterBookstoreHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalln("failed to register gateway ", err)
	}
	//创建HTTP接口
	gwServer := &http.Server{  //使用指针类型避免不必要的复制
		Addr:    "127.0.0.1:8080",
		Handler: gwmux,
	}
	//启动HTTP服务
	log.Println("Serving grpc-Gateway on http://127.0.0.1:8080")
	log.Fatalln(gwServer.ListenAndServe())
}
