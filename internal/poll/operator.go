package poll

import "syscall"

func (ep *Epoll) AddRead(fd int) error {
	return syscall.EpollCtl(ep.fd, syscall.EPOLL_CTL_ADD, fd, &syscall.EpollEvent{Fd: int32(fd), Events: syscall.EPOLLIN})
}

func (ep *Epoll) AddWriteAndRead(fd int) error {
	return syscall.EpollCtl(ep.fd, syscall.EPOLL_CTL_ADD, fd, &syscall.EpollEvent{Fd: int32(fd), Events: syscall.EPOLLIN | syscall.EPOLLOUT})
}

func (ep *Epoll) ModRead(fd int) error {
	return syscall.EpollCtl(ep.fd, syscall.EPOLL_CTL_MOD, fd, &syscall.EpollEvent{Fd: int32(fd), Events: syscall.EPOLLIN})
}

func (ep *Epoll) ModWriteAndRead(fd int) error {
	return syscall.EpollCtl(ep.fd, syscall.EPOLL_CTL_MOD, fd, &syscall.EpollEvent{Fd: int32(fd), Events: syscall.EPOLLIN | syscall.EPOLLOUT})
}

func (ep *Epoll) Del(fd int) error {
	return syscall.EpollCtl(ep.fd, syscall.EPOLL_CTL_DEL, fd, &syscall.EpollEvent{Fd: int32(fd), Events: syscall.EPOLLIN | syscall.EPOLLOUT})
}
