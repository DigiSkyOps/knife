package util

import (
	"agent/middleware"
	"github.com/gin-gonic/gin"
)

//var (
//	Engine *gin.Engine
//)
//
//func init(){
//	Engine = InitGin(GinLoggerHandler(InitGinLogger()))
//}

// 初始化gin
func InitGin(LoggerHandler gin.HandlerFunc,Mode string) *gin.Engine {
	// zerolog 方式
	//zerolog.SetGlobalLevel(zerolog.InfoLevel)
	//if gin.IsDebugging() {
	//	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	//}

	//f, _ := os.Create("gin.log")
	//subLog := zerolog.New(os.Stdout).With().Caller().Logger().Output(io.MultiWriter(f, os.Stdout))

	//engine := gin.New()
	//engine.Use(logger.SetLogger(logger.Config{
	//	Logger:  &subLog,
	//	UTC:     true,
	//}))
	// logrus 方式
	if Mode == "release"{
		gin.SetMode(gin.ReleaseMode)
	}
	engine := gin.New()
	engine.Use(LoggerHandler)
	engine.Use(gin.Recovery())
	engine.Use(middleware.Cors())
	//engine.Use(middleware.Auth())
	return engine
}
