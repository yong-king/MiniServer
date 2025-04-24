package urltool

import (
	"errors"
	"net/url"
	"path"
)

/*
url.Parse 是 Go 标准库中 net/url 包提供的一个函数，用于解析一个 URL 字符串并返回一个 *url.URL 结构体。


它的作用是解析一个完整的 URL 字符串，
将其分解成不同的部分（例如协议、主机、路径、查询字符串等），并将这些部分保存在一个结构体中，供后续操作使用。

url.URL ： Scheme（协议）、Host、Path、RawQuery（查询字符串部分）、Fragment（片段部分）（例如 #section1）

path.Base 是 Go 标准库中 path 包提供的一个函数，用于返回路径的“基本”部分，即路径的最后一部分（通常是文件名或文件夹名）。
path.Base 会返回给定路径字符串的最后一部分。如果路径以 / 结尾，它会返回路径中去除 / 后的最后部分。
*/

func GetBasePath(targetUrl string) (string, error) {
	// url.Parse(targetUrl) 会尝试解析传入的 URL 字符串。
	// 如果传入的字符串不是有效的 URL，它会返回一个错误。
	myUrl, err := url.Parse(targetUrl)
	if err != nil {
		return "", err
	}
	if len(myUrl.Host) == 0 {
		return "", errors.New("no host")
	}
	// path.Base(myUrl.Path) 来获取 URL 路径部分的“基本”部分。
	// myUrl.Path 返回 URL 的路径（比如 /path/to/file），然后 path.Base 会提取路径的最后一部分
	return path.Base(myUrl.Path), nil
}