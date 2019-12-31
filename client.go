package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", ":19903")
	if err != nil {
		return
	}

	defer conn.Close()

	n, err := conn.Write([]byte("123456"))
	if err != nil {
		return
	}
	fmt.Println("len: ", n)
	for {
		time.Sleep(3 * time.Second)
	}
}
