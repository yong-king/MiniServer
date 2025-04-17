package transport

import (
	"context"
	"encoding/json"
	"net/http"
)

// 编码
func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
