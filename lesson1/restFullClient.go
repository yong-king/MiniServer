package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type addParams struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type addResponse struct {
	Code int `json:"code"`
	Data int `json:"data"`
}

func main() {
	url := "http://localhost:9090/add"
	param := addParams{X: 10, Y: 20}

	paramBytes, _ := json.Marshal(param)
	resp, _ := http.Post(url, "application/json", bytes.NewBuffer(paramBytes))
	defer resp.Body.Close()

	all, _ := ioutil.ReadAll(resp.Body)
	var response addResponse
	json.Unmarshal(all, &response)
	fmt.Println(response)
}
