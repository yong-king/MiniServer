package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

// 参数
type addParam struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type addResult struct {
	Code int `json:"code"`
	Data int `json:"data"`
}

func add(x, y int) int {
	return x + y
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	// 解析参数
	all, _ := ioutil.ReadAll(r.Body)
	var param addParam
	json.Unmarshal(all, &param)
	// 业务逻辑
	ret := add(param.X, param.Y)
	respBytes, _ := json.Marshal(addResult{0, ret})
	w.Write(respBytes)
}

func main() {
	http.HandleFunc("/add", addHandler)
	log.Fatal(http.ListenAndServe(":9090", nil))
}
