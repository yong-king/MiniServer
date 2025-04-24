package connect

import (
	"net/http"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

var client = &http.Client{
	Transport: &http.Transport{
		DisableKeepAlives: true,
	},
	Timeout: 2 * time.Second,
}

// http.Client 是 Go 的标准库 net/http 中的一个结构体，表示一个 HTTP 客户端。
// Transport 是 http.Client 的一个字段，它定义了如何发送 HTTP 请求。
// 这里通过 http.Transport 来创建传输层配置。
// DisableKeepAlives: true：禁用 HTTP 的 keep-alive 机制。
// keep-alive 会让连接在多次请求之间保持打开状态，避免每次请求都建立新的连接。
// 这里禁用 keep-alive，意味着每次请求后连接都会被关闭。
// Timeout: 2 * time.Second：设置客户端的请求超时时间为 2 秒。
// 如果请求在 2 秒内没有完成，将会返回超时错误。

// Get 判断长链接是否能请求通
func Get(url string) bool {
	// Get 函数的作用是尝试发送一个 HTTP GET 请求到给定的 url，然后判断这个请求是否成功。
	resp, err := client.Get(url)
	if err != nil {
		logx.Errorw(
			"connect client.Get failed",
			logx.LogField{Key: "longUrl", Value: url},
			logx.LogField{Key: "err", Value: err.Error()},
		)
		return false
	}
	// resp.Body.Close()：关闭响应体，避免资源泄漏。每个响应体的 Body 必须被关闭。
	resp.Body.Close()
	return resp.StatusCode == http.StatusOK // 别人给我发一个跳转响应这里也不算过
}
