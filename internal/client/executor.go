package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type SpadeExecutorCreateRequest struct {
	Name                    string `json:"name"`
	Description             string `json:"description"`
	Callable                string `json:"callable"`
	HistoryProviderCallable string `json:"history_provider_callable"`
}

type SpadeExecutorReadResponse struct {
	Id                      int64  `json:"id"`
	Name                    string `json:"name"`
	Description             string `json:"description"`
	Callable                string `json:"callable"`
	HistoryProviderCallable string `json:"history_provider_callable"`
}

func (c *SpadeClient) CreateExecutor(name, description, callable, historyProviderCallable string) (*SpadeExecutorReadResponse, error) {
	httpReqBody, err := json.Marshal(SpadeExecutorCreateRequest{
		Name:                    name,
		Description:             description,
		Callable:                callable,
		HistoryProviderCallable: historyProviderCallable,
	})
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest(
		"POST",
		c.ApiUrl+"/api/v1/executors",
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
			return nil, fmt.Errorf("create executor failed with status code %d", httpResp.StatusCode)
		}
		return nil, fmt.Errorf("create executor failed with status code %d, response %s", httpResp.StatusCode, string(bodyData))
	}
	resp := SpadeExecutorReadResponse{}
	err = json.NewDecoder(httpResp.Body).Decode(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *SpadeClient) ReadExecutor(id int64) (*SpadeExecutorReadResponse, error) {
	httpReq, err := http.NewRequest(
		"GET",
		c.ApiUrl+"/api/v1/executors/"+fmt.Sprint(id),
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
	if !(httpResp.StatusCode >= 200 && httpResp.StatusCode < 300) {
		bodyData, err := io.ReadAll(httpResp.Body)
		if err != nil {
			return nil, fmt.Errorf("read executor failed with status code %d", httpResp.StatusCode)
		}
		return nil, fmt.Errorf("read executor failed with status code %d, response %s", httpResp.StatusCode, string(bodyData))
	}
	resp := SpadeExecutorReadResponse{}
	err = json.NewDecoder(httpResp.Body).Decode(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *SpadeClient) UpdateExecutor(id int64, name, description, callable, historyProviderCallable string) (*SpadeExecutorReadResponse, error) {
	httpReqBody, err := json.Marshal(SpadeExecutorCreateRequest{
		Name:                    name,
		Description:             description,
		Callable:                callable,
		HistoryProviderCallable: historyProviderCallable,
	})
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest(
		"PATCH",
		c.ApiUrl+"/api/v1/executors/"+fmt.Sprint(id),
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
			return nil, fmt.Errorf("update executor failed with status code %d", httpResp.StatusCode)
		}
		return nil, fmt.Errorf("update executor failed with status code %d, response %s", httpResp.StatusCode, string(bodyData))
	}
	resp := SpadeExecutorReadResponse{}
	err = json.NewDecoder(httpResp.Body).Decode(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *SpadeClient) DeleteExecutor(id int64) error {
	url, err := url.Parse(c.ApiUrl + "/api/v1/executors/" + fmt.Sprint(id))
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
			return fmt.Errorf("delete executor failed with status code %d", httpResp.StatusCode)
		}
		return fmt.Errorf("delete executor failed with status code %d, response %s", httpResp.StatusCode, string(bodyData))
	}
	return nil
}
