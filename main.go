package main

import (
	"github.com/barrettzjh/DingDingAlertWebHook/loki"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	go loki.Server()

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.POST("/api/v2/alerts", func(c *gin.Context) {
		var alert []loki.LokiRuleAlertStruct
		if err := c.BindJSON(&alert); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"code": 10001,
				"error": err.Error(),
			})
			return
		}
		loki.AlertChan <- alert
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"code":    0,
		})
	})

	router.Run(":8080")
}
