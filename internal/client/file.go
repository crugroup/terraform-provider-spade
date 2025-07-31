package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type SpadeFileCreateRequest struct {
	Code          string                 `json:"code"`
	Description   string                 `json:"description"`
	Tags          []string               `json:"tags"`
	Format        int64                  `json:"format"`
	Processor     int64                  `json:"processor"`
	SystemParams  map[string]interface{} `json:"system_params"`
	UserParams    map[string]interface{} `json:"user_params"`
	LinkedProcess *int64                 `json:"linked_process"`
	VariableSets  []int64                `json:"variable_sets"`
}

type SpadeFileReadResponse struct {
	Id            int64                  `json:"id"`
	Code          string                 `json:"code"`
	Description   string                 `json:"description"`
	Tags          []string               `json:"tags"`
	Format        int64                  `json:"format"`
	Processor     int64                  `json:"processor"`
	SystemParams  map[string]interface{} `json:"system_params"`
	UserParams    map[string]interface{} `json:"user_params"`
	LinkedProcess int64                  `json:"linked_process"`
	VariableSets  []int64                `json:"variable_sets"`
}

func (c *SpadeClient) CreateFile(code, description string, tags []string, format, processor int64, systemParams, userParams map[string]interface{}, linkedProcess int64, variableSets []int64) (*SpadeFileReadResponse, error) {
	linkedProcessPtr := &linkedProcess
	if linkedProcess == 0 {
		linkedProcessPtr = nil
	}
	httpReqBody, err := json.Marshal(SpadeFileCreateRequest{
		Code:          code,
		Description:   description,
		Tags:          tags,
		Format:        format,
		Processor:     processor,
		SystemParams:  systemParams,
		UserParams:    userParams,
		LinkedProcess: linkedProcessPtr,
		VariableSets:  variableSets,
	})
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest(
		"POST",
		c.ApiUrl+"/api/v1/files",
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
	if !(httpResp.StatusCode >= 200 && httpResp.StatusCode < 300) {
		bodyData, err := io.ReadAll(httpResp.Body)
		if err != nil {
			return nil, fmt.Errorf("create file failed with status code %d", httpResp.StatusCode)
		}
		return nil, fmt.Errorf("create file failed with status code %d, response %s", httpResp.StatusCode, string(bodyData))
	}
	resp := SpadeFileReadResponse{}
	err = json.NewDecoder(httpResp.Body).Decode(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *SpadeClient) ReadFile(id int64) (*SpadeFileReadResponse, error) {
	httpReq, err := http.NewRequest(
		"GET",
		c.ApiUrl+"/api/v1/files/"+fmt.Sprint(id),
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
	if httpResp.StatusCode == 404 {
		return nil, nil
	}
	if !(httpResp.StatusCode >= 200 && httpResp.StatusCode < 300) {
		bodyData, err := io.ReadAll(httpResp.Body)
		if err != nil {
			return nil, fmt.Errorf("read file failed with status code %d", httpResp.StatusCode)
		}
		return nil, fmt.Errorf("read file failed with status code %d, response %s", httpResp.StatusCode, string(bodyData))
	}
	resp := SpadeFileReadResponse{}
	err = json.NewDecoder(httpResp.Body).Decode(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *SpadeClient) UpdateFile(id int64, code, description string, tags []string, format, processor int64, systemParams, userParams map[string]interface{}, linkedProcess int64, variableSets []int64) (*SpadeFileReadResponse, error) {
	linkedProcessPtr := &linkedProcess
	if linkedProcess == 0 {
		linkedProcessPtr = nil
	}
	httpReqBody, err := json.Marshal(SpadeFileCreateRequest{
		Code:          code,
		Description:   description,
		Tags:          tags,
		Format:        format,
		Processor:     processor,
		SystemParams:  systemParams,
		UserParams:    userParams,
		LinkedProcess: linkedProcessPtr,
		VariableSets:  variableSets,
	})
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest(
		"PATCH",
		c.ApiUrl+"/api/v1/files/"+fmt.Sprint(id),
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
	if !(httpResp.StatusCode >= 200 && httpResp.StatusCode < 300) {
		bodyData, err := io.ReadAll(httpResp.Body)
		if err != nil {
			return nil, fmt.Errorf("update file failed with status code %d", httpResp.StatusCode)
		}
		return nil, fmt.Errorf("update file failed with status code %d, response %s", httpResp.StatusCode, string(bodyData))
	}
	resp := SpadeFileReadResponse{}
	err = json.NewDecoder(httpResp.Body).Decode(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *SpadeClient) DeleteFile(id int64) error {
	url, err := url.Parse(c.ApiUrl + "/api/v1/files/" + fmt.Sprint(id))
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
	if !(httpResp.StatusCode >= 200 && httpResp.StatusCode < 300) {
		bodyData, err := io.ReadAll(httpResp.Body)
		if err != nil {
			return fmt.Errorf("delete file failed with status code %d", httpResp.StatusCode)
		}
		return fmt.Errorf("delete file failed with status code %d, response %s", httpResp.StatusCode, string(bodyData))
	}
	return nil
}
