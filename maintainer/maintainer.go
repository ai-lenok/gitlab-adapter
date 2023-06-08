package maintainer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ai-lenok/gitlab-adapter/properties"
	"net/http"
)

// https://docs.gitlab.com/ee/api/projects.html#create-project
func CreateRepo(config *properties.GitLabConfig,
	req *properties.ReqCreateRepo) (*properties.RespCreateRepo, error) {
	data := map[string]string{
		"name":         req.Name,
		"path":         req.Path,
		"description":  req.Description,
		"namespace_id": req.NamespaceId,
		//"import_url":   importUrl,
	}

	resp, err := request("POST", config.UrlProjects(), config, data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK ||
		resp.StatusCode == http.StatusCreated {
		var respCreateRepo properties.RespCreateRepo
		err = json.NewDecoder(resp.Body).Decode(&respCreateRepo)
		if err != nil {
			return nil, err
		}
		return &respCreateRepo, nil
	} else {
		return nil, fmt.Errorf("Bad status: %s", resp.Status)
	}
}

// https://docs.gitlab.com/ee/api/projects.html#delete-project
func DeleteRepo(config *properties.GitLabConfig, req *properties.ReqDeleteRepo) (*http.Response, error) {
	resp, err := request("DELETE", config.UrlDeleteProject(req.ProjectId), config, map[string]string{})
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusAccepted {
		return resp, err
	} else {
		return nil, fmt.Errorf("Bad status: %s", resp.Status)
	}
}

func LastBuildIsSuccess(conf *properties.GitLabConfig, req *properties.ReqListPipelines) (bool, error) {
	resp, err := ListPipelines(conf, req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	var pipelines []properties.RespListPipeline
	err = json.NewDecoder(resp.Body).Decode(&pipelines)
	if err != nil {
		return false, err
	}
	return pipelines[0].Status == "success", nil

}

// https://docs.gitlab.com/ee/api/pipelines.html#list-project-pipelines
func ListPipelines(config *properties.GitLabConfig, req *properties.ReqListPipelines) (*http.Response, error) {
	return request("GET", config.UrlPipelines(req.ProjectId), config, map[string]string{})
}

func request(method string,
	url string,
	config *properties.GitLabConfig,
	data map[string]string) (*http.Response, error) {
	jsonData, _ := json.Marshal(data)
	buffer := bytes.NewBuffer(jsonData)
	client := &http.Client{}

	req, _ := http.NewRequest(method, url, buffer)
	req.Header.Add("Authorization", "Bearer "+config.AuthToken)
	req.Header.Add("Content-Type", "application/json")
	return client.Do(req)
}
