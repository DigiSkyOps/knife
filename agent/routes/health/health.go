package health

import (
	"agent/util"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
)

type SHealth struct {
	Data [][]interface{} `json:"data" binding:"required"`
}

func (conf *SHealth) String() string {
	b, err := json.Marshal(*conf)
	if err != nil {
		return fmt.Sprintf("%+v", *conf)
	}
	var out bytes.Buffer
	err = json.Indent(&out, b, "", "    ")
	if err != nil {
		return fmt.Sprintf("%+v", *conf)
	}
	return out.String()
}

func HandleHealth(c *gin.Context){
	var health SHealth
	if err := c.ShouldBindJSON(&health); err != nil {
		util.Logger.Errorf("Query Error: %s",err.Error())
		c.JSON(200, gin.H{
			"code": -200,
			"data": "缺少参数:" + err.Error(),
		})
		return
	}
	println(health.String())
}