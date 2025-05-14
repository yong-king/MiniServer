package main

import (
	"math/rand"
	"net/http"
	"time"
)


func doGet(){
	for {
		_, _ = http.Get("http://127.0.0.1:8083/ping")
		time.Sleep(time.Duration(rand.Intn(1000)+800)*time.Millisecond)
	}
}