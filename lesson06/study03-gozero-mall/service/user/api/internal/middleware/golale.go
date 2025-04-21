package middleware

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
)

// 功能
// 记录所有请求的响应信息

// 思路
// 将请求的响应信息写到我们的指定的文件当中
// 然后再往HTTP响应

// rest.Middleware --> Middleware func(next http.HandlerFunc) http.HandlerFunc
// type HandlerFunc func(ResponseWriter, *Request)

type bodyCopy struct{
	http.ResponseWriter	// 结构体嵌入接口类型，默认实现了接口的所有方法
	body *bytes.Buffer // 我们记录响应体的内容
}

func NewbodyCopy(w http.ResponseWriter) *bodyCopy {
	return &bodyCopy{
		ResponseWriter: w,
		body: bytes.NewBuffer([]byte{}),
	}
}

// 这里相当于重写了 ResponseWriter.Write
func (bc bodyCopy) Write(b []byte) (int, error) {
	// 1. 先记录到我们的这里
	bc.body.Write(b)
	// 2. 再往HTTP响应体写内容
	return bc.ResponseWriter.Write(b)
}

func CopyResq(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 初始化一个自定义的 ResponseWriter.Write
		bc := NewbodyCopy(w)
		// 实际执行完后，会执行bc.body.Write(b)， 然后再往HTTP响应体写内容
		next(bc, r)
		// 处理后的请求
		logx.Infof("-->reqL%v resp:%v\n", r.URL, bc.body.String())
	}
}


func MiddlewareWithAnotherService(ok bool) rest.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if ok {
				fmt.Println("ok!")
			}
			next(w, r)
		}
	}
}