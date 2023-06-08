package server

import (
	"fmt"
	"github.com/ai-lenok/gitlab-adapter/maintainer"
	"github.com/ai-lenok/gitlab-adapter/properties"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

var gitLabConfig *properties.GitLabConfig

func StartServer(config *properties.GitLabConfig, port int) error {
	gitLabConfig = config

	server := gin.Default()
	server.GET("/health", HealthCheck)
	server.POST("/api/v1/create-repo", creteRepo)
	server.POST("/api/v1/delete-repo", deleteRepo)
	server.POST("/api/v1/verify-pipeline-status", verifyPipelineStatus)

	return server.Run(fmt.Sprintf(":%d", port))
}

func creteRepo(context *gin.Context) {
	var reqCreateRepo properties.ReqCreateRepo

	if err := context.BindJSON(&reqCreateRepo); err != nil {
		return
	}

	resp, err := maintainer.CreateRepo(gitLabConfig, &reqCreateRepo)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(resp)

	context.IndentedJSON(http.StatusCreated, resp)
}

func deleteRepo(context *gin.Context) {
	var reqDeleteRepo properties.ReqDeleteRepo

	if err := context.BindJSON(&reqDeleteRepo); err != nil {
		return
	}
	resp, err := maintainer.DeleteRepo(gitLabConfig, &reqDeleteRepo)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	resp.Body.Close()

	log.Println(string(body))
	context.IndentedJSON(http.StatusAccepted, resp)
}

func verifyPipelineStatus(context *gin.Context) {
	var reqListPipelines properties.ReqListPipelines

	if err := context.BindJSON(&reqListPipelines); err != nil {
		return
	}

	isSuccess, err := maintainer.LastBuildIsSuccess(gitLabConfig, &reqListPipelines)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("is Success: " + strconv.FormatBool(isSuccess))
	if isSuccess {
		context.IndentedJSON(http.StatusNoContent, nil)
	} else {
		context.IndentedJSON(http.StatusConflict, nil)
	}
}

func HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "UP",
	})
}
