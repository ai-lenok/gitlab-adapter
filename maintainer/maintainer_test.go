package maintainer_test

import (
	"encoding/json"
	"github.com/ai-lenok/gitlab-adapter/maintainer"
	"github.com/ai-lenok/gitlab-adapter/properties"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
)

func TestCreateRepoSuccess(t *testing.T) {
	err := mockResponse("response/createRepoSuccess.json")
	assert.Nil(t, err)
	m := generateTestMaintainer()
	resp, err := m.CreateRepo(&properties.ReqCreateRepo{})
	assert.NotNil(t, resp)
	assert.Nil(t, err)
	assert.EqualValues(t, 46690708, resp.Id)
	assert.EqualValues(t, "test-name", resp.Name)
	assert.EqualValues(t, "test-path", resp.Path)
	assert.EqualValues(t, "Test repo", resp.Description)
	assert.EqualValues(t, "https://gitlab.com/java-school-courses/testing/test-path.git", resp.HttpUrlToRepo)
	assert.EqualValues(t, "git@gitlab.com:java-school-courses/testing/test-path.git", resp.SshUrlToRepo)
	assert.EqualValues(t, 10598392, resp.CreatorId)
	assert.EqualValues(t, "https://gitlab.com/java-school-courses/testing/test-path", resp.WebUrl)
	assert.EqualValues(t, "java-school-courses/testing/test-path", resp.PathWithNamespace)
}

func TestCreateRepoBadToken(t *testing.T) {
	err := mockResponse("response/createRepoBadToken.json")
	assert.Nil(t, err)
	m := generateTestMaintainer()
	resp, err := m.CreateRepo(&properties.ReqCreateRepo{})
	assert.Nil(t, resp)
	assert.NotNil(t, err)
	assert.EqualValues(t, "Bad status: 401", err.Error())
}

func TestDeleteRepoSuccess(t *testing.T) {
	err := mockResponse("response/deleteRepoSuccess.json")
	assert.Nil(t, err)
	m := generateTestMaintainer()
	resp, err := m.DeleteRepo(&properties.ReqDeleteRepo{})
	assert.NotNil(t, resp)
	assert.Nil(t, err)
	assert.EqualValues(t, "202 Accepted", resp.Message)
}

func TestDeleteRepoBadToken(t *testing.T) {
	err := mockResponse("response/deleteRepoBadToken.json")
	assert.Nil(t, err)
	m := generateTestMaintainer()
	resp, err := m.DeleteRepo(&properties.ReqDeleteRepo{})
	assert.Nil(t, resp)
	assert.NotNil(t, err)
	assert.EqualValues(t, "Bad status: 401", err.Error())
}

func TestVerifyPipelineSuccessTrue(t *testing.T) {
	err := mockResponseArray("response/verifyPipelineSuccessTrue.json")
	assert.Nil(t, err)
	m := generateTestMaintainer()
	resp, err := m.LastBuildIsSuccess(&properties.ReqListPipelines{})
	assert.Nil(t, err)
	assert.True(t, resp)
}

func TestVerifyPipelineSuccessFalse(t *testing.T) {
	err := mockResponseArray("response/verifyPipelineSuccessFalse.json")
	assert.Nil(t, err)
	m := generateTestMaintainer()
	resp, err := m.LastBuildIsSuccess(&properties.ReqListPipelines{})
	assert.Nil(t, err)
	assert.False(t, resp)
}

func TestVerifyPipelineBadToken(t *testing.T) {
	err := mockResponseArray("response/verifyPipelineBadToken.json")
	assert.Nil(t, err)
	m := generateTestMaintainer()
	resp, err := m.LastBuildIsSuccess(&properties.ReqListPipelines{})
	assert.False(t, resp)
	assert.NotNil(t, err)
	assert.EqualValues(t, "Bad status: 401", err.Error())
}

type Response struct {
	Status int
	Body   map[string]interface{}
}

type ResponseArray struct {
	Status int
	Body   []map[string]interface{}
}

type MockClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

var (
	GetDoFunc func(req *http.Request) (*http.Response, error)
)

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	return GetDoFunc(req)
}

func generateTestMaintainer() maintainer.Maintainer {
	return maintainer.Maintainer{
		Client: &MockClient{},
		Config: &properties.GitLabConfig{
			Host:      "test_host",
			AuthToken: "test_token",
		},
	}
}

func mockResponse(pathToFile string) error {
	response, err := readResponse(pathToFile)
	if err != nil {
		return err
	}
	GetDoFunc = func(*http.Request) (*http.Response, error) {
		return response, nil
	}
	return nil
}

func readResponse(pathToFile string) (*http.Response, error) {
	file, err := os.ReadFile(pathToFile)
	if err != nil {
		return nil, err
	}

	var response Response
	err = json.Unmarshal(file, &response)
	body, err := json.Marshal(response.Body)
	if err != nil {
		return nil, err
	}
	strBody := string(body)
	return &http.Response{
		StatusCode: response.Status,
		Body:       io.NopCloser(strings.NewReader(strBody)),
	}, nil
}

func mockResponseArray(pathToFile string) error {
	response, err := readResponseArray(pathToFile)
	if err != nil {
		return err
	}
	GetDoFunc = func(*http.Request) (*http.Response, error) {
		return response, nil
	}
	return nil
}

func readResponseArray(pathToFile string) (*http.Response, error) {
	file, err := os.ReadFile(pathToFile)
	if err != nil {
		return nil, err
	}

	var response ResponseArray
	err = json.Unmarshal(file, &response)
	body, err := json.Marshal(response.Body)
	if err != nil {
		return nil, err
	}
	strBody := string(body)
	return &http.Response{
		StatusCode: response.Status,
		Body:       io.NopCloser(strings.NewReader(strBody)),
	}, nil
}
