package main

import (
	"log"
	"syscall"

	"github.com/zxq97/netpoll/internal/socket"
)

func read(fd int) error {
	buf := make([]byte, 1024)
	n, err := syscall.Read(fd, buf)
	log.Println("read", string(buf[:n]))
	return err
}

func write(fd int) error {
	n, err := syscall.Write(fd, []byte("qwerty"))
	log.Println("write", fd, n)
	return err
}

func closeFD(epfd, fd int) {
	if err := syscall.EpollCtl(epfd, syscall.EPOLL_CTL_DEL, fd, &syscall.EpollEvent{Fd: int32(fd), Events: syscall.EPOLLOUT}); err != nil {
		log.Println("EpollCtl", fd, err)
	}
	if err := syscall.Close(fd); err != nil {
		log.Println("Close", fd, err)
	}
}

func main() {
	fd, err := socket.ListenTCP(8000)
	if err != nil {
		panic(err)
	}

	epfd, err := syscall.EpollCreate(4096)
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
				if err = read(int(evs[i].Fd)); err != nil {
					log.Println("readhandle", evs[i].Fd, err)
					closeFD(epfd, int(evs[i].Fd))
					continue
				}
				if err = syscall.EpollCtl(epfd, syscall.EPOLL_CTL_MOD, int(evs[i].Fd), &syscall.EpollEvent{Fd: evs[i].Fd, Events: syscall.EPOLLOUT | syscall.EPOLLERR | syscall.EPOLLHUP}); err != nil {
					log.Println("EpollCtl", evs[i].Fd, err)
				}
			} else if evs[i].Events&syscall.EPOLLOUT != 0 {
				if err = write(int(evs[i].Fd)); err != nil {
					log.Println("writehandle", evs[i].Fd, err)
					closeFD(epfd, int(evs[i].Fd))
					continue
				}
				if err = syscall.EpollCtl(epfd, syscall.EPOLL_CTL_MOD, int(evs[i].Fd), &syscall.EpollEvent{Fd: evs[i].Fd, Events: syscall.EPOLLIN | syscall.EPOLLERR | syscall.EPOLLHUP}); err != nil {
					log.Println("EpollCtl", evs[i].Fd, err)
				}
			} else if evs[i].Events&syscall.EPOLLHUP != 0 {
				closeFD(epfd, int(evs[i].Fd))
			}
		}
	}
}
