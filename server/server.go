package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func runRest(port int) error {
	server := gin.New()
	server.GET("/health", HealthCheck)
	return server.Run(fmt.Sprintf(":%d", port))
}

func HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "UP",
	})
}
