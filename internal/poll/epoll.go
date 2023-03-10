package poll

import (
	"context"
	"log"
	"runtime"
	"syscall"

	"github.com/zxq97/netpoll/internal/socket"
)

type Epoll struct {
	fd  int
	out []byte
}

func OpenEpoll() (*Epoll, error) {
	fd, err := syscall.EpollCreate(1024)
	if err != nil {
		return nil, err
	}
	return &Epoll{fd: fd}, nil
}

func (ep *Epoll) Wait(ctx context.Context, fd int, fn func([]byte) ([]byte, error)) error {
	for {
		evs := make([]syscall.EpollEvent, 64)
		n, err := syscall.EpollWait(ep.fd, evs, 0)
		if err != nil {
			return err
		} else if n == 0 {
			runtime.Gosched()
			continue
		}
		for i := 0; i < n; i++ {
			if fd == int(evs[i].Fd) {
				// accept
				nfd, err := socket.Accept(fd)
				if err != nil {
					log.Println("Accept", fd, err)
					continue
				}
				if err = ep.AddRead(nfd); err != nil {
					log.Println("AddRead", nfd, err)
					socket.Close(nfd)
					continue
				}
			} else if evs[i].Events&syscall.EPOLLIN != 0 {
				// read
				buf := make([]byte, 1024)
				cfd := int(evs[i].Fd)
				l, err := syscall.Read(cfd, buf)
				if err != nil {
					log.Println("Read", cfd, err)
					socket.Close(cfd)
					continue
				}
				// handle
				out, err := fn(buf[:l])
				if err != nil {
					log.Println("handle", cfd, err)
					socket.Close(cfd)
					continue
				}
				ep.out = append(ep.out, out...)
				if err = ep.ModWriteAndRead(cfd); err != nil {
					log.Println("ModWriteAndRead", cfd, err)
					socket.Close(cfd)
					continue
				}
			} else if evs[i].Events&syscall.EPOLLOUT != 0 {
				if len(ep.out) != 0 {
					cfd := int(evs[i].Fd)
					l, err := syscall.Write(cfd, ep.out)
					if err != nil {
						log.Println("Write", cfd, err)
						socket.Close(cfd)
						continue
					}
					if l == len(ep.out) {
						ep.out = ep.out[:0]
					} else {
						ep.out = ep.out[l:]
					}
					if len(ep.out) == 0 {
						if err = ep.ModRead(cfd); err != nil {
							log.Println("ModRead", cfd, err)
							socket.Close(cfd)
							continue
						}
					}
				}
			}
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
	}
}
