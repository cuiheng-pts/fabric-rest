package main

import (
	"fabric-rest/fabric"
	"fabric-rest/web"
	"github.com/gin-gonic/gin"
)

const infoConfig = "info.yaml"
const sdkConfig = "config.yaml"

func main() {
	// read config file
	// init fabric sdk
	initInfo, err := fabric.ConstructorFromYaml(infoConfig)
	if err != nil {
		return
	}
	_, err = fabric.InitSDK(sdkConfig, false, initInfo)
	if err != nil {
		return
	}
	// start service
	web.SetInitInfo(initInfo)
	r := gin.Default()
	r.POST("/create", web.Create)
	r.GET("/query", web.Query)
	r.Run(":13000")
}
