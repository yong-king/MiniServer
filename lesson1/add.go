package main

import "fmt"

// 本地调用
func Add(a int, b int) int {
	return a + b
}

func main() {
	a := 10
	b := 20
	fmt.Println(Add(a, b))
}
