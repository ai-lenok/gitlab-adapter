package maintainer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ai-lenok/gitlab-adapter/properties"
	"net/http"
)

type Maintainer struct {
	Config *properties.GitLabConfig
	Client HttpClient
}

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func (m Maintainer) Request(method string,
	url string,
	data map[string]string) (*http.Response, error) {
	jsonData, _ := json.Marshal(data)
	buffer := bytes.NewBuffer(jsonData)

	req, _ := http.NewRequest(method, url, buffer)
	req.Header.Add("Authorization", "Bearer "+m.Config.AuthToken)
	req.Header.Add("Content-Type", "application/json")
	return m.Client.Do(req)
}

// https://docs.gitlab.com/ee/api/projects.html#create-project
func (m Maintainer) CreateRepo(req *properties.ReqCreateRepo) (*properties.RespCreateRepo, error) {
	data := map[string]string{
		"name":         req.Name,
		"path":         req.Path,
		"description":  req.Description,
		"namespace_id": req.NamespaceId,
		//"import_url":   importUrl,
	}

	resp, err := m.Request("POST", m.Config.UrlProjects(), data)
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
		return nil, fmt.Errorf("Bad status: %d", resp.StatusCode)
	}
}

// https://docs.gitlab.com/ee/api/projects.html#delete-project
func (m Maintainer) DeleteRepo(req *properties.ReqDeleteRepo) (*http.Response, error) {
	resp, err := m.Request("DELETE", m.Config.UrlDeleteProject(req.ProjectId), map[string]string{})
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusAccepted {
		return resp, err
	} else {
		return nil, fmt.Errorf("Bad status: %s", resp.Status)
	}
}

func (m Maintainer) LastBuildIsSuccess(req *properties.ReqListPipelines) (bool, error) {
	resp, err := m.ListPipelines(req)
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
func (m Maintainer) ListPipelines(req *properties.ReqListPipelines) (*http.Response, error) {
	return m.Request("GET", m.Config.UrlPipelines(req.ProjectId), map[string]string{})
}
