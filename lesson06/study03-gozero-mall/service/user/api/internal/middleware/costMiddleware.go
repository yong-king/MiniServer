package middleware

import (
	"net/http"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)




type CostMiddleware struct {
}

func NewCostMiddleware() *CostMiddleware {
	return &CostMiddleware{}
}

func (m *CostMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO generate middleware implement function, delete after code implementation
		// 中间件逻辑
		now := time.Now()
		// Passthrough to next handler if need
		next(w, r) // 实际执行的后续接口handler处理函数
		logx.Infof("-->cost:%v\n", time.Since(now))
	}
}
