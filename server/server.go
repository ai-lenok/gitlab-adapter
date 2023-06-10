package server

import (
	"fmt"
	"github.com/ai-lenok/gitlab-adapter/maintainer"
	"github.com/ai-lenok/gitlab-adapter/properties"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

var m *maintainer.Maintainer

func StartServer(manager *maintainer.Maintainer, port int) error {
	log.Printf("Start server. Port: %d", port)
	m = manager

	server := gin.Default()
	server.GET("/health", HealthCheck)
	server.POST("/api/v1/project", creteRepo)
	server.DELETE("/api/v1/project", deleteRepo)
	server.POST("/api/v1/project/verify-pipeline", verifyPipelineStatus)

	return server.Run(fmt.Sprintf(":%d", port))
}

func creteRepo(context *gin.Context) {
	var reqCreateRepo properties.ReqCreateRepo

	if err := context.BindJSON(&reqCreateRepo); err != nil {
		log.Println(err)
		context.Status(http.StatusInternalServerError)
		return
	}

	resp, err := m.CreateRepo(&reqCreateRepo)
	if err != nil {
		log.Println(err)
		context.Status(http.StatusInternalServerError)
		return
	}

	log.Println(resp)

	context.IndentedJSON(http.StatusCreated, resp)
}

func deleteRepo(context *gin.Context) {
	var reqDeleteRepo properties.ReqDeleteRepo

	if err := context.BindJSON(&reqDeleteRepo); err != nil {
		log.Println(err)
		context.Status(http.StatusInternalServerError)
		return
	}
	resp, err := m.DeleteRepo(&reqDeleteRepo)
	if err != nil {
		log.Println(err)
		context.Status(http.StatusInternalServerError)
		return
	}
	log.Println(resp)
	context.Status(http.StatusAccepted)
}

func verifyPipelineStatus(context *gin.Context) {
	var reqListPipelines properties.ReqListPipelines

	if err := context.BindJSON(&reqListPipelines); err != nil {
		log.Println(err)
		context.Status(http.StatusInternalServerError)
		return
	}

	isSuccess, err := m.LastBuildIsSuccess(&reqListPipelines)
	if err != nil {
		log.Println(err)
		context.Status(http.StatusInternalServerError)
		return
	}
	log.Print("is Success: " + strconv.FormatBool(isSuccess))
	if isSuccess {
		context.Status(http.StatusNoContent)
	} else {
		context.Status(http.StatusConflict)
	}
}

func HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "UP",
	})
}
