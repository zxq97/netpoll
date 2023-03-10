package main

import (
	"context"

	"github.com/zxq97/netpoll/internal/poll"
	"github.com/zxq97/netpoll/internal/socket"
)

func main() {
	fd, err := socket.ListenTCP(8000)
	if err != nil {
		panic(err)
	}

	ep, err := poll.OpenEpoll()
	if err != nil {
		panic(err)
	}

	if err = ep.AddRead(fd); err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()
	if err = ep.Wait(ctx, fd, func(in []byte) ([]byte, error) {
		return in, nil
	}); err != nil {
		panic(err)
	}
}
