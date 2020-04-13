package file

import (
	"agent/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type SPath struct {
	Path string `form:"filepath" binding:"required"`
}

type SFile struct {
	Path string
	ModTime time.Time
	Size int64
}

func HandleGetFile(c *gin.Context) {
	var filePath SPath
	if err := c.ShouldBindQuery(&filePath); err != nil {
		util.Logger.Error(err.Error())
		c.JSON(200, gin.H{
			"code": -200,
			"data": "缺少参数:" + err.Error(),
		})
		return
	}
	util.Logger.Infof("文件目录：%s", filePath.Path)
	var s []SFile

	s, err := GetAllFile(filePath.Path, s)
	if err != nil{
		c.JSON(200, gin.H{
			"code": -200,
			"data": "获取失败:" + err.Error(),
		})
	}else{
		c.JSON(200, gin.H{
			"code": 200,
			"data": s,
		})
	}
}

func HandleGetLastFile(c *gin.Context){
	var filePath SPath
	if err := c.ShouldBindQuery(&filePath); err != nil {
		util.Logger.Error(err.Error())
		c.JSON(200, gin.H{
			"code": -200,
			"data": "缺少参数:" + err.Error(),
		})
		return
	}
	util.Logger.Infof("文件目录：%s", filePath.Path)
	var s []SFile

	s, err := GetAllFile(filePath.Path, s)
	if err != nil{
		c.JSON(200, gin.H{
			"code": -200,
			"data": "获取失败:" + err.Error(),
		})
	}else{
		sort.Slice(s, func(i, j int) bool {
			return s[i].ModTime.After(s[j].ModTime)
		})
		c.JSON(200, gin.H{
			"code": 200,
			"data": s[0],
		})
	}
}

func GetAllFile(pathname string, s []SFile) ([]SFile, error) {
	//fromSlash := filepath.FromSlash(pathname)
	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		fmt.Println("read dir fail:", err)
		return s, err
	}
	for _, fi := range rd {
		if fi.IsDir() {
			fullDir := filepath.Join(pathname, fi.Name())
			s, err = GetAllFile(fullDir, s)
			if err != nil {
				fmt.Println("read dir fail:", err)
				return s, err
			}
		} else {
			fullName := strings.ReplaceAll(filepath.Join(pathname, fi.Name()), "\\", "/")
			finfo, _ := os.Stat(fullName)
			s = append(s, SFile{
				Path: fullName,
				ModTime: finfo.ModTime(),
				Size: finfo.Size(),
			})
		}
	}
	return s, nil
}
