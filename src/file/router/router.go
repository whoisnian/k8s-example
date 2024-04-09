package router

import (
	"github.com/gin-gonic/gin"
	"github.com/whoisnian/k8s-example/src/file/global"
)

func Setup() *gin.Engine {
	if global.CFG.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())
	engine.GET("/", func(c *gin.Context) { c.String(200, "ok") })
	return engine
}
