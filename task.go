package qencode

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// CreateTaskResponse is the create task endpoint response.
type CreateTaskResponse struct {
	Error     int    `json:"error"`
	UploadURL string `json:"upload_url"`
	TaskToken string `json:"task_token"`
}

// CreateTask creates a new qEncode task.
func (c *Client) CreateTask(token string) (CreateTaskResponse, error) {
	createTaskEndpoint := c.baseURL + "/v1/create_task"

	bd := url.Values{}
	bd.Add("token", token)
	r, err := http.NewRequest(http.MethodPost, createTaskEndpoint, strings.NewReader(bd.Encode()))
	if err != nil {
		return CreateTaskResponse{}, err
	}
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.hClient.Do(r)
	if err != nil {
		return CreateTaskResponse{}, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode > 299 {
		b := new(bytes.Buffer)
		_, _ = io.Copy(b, resp.Body)
		return CreateTaskResponse{}, RequestError{
			Message:      "error creating qEncode task",
			ResponseBody: b,
			StatusCode:   resp.StatusCode,
		}
	}

	var createTaskResp CreateTaskResponse
	if err := json.NewDecoder(resp.Body).Decode(&createTaskResp); err != nil {
		return CreateTaskResponse{}, err
	}

	return createTaskResp, nil
}

// StartTaskResponse is the start task endpoint response.
type StartTaskResponse struct {
	StatusURL string `json:"status_url"`
	Error     int    `json:"error"`
}

// StartTask will start an encoding task. A task token generated by
// CreateTask or manually must be provided. An optional payload argument
// can be passed so that it can be accessed later on a callback/webhook.
// You can provide a json object formatted as a string.
// The StartTaskQuery is used to define the options of the encoding task.
func (c *Client) StartTask(taskToken, payload string, q StartTaskQuery) (StartTaskResponse, error) {
	startTaskEndpoint := c.baseURL + "/v1/start_encode2"

	query, err := json.Marshal(q)
	if err != nil {
		return StartTaskResponse{}, err
	}

	bd := url.Values{}
	bd.Add("task_token", taskToken)
	bd.Add("payload", payload)
	bd.Add("query", string(query))

	r, err := http.NewRequest(http.MethodPost, startTaskEndpoint, strings.NewReader(bd.Encode()))
	if err != nil {
		return StartTaskResponse{}, err
	}
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.hClient.Do(r)
	if err != nil {
		return StartTaskResponse{}, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode > 299 {
		b := new(bytes.Buffer)
		_, _ = io.Copy(b, resp.Body)
		return StartTaskResponse{}, RequestError{
			Message:      "error starting qEncode task",
			ResponseBody: b,
			StatusCode:   resp.StatusCode,
		}
	}

	var startTaskResp StartTaskResponse
	if err := json.NewDecoder(resp.Body).Decode(&startTaskResp); err != nil {
		return StartTaskResponse{}, err
	}

	return startTaskResp, nil
}
