package server

import (
	"net/http"

	kratoshttp "github.com/go-kratos/kratos/v2/transport/http"
	kratosstatus "github.com/go-kratos/kratos/v2/transport/http/status"
	"google.golang.org/grpc/status"
)

type httpResonse struct {
	Code int `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

func responseEncoder(w http.ResponseWriter, r *http.Request, v interface{}) error {
	if v == nil {
		return nil
	}
	if rd, ok := v.(kratoshttp.Redirector); ok {
		url, code := rd.Redirect()
		http.Redirect(w, r, url, code)
		return nil
	}
	codec, _ := kratoshttp.CodecForRequest(r, "Accept")

	// 构造自定义结构体
	resp := &httpResonse{
		Code: http.StatusOK,
		Msg: "success",
		Data: v,
	}

	data, err := codec.Marshal(resp)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/" + codec.Name())
	_, err = w.Write(data)
	if err != nil {
		return err
	}
	return nil
}

// DefaultErrorEncoder encodes the error to the HTTP response.
func errorEncoder(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		return
	}
	// 判断err是否已经是一个err类型
	// se := errors.FromError(err)
	resp := new(httpResonse)
	// 检查 err 是否是 gRPC 错误，从错误中提取出 gRPC 错误信息
	if gs, ok := status.FromError(err); ok{
		resp = &httpResonse{
			// httpstatus.FromGRPCCode 将其转换为 HTTP 状态码
			Code: kratosstatus.FromGRPCCode(gs.Code()),
			Msg: gs.Message(),
			Data: nil,
		}
	}else {
		resp = &httpResonse{
			Code: http.StatusInternalServerError, //500
			Msg: "内部错误",
			Data: nil,
		}
	}
	codec, _ := kratoshttp.CodecForRequest(r, "Accept")
	body, err := codec.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/" + codec.Name())
	w.WriteHeader(int(resp.Code))
	_, _ = w.Write(body)
}
