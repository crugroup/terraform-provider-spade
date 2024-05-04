package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type SpadeClient struct {
	ApiUrl     string
	HttpClient *http.Client
	Token      string
}

type SpadeLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SpadeLoginResponse struct {
	AccessToken  string `json:"access"`
	RefreshToken string `json:"refresh"`
}

func (c *SpadeClient) Login(email, password string) error {
	httpReqBody, err := json.Marshal(SpadeLoginRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequest(
		"POST",
		c.ApiUrl+"/api/v1/token",
		bytes.NewBuffer(httpReqBody),
	)
	httpReq.Header.Set("Content-Type", "application/json")
	if err != nil {
		return err
	}
	httpResp, err := c.HttpClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer httpResp.Body.Close()
	resp := SpadeLoginResponse{}
	err = json.NewDecoder(httpResp.Body).Decode(&resp)
	if err != nil {
		return err
	}
	c.Token = resp.AccessToken
	return nil
}

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
	return nil
}

type SpadeFileFormatCreateRequest struct {
	Format string `json:"format"`
}

type SpadeFileFormatReadResponse struct {
	Id     int64  `json:"id"`
	Format string `json:"format"`
}

func (c *SpadeClient) CreateFileFormat(format string) (*SpadeFileFormatReadResponse, error) {
	httpReqBody, err := json.Marshal(SpadeFileFormatCreateRequest{
		Format: format,
	})
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest(
		"POST",
		c.ApiUrl+"/api/v1/fileformats",
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
	resp := SpadeFileFormatReadResponse{}
	err = json.NewDecoder(httpResp.Body).Decode(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *SpadeClient) ReadFileFormat(id int64) (*SpadeFileFormatReadResponse, error) {
	httpReq, err := http.NewRequest(
		"GET",
		c.ApiUrl+"/api/v1/fileformats/"+fmt.Sprint(id),
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
	resp := SpadeFileFormatReadResponse{}
	err = json.NewDecoder(httpResp.Body).Decode(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *SpadeClient) UpdateFileFormat(id int64, format string) (*SpadeFileFormatReadResponse, error) {
	httpReqBody, err := json.Marshal(SpadeFileFormatCreateRequest{
		Format: format,
	})
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest(
		"PATCH",
		c.ApiUrl+"/api/v1/fileformats/"+fmt.Sprint(id),
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
	resp := SpadeFileFormatReadResponse{}
	err = json.NewDecoder(httpResp.Body).Decode(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *SpadeClient) DeleteFileFormat(id int64) error {
	url, err := url.Parse(c.ApiUrl + "/api/v1/fileformats/" + fmt.Sprint(id))
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
	return nil
}

type SpadeFileCreateRequest struct {
	Code         string                 `json:"code"`
	Description  string                 `json:"description"`
	Tags         []string               `json:"tags"`
	Format       int64                  `json:"format"`
	Processor    int64                  `json:"processor"`
	SystemParams map[string]interface{} `json:"system_params"`
	UserParams   map[string]interface{} `json:"user_params"`
}

type SpadeFileReadResponse struct {
	Id           int64                  `json:"id"`
	Code         string                 `json:"code"`
	Description  string                 `json:"description"`
	Tags         []string               `json:"tags"`
	Format       int64                  `json:"format"`
	Processor    int64                  `json:"processor"`
	SystemParams map[string]interface{} `json:"system_params"`
	UserParams   map[string]interface{} `json:"user_params"`
}

func (c *SpadeClient) CreateFile(code, description string, tags []string, format, processor int64, systemParams, userParams map[string]interface{}) (*SpadeFileReadResponse, error) {
	httpReqBody, err := json.Marshal(SpadeFileCreateRequest{
		Code:         code,
		Description:  description,
		Tags:         tags,
		Format:       format,
		Processor:    processor,
		SystemParams: systemParams,
		UserParams:   userParams,
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
	resp := SpadeFileReadResponse{}
	err = json.NewDecoder(httpResp.Body).Decode(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *SpadeClient) UpdateFile(id int64, code, description string, tags []string, format, processor int64, systemParams, userParams map[string]interface{}) (*SpadeFileReadResponse, error) {
	httpReqBody, err := json.Marshal(SpadeFileCreateRequest{
		Code:         code,
		Description:  description,
		Tags:         tags,
		Format:       format,
		Processor:    processor,
		SystemParams: systemParams,
		UserParams:   userParams,
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
	return nil
}

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
