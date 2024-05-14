package main

// bookstore grpc服务 实现增删改查四个方法

import (
	"bookstore/pb"
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type server struct {
	pb.UnimplementedBookstoreServer

	bs *bookstore // data.go里面的
}

// ListShelves 获取书架列表
func (s *server) ListShelves(ctx context.Context, in *emptypb.Empty) (*pb.ListShelvesResponse, error) {
	// 调用orm操作的方法
	sl, err := s.bs.ListShelves(ctx)
	//如果切片为空
	if err == gorm.ErrEmptySlice {
		return &pb.ListShelvesResponse{}, err
	}
	//查询数据库失败
	if err != nil {
		return nil, status.Error(codes.Internal, "query failed")
	}
	// 封装返回数据
	nsl := make([]*pb.Shelf, 0, len(sl)) //创建新切片需要使用make
	for _, s := range sl {
		nsl = append(nsl, &pb.Shelf{
			Id:    s.ID,
			Theme: s.Theme,
			Size:  s.Size,
		})
	}
	return &pb.ListShelvesResponse{Shelves: nsl}, nil
}

// CreateShelf 创建书架
func (s *server) CreateShelf(ctx context.Context, in *pb.CreateShelfRequest) (*pb.Shelf, error) {
	//参数校验（勿忘）
	if len(in.GetShelf().GetTheme()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid theme")
	}
	//更新数据
	data := Shelf{
		Theme: in.GetShelf().GetTheme(),
		Size:  in.GetShelf().GetSize(),
	}
	// 去数据库创建
	ns, err := s.bs.CreateShelf(ctx, data)
	if err != nil {
		return nil, status.Error(codes.Internal, "create failed")
	}
	return &pb.Shelf{Id: ns.ID, Theme: ns.Theme, Size: ns.Size}, nil
}

// GetShelf 获取书架
func (s *server) GetShelf(ctx context.Context, in *pb.GetShelfRequest) (*pb.Shelf, error) {
	//调用orm方法
	data, err := s.bs.GetShelf(ctx, in.GetShelf())
	if err != nil {
		return nil, status.Error(codes.Internal, "query failed")
	}
	//封装返回数据
	return &pb.Shelf{Id: data.ID, Theme: data.Theme, Size: data.Size}, nil
}

// DeleteShelf 根据id删除书架
func (s *server) DeleteShelf(ctx context.Context, in *pb.DeleteShelfRequest) (*emptypb.Empty, error) {
	//参数check
	if in.GetShelf() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid shelf id")
	}
	//调用orm方法
	err := s.bs.DeleteShelf(ctx, in.GetShelf())
	if err != nil { //数据库删除失败
		return nil, status.Error(codes.Internal, "delete failed")
	}
	//删除成功
	fmt.Printf("删除成功,删除了id为:%d的书架", in.GetShelf())
	return &emptypb.Empty{}, nil
}

// ListBooks 获取某个书架上的书籍列表
func (s *server) ListBooks(ctx context.Context, in *pb.ListBooksRequest) (*pb.ListBooksResponse, error) {
	// 调用orm操作的方法
	sl, err := s.bs.ListBooks(ctx, in.Shelf)
	//如果切片为空
	if err == gorm.ErrEmptySlice {
		return &pb.ListBooksResponse{}, err
	}
	//如果查询数据库失败
	if err != nil {
		return nil, status.Error(codes.Internal, "query failed")
	}
	//封装返回数据
	nsl := make([]*pb.Book, 0, len(sl))
	for _, s := range sl {
		nsl = append(nsl, &pb.Book{
			Id:     s.ID,
			Author: s.Author,
			Title:  s.Titile,
		})
	}
	return &pb.ListBooksResponse{Books: nsl}, nil
}

// CreateBook 创建书架上的书籍
func (s *server) CreateBook(ctx context.Context, in *pb.CreateBookRequest) (*pb.Book, error) {
	//参数校验
	if len(in.GetBook().Title) <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid title")
	}
	//更新数据
	data := Book{
		Author: in.GetBook().GetAuthor(),
		Titile: in.GetBook().GetTitle(),
		ShelfID: in.GetShelf(),
	}
	//到数据库里操作
	nb, err := s.bs.CreateBook(ctx, data)
	if err != nil {
		return nil, status.Error(codes.Internal, "create failed")
	}
	//返回更新操作
	return &pb.Book{Id: nb.ID,Author: nb.Author, Title: nb.Titile}, nil
}
