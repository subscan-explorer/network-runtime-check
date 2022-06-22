package utils

import (
	"runtime"
	"syscall"
	"unsafe"
)

const (
	TiocgwinszOsx = 1074295912
)

type window struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

func TerminalWidth() int {
	w := new(window)
	tio := syscall.TIOCGWINSZ
	if runtime.GOOS == "darwin" {
		tio = TiocgwinszOsx
	}
	res, _, _ := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(tio),
		uintptr(unsafe.Pointer(w)),
	)
	if int(res) == -1 {
		return 0
	}
	return int(w.Col)
}
