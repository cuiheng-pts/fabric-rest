package web

import (
	"fabric-rest/fabric"
	"github.com/gin-gonic/gin"
	"net/http"
)

var initInfo *fabric.InitInfo

func SetInitInfo(info *fabric.InitInfo) {
	initInfo = info
	return
}

type CreateRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func Create(c *gin.Context) {
	req := &CreateRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	err := fabric.Create(initInfo, req.Key, req.Value)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{})
}

func Query(c *gin.Context) {
	key := c.Query("key")
	query, err := fabric.Query(initInfo, key)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
	c.JSON(http.StatusOK, gin.H{"value": query})
}
