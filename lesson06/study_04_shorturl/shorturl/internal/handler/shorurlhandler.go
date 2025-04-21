package handler

import (
	"net/http"

	"shorurl/internal/logic"
	"shorurl/internal/svc"
	"shorurl/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func ShorurlHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Request
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewShorurlLogic(r.Context(), svcCtx)
		resp, err := l.Shorurl(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			// httpx.OkJsonCtx(r.Context(), w, resp)

			// w.Header().Set("location", resp.LongURL) // 重定向
			// w.WriteHeader(http.StatusFound) // 状态码

			http.Redirect(w, r, resp.LongURL, http.StatusFound) // 重定向
		}
	}
}
