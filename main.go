package main

import (
	"fmt"
	"qq-go-/Server"
)

func main() {
	newserver := Server.NewServer("127.0.0.1", 9999)
	fmt.Println("22", newserver)
	newserver.Start()
}
