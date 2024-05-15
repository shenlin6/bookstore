package main

import (
	"bookstore/pb"
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
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

	//同一个端口分别处理gRPC和HTTP(重要)
	// 1. 创建gRPC-Gateway mux
	gwmux := runtime.NewServeMux()
	dops := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := pb.RegisterBookstoreHandlerFromEndpoint(context.Background(),
		gwmux, "127.0.0.1:8080", dops); err != nil {
		log.Fatalf("RegisterBookstoreHandlerFromEndpoint failed,err:%v\n", err)
		return
	}

	// 2. 创建HTTP mux
	mux := http.NewServeMux()
	//注册路由
	mux.Handle("/", gwmux) // "/" 表示根路径，即匹配所有的 HTTP 请求。

	// 3. 定义HTTP server
	gwServer := &http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: grpcHandlerFunc(s, mux), //分别传入 grpc和HTTP的handler
	}
	// 4. 启动
	log.Println("Serving on http://127.0.0.1:8080")
	gwServer.Serve(l)

}

// grpcHandlerFunc :自定义函数 将 gRPC 请求和 HTTP 请求分别调用不同的 handler 处理 （可直接拿来主义）
func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	}), &http2.Server{})
}
