package main

import (
	"log"
	"syscall"
	"time"

	"github.com/zxq97/netpoll/internal/socket"
)

func main() {
	fd, err := socket.ListenTCP(8000)
	if err != nil {
		panic(err)
	}

	epfd, err := syscall.EpollCreate1(0)
	if err != nil {
		panic(err)
	}
	ev := &syscall.EpollEvent{Fd: int32(fd), Events: syscall.EPOLLIN | syscall.EPOLLHUP | syscall.EPOLLERR}
	if err = syscall.EpollCtl(epfd, syscall.EPOLL_CTL_ADD, fd, ev); err != nil {
		panic(err)
	}

	for {
		evs := make([]syscall.EpollEvent, 1024)
		n, err := syscall.EpollWait(epfd, evs, 0)
		if err != nil {
			log.Println("EpollWait", err)
			continue
		}
		log.Println("EpollWait", n, err)
		for i := 0; i < n; i++ {
			if evs[i].Fd == int32(fd) {
				log.Println("accept")
				nfd, err := socket.Accept(fd)
				if err != nil {
					log.Println("Accept", err)
					continue
				}
				if err = syscall.EpollCtl(epfd, syscall.EPOLL_CTL_ADD, nfd, &syscall.EpollEvent{Fd: int32(nfd), Events: syscall.EPOLLIN | syscall.EPOLLHUP | syscall.EPOLLERR}); err != nil {
					log.Println("EpollCtl EPOLL_CTL_ADD", nfd, err)
					syscall.Close(nfd)
					continue
				}
			} else if evs[i].Events&syscall.EPOLLIN != 0 {
				buf := make([]byte, 4096)
				l, err := syscall.Read(int(evs[i].Fd), buf)
				if err != nil {
					log.Println("Read", evs[i].Fd, err)
					continue
				}
				log.Println("read", string(buf[:l]))
				if err = syscall.EpollCtl(epfd, syscall.EPOLL_CTL_MOD, int(evs[i].Fd), &syscall.EpollEvent{Fd: evs[i].Fd, Events: syscall.EPOLLERR | syscall.EPOLLHUP | syscall.EPOLLOUT}); err != nil {
					log.Println("EpollCtl", evs[i].Fd, err)
				}
			} else if evs[i].Events&syscall.EPOLLOUT != 0 {
				log.Println("write", evs[i].Fd, evs[i].Events)
			} else if evs[i].Events&syscall.EPOLLHUP != 0 {
				if err = syscall.EpollCtl(epfd, syscall.EPOLL_CTL_DEL, int(evs[i].Fd), &syscall.EpollEvent{Fd: evs[i].Fd, Events: syscall.EPOLLOUT}); err != nil {
					log.Println("EpollCtl EPOLL_CTL_DEL", evs[i].Fd, err)
				}
				if err = syscall.Close(int(evs[i].Fd)); err != nil {
					log.Println("Close", evs[i].Fd, err)
				}
			}
		}
		time.Sleep(time.Second * 5)
	}
}
