package connect

import (
	"testing"

	c "github.com/smartystreets/goconvey/convey"// 别名导入
) 

func TestGet(t *testing.T) {
	c.Convey("基础用例", t, func() {
		targetUrl := "https://space.bilibili.com/290551837"
		got := Get(targetUrl)
		c.So(got, c.ShouldEqual, true) // 断言
	})

	c.Convey("url请求不通过的示例", t, func() {
		targetUrl := "ysh/dlrb/"
		got := Get(targetUrl)
		c.ShouldBeFalse(got) // 断言
	})
}
