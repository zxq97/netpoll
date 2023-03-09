package main

import (
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", ":8000")
	if err != nil {
		panic(err)
	}
	for i := 0; i < 10; i++ {
		if _, err = conn.Write([]byte("123456")); err != nil {
			panic(err)
		}
		buf := make([]byte, 1024)
		if _, err = conn.Read(buf); err != nil {
			panic(err)
		}
		fmt.Println(string(buf))
	}
	if err = conn.Close(); err != nil {
		panic(err)
	}
}
