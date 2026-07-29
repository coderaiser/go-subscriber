//go:build linux || darwin

package formatter

import (
	"os"
	"syscall"
	"unsafe"
)

func termWidth() int {
	type winsize struct {
		Row, Col        uint16
		Xpixel, Ypixel uint16
	}
	ws := &winsize{}
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, os.Stderr.Fd(), syscall.TIOCGWINSZ, uintptr(unsafe.Pointer(ws)))
	if errno != 0 || ws.Col == 0 {
		return 80
	}
	return int(ws.Col)
}
