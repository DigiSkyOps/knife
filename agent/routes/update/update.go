package update

import (
	"agent/global"
	"agent/util"
	"github.com/gin-gonic/gin"
	"github.com/inconshreveable/go-update"
	"net/http"
)

type SUrl struct {
	Url string `json:"url" binding:"required"`
}

func HandleUpdate(c *gin.Context){
	var updateUrl SUrl
	if err := c.ShouldBindJSON(&updateUrl); err != nil {
		util.Logger.Error(err.Error())
		c.JSON(200, gin.H{
			"code": -200,
			"data": "缺少参数:" + err.Error(),
		})
		return
	}
	resp, err := http.Get(updateUrl.Url)
	if err != nil {
		c.JSON(200, gin.H{
			"code": -200,
			"data": "获取失败:" + err.Error(),
		})
	}
	defer resp.Body.Close()
	err = update.Apply(resp.Body, update.Options{})
	if err != nil {
		c.JSON(200, gin.H{
			"code": -200,
			"data": "升级失败:" + err.Error(),
		})
	}

	global.ServerSig <- 1

	c.JSON(200, gin.H{
		"code": 200,
		"data": "升级成功",
	})
}