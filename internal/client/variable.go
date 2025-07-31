package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type SpadeVariableCreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Value       string `json:"value"`
	IsSecret    bool   `json:"is_secret"`
}
type SpadeVariableUpdateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Value       string `json:"value"`
}

type SpadeVariableReadResponse struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Value       string `json:"value"`
	IsSecret    bool   `json:"is_secret"`
}

func (c *SpadeClient) CreateVariable(name, description, value string, isSecret bool) (*SpadeVariableReadResponse, error) {
	httpReqBody, err := json.Marshal(SpadeVariableCreateRequest{
		Name:        name,
		Description: description,
		Value:       value,
		IsSecret:    isSecret,
	})
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest(
		"POST",
		c.ApiUrl+"/api/v1/variables",
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
			return nil, fmt.Errorf("create variable failed with status code %d", httpResp.StatusCode)
		}
		return nil, fmt.Errorf("create variable failed with status code %d, response %s", httpResp.StatusCode, string(bodyData))
	}
	resp := SpadeVariableReadResponse{}
	err = json.NewDecoder(httpResp.Body).Decode(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *SpadeClient) ReadVariable(id int64) (*SpadeVariableReadResponse, error) {
	httpReq, err := http.NewRequest(
		"GET",
		c.ApiUrl+"/api/v1/variables/"+fmt.Sprint(id),
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
			return nil, fmt.Errorf("read variable failed with status code %d", httpResp.StatusCode)
		}
		return nil, fmt.Errorf("read variable failed with status code %d, response %s", httpResp.StatusCode, string(bodyData))
	}
	resp := SpadeVariableReadResponse{}
	err = json.NewDecoder(httpResp.Body).Decode(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *SpadeClient) UpdateVariable(id int64, name, description, value string) (*SpadeVariableReadResponse, error) {
	httpReqBody, err := json.Marshal(SpadeVariableUpdateRequest{
		Name:        name,
		Description: description,
		Value:       value,
	})
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest(
		"PATCH",
		c.ApiUrl+"/api/v1/variables/"+fmt.Sprint(id),
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
			return nil, fmt.Errorf("update variable failed with status code %d", httpResp.StatusCode)
		}
		return nil, fmt.Errorf("update variable failed with status code %d, response %s", httpResp.StatusCode, string(bodyData))
	}
	resp := SpadeVariableReadResponse{}
	err = json.NewDecoder(httpResp.Body).Decode(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *SpadeClient) DeleteVariable(id int64) error {
	url, err := url.Parse(c.ApiUrl + "/api/v1/variables/" + fmt.Sprint(id))
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
			return fmt.Errorf("delete variable failed with status code %d", httpResp.StatusCode)
		}
		return fmt.Errorf("delete variable failed with status code %d, response %s", httpResp.StatusCode, string(bodyData))
	}
	return nil
}

type SpadeVariableSearchResponse struct {
	Results []SpadeVariableReadResponse `json:"results"`
}

func (c *SpadeClient) SearchVariable(name string) (*SpadeVariableReadResponse, error) {
	httpReq, err := http.NewRequest(
		"GET",
		c.ApiUrl+"/api/v1/variables?search="+url.QueryEscape(name),
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
			return nil, fmt.Errorf("search variable failed with status code %d", httpResp.StatusCode)
		}
		return nil, fmt.Errorf("search variable failed with status code %d, response %s", httpResp.StatusCode, string(bodyData))
	}
	resp := SpadeVariableSearchResponse{}
	err = json.NewDecoder(httpResp.Body).Decode(&resp)
	if err != nil {
		return nil, err
	}
	for _, v := range resp.Results {
		if v.Name == name {
			return &v, nil
		}
	}
	return nil, fmt.Errorf("cannot find variable with name: %s", name)
}
