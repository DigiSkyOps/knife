package util

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

type (
	HWND uintptr
)

var (
	libuser32 *windows.LazyDLL

	findWindowW          *windows.LazyProc
	getWindowTextW       *windows.LazyProc
	getWindowTextLengthW *windows.LazyProc
	getParent            *windows.LazyProc
	getSystemMetrics     *windows.LazyProc
	setWindowPos         *windows.LazyProc
	closeWindow          *windows.LazyProc
	destroyWindow        *windows.LazyProc
	postMessageW         *windows.LazyProc
	enumWindows          *windows.LazyProc
)

const (
	SM_CXSCREEN    = 0
	SM_CYSCREEN    = 1
	HWND_TOP       = HWND(0)
	HWND_TOPMOST   = ^HWND(0)
	SWP_SHOWWINDOW = 0x0040
	WM_CLOSE       = 16
)

func init() {
	libuser32 = windows.NewLazySystemDLL("user32.dll")

	findWindowW = libuser32.NewProc("FindWindowW")
	getWindowTextW = libuser32.NewProc("GetWindowTextW")
	getWindowTextLengthW = libuser32.NewProc("GetWindowTextLengthW")
	getParent = libuser32.NewProc("GetParent")
	getSystemMetrics = libuser32.NewProc("GetSystemMetrics")
	setWindowPos = libuser32.NewProc("SetWindowPos")
	closeWindow = libuser32.NewProc("CloseWindow")
	destroyWindow = libuser32.NewProc("DestroyWindow")
	postMessageW = libuser32.NewProc("PostMessageW")
	enumWindows = libuser32.NewProc("EnumWindows")
}

func FindWindow(lpClassName, lpWindowName *uint16) HWND {
	ret, _, _ := syscall.Syscall(findWindowW.Addr(), 2,
		uintptr(unsafe.Pointer(lpClassName)),
		uintptr(unsafe.Pointer(lpWindowName)),
		0,
	)

	return HWND(ret)
}

func GetWindowText(hWnd HWND) string {
	textLen := GetWindowTextLength(hWnd) + 1

	buf := make([]uint16, textLen)
	_, _, _ = syscall.Syscall(getWindowTextW.Addr(), 3,
		uintptr(hWnd),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(textLen),
	)

	return syscall.UTF16ToString(buf)
}

func GetWindowTextLength(hWnd HWND) int {
	ret, _, _ := syscall.Syscall(getWindowTextLengthW.Addr(), 1,
		uintptr(hWnd),
		0,
		0,
	)

	return int(ret)
}

func GetParent(hWnd HWND) HWND {
	ret, _, _ := syscall.Syscall(getParent.Addr(), 1,
		uintptr(hWnd),
		0,
		0)

	return HWND(ret)
}

func GetSystemMetrics(nIndex int32) int32 {
	ret, _, _ := syscall.Syscall(getSystemMetrics.Addr(), 1,
		uintptr(nIndex),
		0,
		0,
	)

	return int32(ret)
}

func SetWindowPos(hWnd, hWndInsertAfter HWND, x, y, width, height int32, flags uint32) bool {
	ret, _, _ := syscall.Syscall9(setWindowPos.Addr(), 7,
		uintptr(hWnd),
		uintptr(hWndInsertAfter),
		uintptr(x),
		uintptr(y),
		uintptr(width),
		uintptr(height),
		uintptr(flags),
		0,
		0)

	return ret != 0
}

func CloseWindow(hWnd HWND) bool {
	ret, _, _ := syscall.Syscall(closeWindow.Addr(), 1,
		uintptr(hWnd),
		0,
		0,
	)
	return ret != 0
}

func DestroyWindow(hWnd HWND) bool {
	ret, _, _ := syscall.Syscall(destroyWindow.Addr(), 1,
		uintptr(hWnd),
		0,
		0,
	)
	if ret == 0 {
		err := syscall.GetLastError()
		if err != nil {
			println(err.Error())
		}
	}
	return ret != 0
}

func PostMessage(hWnd HWND, msg, wParam, lParam uintptr) bool {
	ret, _, _ := syscall.Syscall6(postMessageW.Addr(), 4,
		uintptr(hWnd),
		uintptr(msg),
		uintptr(wParam),
		uintptr(lParam),
		0,
		0,
	)
	return ret != 0
}

func EnumWindows(enumFunc uintptr, lparam uintptr) (err error) {
	r1, _, e1 := syscall.Syscall(enumWindows.Addr(), 2,
		uintptr(enumFunc),
		uintptr(lparam),
		0,
	)
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}
