package main

import (
	"fmt"

	// "google.golang.org/protobuf/types/known/wrapperspb"
	// "github.com/gohugoio/hugo/common/paths"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"oneof.demo/api"
"github.com/iancoleman/strcase"
fieldmask_utils "github.com/mennanov/fieldmask-utils"
)

func NoticeHandeler() {

	// clinent
	// req := api.NoticeReaderRequest{
	// 	Msg: "yk又来消息啦",
	// 	NoticeWay: &api.NoticeReaderRequest_Email{
	// 		Email: "xxxxx@mail.com",
	// 	},
	// }

	req := api.NoticeReaderRequest{
		Msg: "yk又来消息啦",
		NoticeWay: &api.NoticeReaderRequest_Phone{
			Phone: "156xxxx",
		},
	}

	// server
	switch v := req.NoticeWay.(type) {
	case *api.NoticeReaderRequest_Email:
		noticeWithEmail(v)
	case *api.NoticeReaderRequest_Phone:
		noticeWithPhone(v)
	}

}

// 发送通知相关的功能函数
func noticeWithEmail(in *api.NoticeReaderRequest_Email) {
	fmt.Printf("notice reader by email:%v\n", in.Email)
}

func noticeWithPhone(in *api.NoticeReaderRequest_Phone) {
	fmt.Printf("notice reader by phone:%v\n", in.Phone)
}

func BookHandelr() {
	// client
	book := api.Book{
		Titile: "《中国》",
		Author: "yk",
		// Price: &wrapperspb.Int64Value{Value: 999},
		Price: proto.Int64(999),
	}

	// server
	if book.Price == nil{
		// 没有赋值
		fmt.Println("book have no price")
	} else {
		fmt.Printf("book with price: %v\n", book.GetPrice())
	}
}

func UpdataBook() {
	// client
	paths := []string{"titile", "price", "info.b"}
	update := api.UpdateBookRequest{
		Opt: "yk",
		Book: &api.Book{
			Titile: "<666>",
			Price: proto.Int64(888),
			Info: &api.Book_Info{
				B: 666,
			},
		},
		UpdateMask: &fieldmaskpb.FieldMask{Paths: paths},
	}

	// server
	mask, _ :=fieldmask_utils.MaskFromProtoFieldMask(update.UpdateMask, strcase.ToCamel)
	var bookDst = make(map[string]interface{})
	fieldmask_utils.StructToMap(mask, update.Book, bookDst)
	fmt.Printf("bookDst:%v\n", bookDst)
}

func main() {
	// NoticeHandeler()
	// BookHandelr()
	UpdataBook()
}
