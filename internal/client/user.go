package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type SpadeUserCreateRequest struct {
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Email     string  `json:"email"`
	IsActive  bool    `json:"is_active"`
	Groups    []int64 `json:"groups"`
}

type SpadeUserReadResponse struct {
	Id        int64   `json:"id"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Email     string  `json:"email"`
	IsActive  bool    `json:"is_active"`
	Groups    []int64 `json:"groups"`
}

func (c *SpadeClient) CreateUser(firstName, lastName, email string, isActive bool, groups []int64) (*SpadeUserReadResponse, error) {
	httpReqBody, err := json.Marshal(SpadeUserCreateRequest{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		IsActive:  isActive,
		Groups:    groups,
	})
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest(
		"POST",
		c.ApiUrl+"/api/v1/users",
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
			return nil, fmt.Errorf("create user failed with status code %d", httpResp.StatusCode)
		}
		return nil, fmt.Errorf("create user failed with status code %d, response %s", httpResp.StatusCode, string(bodyData))
	}
	resp := SpadeUserReadResponse{}
	err = json.NewDecoder(httpResp.Body).Decode(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *SpadeClient) ReadUser(id int64) (*SpadeUserReadResponse, error) {
	httpReq, err := http.NewRequest(
		"GET",
		c.ApiUrl+"/api/v1/users/"+fmt.Sprint(id),
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
			return nil, fmt.Errorf("read user failed with status code %d", httpResp.StatusCode)
		}
		return nil, fmt.Errorf("read user failed with status code %d, response %s", httpResp.StatusCode, string(bodyData))
	}
	resp := SpadeUserReadResponse{}
	err = json.NewDecoder(httpResp.Body).Decode(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *SpadeClient) UpdateUser(id int64, firstName, lastName, email string, isActive bool, groups []int64) (*SpadeUserReadResponse, error) {
	httpReqBody, err := json.Marshal(SpadeUserCreateRequest{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		IsActive:  isActive,
		Groups:    groups,
	})
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest(
		"PATCH",
		c.ApiUrl+"/api/v1/users/"+fmt.Sprint(id),
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
			return nil, fmt.Errorf("update user failed with status code %d", httpResp.StatusCode)
		}
		return nil, fmt.Errorf("update user failed with status code %d, response %s", httpResp.StatusCode, string(bodyData))
	}
	resp := SpadeUserReadResponse{}
	err = json.NewDecoder(httpResp.Body).Decode(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *SpadeClient) DeleteUser(id int64) error {
	url, err := url.Parse(c.ApiUrl + "/api/v1/users/" + fmt.Sprint(id))
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
			return fmt.Errorf("delete user failed with status code %d", httpResp.StatusCode)
		}
		return fmt.Errorf("delete user failed with status code %d, response %s", httpResp.StatusCode, string(bodyData))
	}
	return nil
}
