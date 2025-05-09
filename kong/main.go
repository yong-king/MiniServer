package main

import (
	"flag"
	"fmt"
	"net/http"
)



func main() {
	var (
		port int
		msg string
	)

	flag.IntVar(&port, "port", 9870, "http server port")
	flag.StringVar(&msg, "msg", "resp", "resp message")

	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, msg+" from %d", port)
	})

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Printf("http server failed, err:%v\n", err)
		return
	}
}