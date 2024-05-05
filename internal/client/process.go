package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type SpadeProcessCreateRequest struct {
	Code         string                 `json:"code"`
	Description  string                 `json:"description"`
	Tags         []string               `json:"tags"`
	Executor     int64                  `json:"executor"`
	SystemParams map[string]interface{} `json:"system_params"`
	UserParams   map[string]interface{} `json:"user_params"`
}

type SpadeProcessReadResponse struct {
	Id           int64                  `json:"id"`
	Code         string                 `json:"code"`
	Description  string                 `json:"description"`
	Tags         []string               `json:"tags"`
	Executor     int64                  `json:"executor"`
	SystemParams map[string]interface{} `json:"system_params"`
	UserParams   map[string]interface{} `json:"user_params"`
}

func (c *SpadeClient) CreateProcess(code, description string, tags []string, executor int64, systemParams, userParams map[string]interface{}) (*SpadeProcessReadResponse, error) {
	httpReqBody, err := json.Marshal(SpadeProcessCreateRequest{
		Code:         code,
		Description:  description,
		Tags:         tags,
		Executor:     executor,
		SystemParams: systemParams,
		UserParams:   userParams,
	})
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest(
		"POST",
		c.ApiUrl+"/api/v1/processes",
		bytes.NewBuffer(httpReqBody),
	)
	httpReq.Header.Set("Authorization", "Bearer "+c.Token)
	httpReq.Header.Set("Content-Type", "application/json")
	if err != nil {
		return nil, err
	}
	httpResp, err := c.HttpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()
	resp := SpadeProcessReadResponse{}
	err = json.NewDecoder(httpResp.Body).Decode(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *SpadeClient) ReadProcess(id int64) (*SpadeProcessReadResponse, error) {
	httpReq, err := http.NewRequest(
		"GET",
		c.ApiUrl+"/api/v1/processes/"+fmt.Sprint(id),
		nil,
	)
	httpReq.Header.Set("Authorization", "Bearer "+c.Token)
	httpReq.Header.Set("Content-Type", "application/json")
	if err != nil {
		return nil, err
	}

	httpResp, err := c.HttpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	resp := SpadeProcessReadResponse{}
	err = json.NewDecoder(httpResp.Body).Decode(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *SpadeClient) UpdateProcess(id int64, code, description string, tags []string, executor int64, systemParams, userParams map[string]interface{}) (*SpadeProcessReadResponse, error) {
	httpReqBody, err := json.Marshal(SpadeProcessCreateRequest{
		Code:         code,
		Description:  description,
		Tags:         tags,
		Executor:     executor,
		SystemParams: systemParams,
		UserParams:   userParams,
	})
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest(
		"PATCH",
		c.ApiUrl+"/api/v1/processes/"+fmt.Sprint(id),
		bytes.NewBuffer(httpReqBody),
	)
	httpReq.Header.Set("Authorization", "Bearer "+c.Token)
	httpReq.Header.Set("Content-Type", "application/json")
	if err != nil {
		return nil, err
	}

	httpResp, err := c.HttpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	resp := SpadeProcessReadResponse{}
	err = json.NewDecoder(httpResp.Body).Decode(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *SpadeClient) DeleteProcess(id int64) error {
	url, err := url.Parse(c.ApiUrl + "/api/v1/processes/" + fmt.Sprint(id))
	if err != nil {
		return err
	}
	httpReq := &http.Request{
		Method: "DELETE",
		URL:    url,
		Header: map[string][]string{
			"Authorization": {"Bearer " + c.Token},
			"Content-Type":  {"application/json"},
		},
	}
	httpResp, err := c.HttpClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer httpResp.Body.Close()
	return nil
}
