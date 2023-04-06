package main

import (
	"fmt"
	"runtime"
	"syscall"
	"time"
)

func main() {
	runtime.GOMAXPROCS(1)

	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				fmt.Println(time.Now())
			}
		}
	}()

	var l int = 100
	var msec int = -1
	epfd, _, _ := syscall.RawSyscall(syscall.SYS_EPOLL_CREATE1, 0, 0, 0)
	syscall.RawSyscall(syscall.SYS_EPOLL_CTL, epfd, 0, 0)
	syscall.RawSyscall6(syscall.SYS_EPOLL_WAIT, epfd, 0, uintptr(l), uintptr(msec), 0, 0)
}
