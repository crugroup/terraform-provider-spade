package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type SpadeFileProcessorCreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Callable    string `json:"callable"`
}

type SpadeFileProcessorReadResponse struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Callable    string `json:"callable"`
}

func (c *SpadeClient) CreateFileProcessor(name, description, callable string) (*SpadeFileProcessorReadResponse, error) {
	httpReqBody, err := json.Marshal(SpadeFileProcessorCreateRequest{
		Name:        name,
		Description: description,
		Callable:    callable,
	})
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest(
		"POST",
		c.ApiUrl+"/api/v1/fileprocessors",
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
			return nil, fmt.Errorf("create file processor failed with status code %d", httpResp.StatusCode)
		}
		return nil, fmt.Errorf("create file processor failed with status code %d, response %s", httpResp.StatusCode, string(bodyData))
	}
	resp := SpadeFileProcessorReadResponse{}
	err = json.NewDecoder(httpResp.Body).Decode(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *SpadeClient) ReadFileProcessor(id int64) (*SpadeFileProcessorReadResponse, error) {
	httpReq, err := http.NewRequest(
		"GET",
		c.ApiUrl+"/api/v1/fileprocessors/"+fmt.Sprint(id),
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
			return nil, fmt.Errorf("read file processor failed with status code %d", httpResp.StatusCode)
		}
		return nil, fmt.Errorf("read file processor failed with status code %d, response %s", httpResp.StatusCode, string(bodyData))
	}
	resp := SpadeFileProcessorReadResponse{}
	err = json.NewDecoder(httpResp.Body).Decode(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *SpadeClient) UpdateFileProcessor(id int64, name, description, callable string) (*SpadeFileProcessorReadResponse, error) {
	httpReqBody, err := json.Marshal(SpadeFileProcessorCreateRequest{
		Name:        name,
		Description: description,
		Callable:    callable,
	})
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest(
		"PATCH",
		c.ApiUrl+"/api/v1/fileprocessors/"+fmt.Sprint(id),
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
			return nil, fmt.Errorf("update file processor failed with status code %d", httpResp.StatusCode)
		}
		return nil, fmt.Errorf("update file processor failed with status code %d, response %s", httpResp.StatusCode, string(bodyData))
	}
	resp := SpadeFileProcessorReadResponse{}
	err = json.NewDecoder(httpResp.Body).Decode(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *SpadeClient) DeleteFileProcessor(id int64) error {
	url, err := url.Parse(c.ApiUrl + "/api/v1/fileprocessors/" + fmt.Sprint(id))
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
			return fmt.Errorf("delete file processor failed with status code %d", httpResp.StatusCode)
		}
		return fmt.Errorf("delete file processor failed with status code %d, response %s", httpResp.StatusCode, string(bodyData))
	}
	return nil
}
