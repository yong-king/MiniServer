package main

import (
	"context"
	"testing"

	"com.ysh.blog.booksotre/pb"
)

func Test_server_ListBooks(t *testing.T) {
	db, _ := NewDB()
	s := server{bs: &bookstore{db: db}}

	bl, err := s.ListBooks(context.Background(), &pb.ListBooksRequest{
		Shelf: 3,
	})
	if err != nil{
		t.Fatalf("s.ListBooks failed:%v\n", err)
	}
	t.Logf("next page tokenn:%v\n", bl.GetNextPageToken())

	for i, book := range(bl.GetBooks()){
		t.Logf("%d \n%v", i, book)
	}
}