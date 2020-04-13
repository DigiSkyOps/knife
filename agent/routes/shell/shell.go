package shell

import (
	"agent/util"
	"github.com/gin-gonic/gin"
	"os/exec"
	"strings"
	"syscall"
)

type SEWindow struct {
	Title string
	HWND util.HWND
}

type SShell struct {
	Command  string   `json:"command" binding:"required"`
	FilePath string   `json:"filepath" binding:"required"`
	Params   []string `json:"params"`
}

type SWindow struct {
	Window string `form:"window" binding:"required"`
}

func HandleShell(c *gin.Context) {
	var shell SShell
	if err := c.ShouldBindJSON(&shell); err != nil {
		util.Logger.Errorf("Query Error: %s", err.Error())
		c.JSON(200, gin.H{
			"code": -200,
			"data": "缺少参数:" + err.Error(),
		})
		return
	}
	state := startShell(shell)
	c.JSON(200, gin.H{
		"code": 200,
		"data": state,
	})
}

func startShell(shell SShell) bool {
	var cmd []string
	cmd = append(cmd, shell.FilePath)

	if len(shell.Params) > 0 {
		util.Logger.Infof("Command: %s %s %s", shell.Command, shell.FilePath, strings.Join(shell.Params, " "))
		for _, v := range shell.Params {
			cmd = append(cmd, v)
		}
	} else {
		util.Logger.Infof("Command: %s %s", shell.Command, shell.FilePath)
	}

	go func() {
		res := exec.Command(shell.Command, cmd...)
		_, err := res.Output()
		if err != nil {
			util.Logger.Errorf("Error: %s", err)
		}
		//util.Logger.Infof("OutPut: %s", string(o[:]))
	}()
	return true
}

func HandleEnumWindows(c *gin.Context){
	var windows []SEWindow
	cb := syscall.NewCallback(func(h syscall.Handle, p uintptr) uintptr {
		v := util.GetWindowText(util.HWND(h))
		if strings.Trim(v, " ") != ""{
			windows = append(windows, SEWindow{
				Title: v,
				HWND: util.HWND(p),
			})
		}

		return 1 // continue enumeration
	})
	util.EnumWindows(cb, 0)
	c.JSON(200,gin.H{
		"code": 200,
		"data": windows,
	})
}

func HandleFullWindow(c *gin.Context) {
	var window SWindow
	if err := c.ShouldBindQuery(&window); err != nil {
		util.Logger.Errorf("Query Error: %s", err.Error())
		c.JSON(200, gin.H{
			"code": -200,
			"data": "缺少参数:" + err.Error(),
		})
		return
	}

	screenX := util.GetSystemMetrics(util.SM_CXSCREEN)
	screenY := util.GetSystemMetrics(util.SM_CYSCREEN)

	windowName, _ := syscall.UTF16PtrFromString(window.Window)
	dw := util.FindWindow(nil, windowName)

	f := util.SetWindowPos(dw,util.HWND_TOPMOST,0,0,screenX,screenY,util.SWP_SHOWWINDOW)

	c.JSON(200, gin.H{
		"code": 200,
		"data": f,
	})
}

func HandleCloseWin(c *gin.Context){
	var window SWindow
	if err := c.ShouldBindQuery(&window); err != nil {
		util.Logger.Errorf("Query Error: %s", err.Error())
		c.JSON(200, gin.H{
			"code": -200,
			"data": "缺少参数:" + err.Error(),
		})
		return
	}

	windowName, _ := syscall.UTF16PtrFromString(window.Window)
	dw := util.FindWindow(nil, windowName)

	r := util.PostMessage(dw,util.WM_CLOSE,0,0)
	c.JSON(200, gin.H{
		"code": 200,
		"data": r,
	})
}