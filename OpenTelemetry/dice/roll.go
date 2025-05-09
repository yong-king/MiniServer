package main

import (
	"fmt"
	"math/rand"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var (
	tracer = otel.Tracer("roll")
)


func roll(w http.ResponseWriter, r *http.Request) {
	
	_, span := tracer.Start(r.Context(), "roll") // 开始 span
	defer span.End()                               // 结束 span

	// 业务逻辑
	number := 1 + rand.Intn(6)

	// 往span记录属性
	rollValueAttr := attribute.Int("roll.value", number)
	span.SetAttributes(rollValueAttr) // span 添加属性


	_, _ = fmt.Fprintln(w, number)
}
