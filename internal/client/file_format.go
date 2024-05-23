package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

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
	if !(httpResp.StatusCode >= 200 && httpResp.StatusCode < 300) {
		bodyData, err := io.ReadAll(httpResp.Body)
		if err != nil {
			return nil, fmt.Errorf("create file format failed with status code %d", httpResp.StatusCode)
		}
		return nil, fmt.Errorf("create file format failed with status code %d, response %s", httpResp.StatusCode, string(bodyData))
	}
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
	if httpResp.StatusCode == 404 {
		return nil, nil
	}
	if !(httpResp.StatusCode >= 200 && httpResp.StatusCode < 300) {
		bodyData, err := io.ReadAll(httpResp.Body)
		if err != nil {
			return nil, fmt.Errorf("read file format failed with status code %d", httpResp.StatusCode)
		}
		return nil, fmt.Errorf("read file format failed with status code %d, response %s", httpResp.StatusCode, string(bodyData))
	}
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
	if !(httpResp.StatusCode >= 200 && httpResp.StatusCode < 300) {
		bodyData, err := io.ReadAll(httpResp.Body)
		if err != nil {
			return nil, fmt.Errorf("update file format failed with status code %d", httpResp.StatusCode)
		}
		return nil, fmt.Errorf("update file format failed with status code %d, response %s", httpResp.StatusCode, string(bodyData))
	}
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
	if !(httpResp.StatusCode >= 200 && httpResp.StatusCode < 300) {
		bodyData, err := io.ReadAll(httpResp.Body)
		if err != nil {
			return fmt.Errorf("delete file format failed with status code %d", httpResp.StatusCode)
		}
		return fmt.Errorf("delete file format failed with status code %d, response %s", httpResp.StatusCode, string(bodyData))
	}
	return nil
}
