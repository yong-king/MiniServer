package main

import (
	"context"
	"strconv"
	"time"

	"com.ysh.blog.booksotre/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type server struct{
	pb.UnimplementedBookstoreServer

	bs *bookstore
}
const (
	defaultNextId = "0"
	defaultPageSize = 2
)

// ListShelves 获取书记列表
func (s *server) ListShelves(ctx context.Context, in *emptypb.Empty) (*pb.ListShelvesResponse, error) {
	sl, err := s.bs.ListShelves(ctx)
	if err == gorm.ErrEmptySlice {
		// 没有数据
		return &pb.ListShelvesResponse{}, nil
	}
	if err != nil {
		return nil, status.Error(codes.Internal, "query failed")
	}
	// 返回封装
	nsl := make([]*pb.Shelf, 0, len(sl))
	for _, s := range sl{
		nsl = append(nsl, &pb.Shelf{
			Id: s.ID,
			Theme: s.Theme,
			Size: s.Size,
		})
	}
	return &pb.ListShelvesResponse{Shelves: nsl}, nil
}

// CreateShelf 创建书架
func (s *server) CreateShelf(ctx context.Context, in *pb.CreateShelfRequest) (*pb.Shelf, error) {
	// 检查参数
	if len(in.GetShelf().GetTheme()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid theme")
	}

	// 数据准备
	data := Shelf{Theme: in.GetShelf().GetTheme(), Size: in.GetShelf().GetId()}
	shelf, err := s.bs.CreateShelf(ctx, data)
	if err != nil {
		return nil, status.Error(codes.Internal, "create failed")
	}
	return &pb.Shelf{Id: shelf.ID, Theme: shelf.Theme, Size: shelf.Size}, nil
}

// GetShelf 根据id获取书架
func (s *server) GetShelf(ctx context.Context, in *pb.GetShelfRequest) (*pb.Shelf, error) {
	// 检查参数
	if in.GetShelf() <= 0{
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}
	// 查询数据
	sl, err := s.bs.GetShelf(ctx, in.GetShelf())
	if err != nil {
		return nil, status.Error(codes.Internal, "query failed")
	}
	// 封装返回
	return &pb.Shelf{Id: sl.ID, Theme: sl.Theme, Size: sl.Size}, nil
}

// DeleteShelf 删除书架
func (s *server)  DeleteShelf(ctx context.Context, in *pb.DeleteShelfRequest) (*emptypb.Empty, error){
	// 检查参数
	if in.GetShelf() <= 0{
		return nil, status.Error(codes.InvalidArgument, "invalid id")
	}
	// 删除书架
	err := s.bs.DelateShelf(ctx, in.GetShelf())
	if err != nil {
		return nil, status.Error(codes.Internal, "delate falied")
	}
	// 返回
	return &emptypb.Empty{}, nil
}

// CreateBook 创建书籍
func (s *server)CreateBook(ctx context.Context, in *pb.CreateBookRequest) (*pb.Book, error) {
	// 检查参数
	if in.GetShelf() <= 0{
		return nil, status.Error(codes.InvalidArgument, "invalid shelf id")
	}
	if len(in.GetBook().GetAuthor()) <= 0{
		return nil, status.Error(codes.InvalidArgument, "invalid author")
	}
	if len(in.GetBook().GetTitle()) <= 0{
		return nil, status.Error(codes.InvalidArgument, "invalid title")
	}
	// 处理书籍
	data := Book{Title: in.GetBook().GetTitle(), Author: in.GetBook().GetAuthor(), ShelfID: in.GetShelf()}
	// 添加书籍
	b, err := s.bs.CreateBook(ctx, data)
	if err != nil {
		return nil, status.Error(codes.Internal, "create failed")
	}
	return &pb.Book{Id: b.ID, Title: b.Title, Author: b.Author}, nil
}

func Invalid(p Page) bool{
	return p.NextID == ""|| p.NextTimeAtUTC == 0 || p.NextTimeAtUTC > time.Now().Unix() || p.PageSize <= 0
}

// ListBooks 根据指定书架获取书籍
func(s *server) ListBooks(ctx context.Context, in *pb.ListBooksRequest) (*pb.ListBooksResponse, error) {
	// 检查参数
	if in.GetShelf() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid shelf id")
	}
	var(
		cursor = defaultNextId
		pageSize = defaultPageSize	
	)
	if len(in.GetPageToken()) > 0 {
		// 解码
		p := Token(in.PageToken).Decode()
		// 参数校验
		if Invalid(p){
			return nil, status.Error(codes.InvalidArgument, "invalid page token") 
		}
		cursor = p.NextID
		pageSize = int(p.PageSize)
	}
	// 查询书籍列表
	bl, err := s.bs.ListBooks(ctx, in.GetShelf(), cursor, pageSize+1)
	if err == gorm.ErrEmptySlice{
		return &pb.ListBooksResponse{}, nil
	}
	if err != nil {
		return nil, status.Error(codes.Internal, "query failed")
	}
	var (
		hasNextPageToken bool
		realSize int = len(bl)
		nextPageToken string
	)
	if len(bl) > pageSize {
		hasNextPageToken = true
		realSize = pageSize
	}
	// 封装返回
	nbl := make([]*pb.Book, 0, len(bl))
	for i := range(realSize){
			nbl = append(nbl, &pb.Book{
			Id: bl[i].ID,
			Author: bl[i].Author,
			Title: bl[i].Title,
		})
	}
	// 如果有下一页
	if hasNextPageToken{
		nextPage := Page{
			NextID: strconv.FormatInt(int64(nbl[realSize-1].Id), 10),
			NextTimeAtUTC: time.Now().Unix(),
			PageSize: int64(realSize),
		}
		nextPageToken = string(nextPage.Encode())
	}
	return &pb.ListBooksResponse{Books: nbl, NextPageToken: nextPageToken}, nil
}

// GetBook 获取指定书籍
func (s *server) GetBook(ctx context.Context, in *pb.GetBookRequest) (*pb.GetBookRequesResponse, error) {
	// 检查参数
	if in.GetShelf() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid shelf id")
	}
	if in.GetBook() <= 0{
		return nil, status.Error(codes.InvalidArgument, "invalid book id")
	}
	// 查询书籍
	b, err := s.bs.GetBook(ctx, in.GetBook(), in.GetShelf())
	if err == gorm.ErrEmptySlice {
		return &pb.GetBookRequesResponse{}, nil
	}
	if err != nil{
		return nil, status.Error(codes.Internal, "query failed")
	}
	// 封装返回
	book := pb.Book{Id: b.ID, Title: b.Title, Author: b.Author}
	return &pb.GetBookRequesResponse{Book: &book}, nil
}

// DeleteBook 删除书籍
func (s *server)DeleteBook(ctx context.Context, in *pb.DeleteBookRequest) (*emptypb.Empty, error){
	// 查询参数
	if in.GetShelf() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid shelf id")
	}
	if in.GetBook() <= 0{
		return nil, status.Error(codes.InvalidArgument, "invalid book id")
	}
	// 删除书籍
	err := s.bs.DelateBook(ctx, in.GetBook(), in.GetShelf())
	if err != nil {
		return nil, status.Error(codes.Internal, "delete book failed")
	}
	return &emptypb.Empty{}, nil
}