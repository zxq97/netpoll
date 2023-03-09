package socket

import (
	"syscall"
)

// ListenTCP port 8000
func ListenTCP(port int) (int, error) {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err = syscall.SetNonblock(fd, true); err != nil {
		return 0, err
	}
	if err != nil {
		return 0, err
	}
	if err = syscall.Bind(fd, &syscall.SockaddrInet4{Port: port, Addr: [4]byte{127, 0, 0, 1}}); err != nil {
		return 0, err
	}
	err = syscall.Listen(fd, 1024)
	return fd, err
}

func Accept(fd int) (int, error) {
	nfd, _, err := syscall.Accept(fd)
	if err != nil {
		return 0, err
	}
	err = syscall.SetNonblock(nfd, true)
	return nfd, err
}
