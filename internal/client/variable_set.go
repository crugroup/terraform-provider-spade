package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type SpadeVariableSetCreateRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Variables   []int64 `json:"variables"`
}

type SpadeVariableSetReadResponse struct {
	Id          int64   `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Variables   []int64 `json:"variables"`
}

func (c *SpadeClient) CreateVariableSet(name, description string, variables []int64) (*SpadeVariableSetReadResponse, error) {
	httpReqBody, err := json.Marshal(SpadeVariableSetCreateRequest{
		Name:        name,
		Description: description,
		Variables:   variables,
	})
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest(
		"POST",
		c.ApiUrl+"/api/v1/variable-sets",
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
			return nil, fmt.Errorf("create variable set failed with status code %d", httpResp.StatusCode)
		}
		return nil, fmt.Errorf("create variable set failed with status code %d, response %s", httpResp.StatusCode, string(bodyData))
	}
	resp := SpadeVariableSetReadResponse{}
	err = json.NewDecoder(httpResp.Body).Decode(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *SpadeClient) ReadVariableSet(id int64) (*SpadeVariableSetReadResponse, error) {
	httpReq, err := http.NewRequest(
		"GET",
		c.ApiUrl+"/api/v1/variable-sets/"+fmt.Sprint(id),
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
			return nil, fmt.Errorf("read variable set failed with status code %d", httpResp.StatusCode)
		}
		return nil, fmt.Errorf("read variable set failed with status code %d, response %s", httpResp.StatusCode, string(bodyData))
	}
	resp := SpadeVariableSetReadResponse{}
	err = json.NewDecoder(httpResp.Body).Decode(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *SpadeClient) UpdateVariableSet(id int64, name, description string, variables []int64) (*SpadeVariableSetReadResponse, error) {
	httpReqBody, err := json.Marshal(SpadeVariableSetCreateRequest{
		Name:        name,
		Description: description,
		Variables:   variables,
	})
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest(
		"PATCH",
		c.ApiUrl+"/api/v1/variable-sets/"+fmt.Sprint(id),
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
			return nil, fmt.Errorf("update variable set failed with status code %d", httpResp.StatusCode)
		}
		return nil, fmt.Errorf("update variable set failed with status code %d, response %s", httpResp.StatusCode, string(bodyData))
	}
	resp := SpadeVariableSetReadResponse{}
	err = json.NewDecoder(httpResp.Body).Decode(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *SpadeClient) DeleteVariableSet(id int64) error {
	url, err := url.Parse(c.ApiUrl + "/api/v1/variable-sets/" + fmt.Sprint(id))
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
			return fmt.Errorf("delete variable set failed with status code %d", httpResp.StatusCode)
		}
		return fmt.Errorf("delete variable set failed with status code %d, response %s", httpResp.StatusCode, string(bodyData))
	}
	return nil
}
