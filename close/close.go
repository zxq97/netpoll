package main

import "net"

func main() {
	conn, err := net.Dial("tcp", ":8000")
	if err != nil {
		panic(err)
	}
	if err = conn.Close(); err != nil {
		panic(err)
	}
}
