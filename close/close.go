package main

import (
	"time"

	"github.com/cloudwego/netpoll"
)

func main() {
	conn, err := netpoll.DialConnection("tcp", ":8000", time.Millisecond*100)
	if err != nil {
		panic(err)
	}
	if err = conn.Close(); err != nil {
		panic(err)
	}
}
