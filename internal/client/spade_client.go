package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	if !(httpResp.StatusCode >= 200 && httpResp.StatusCode < 300) {
		bodyData, err := io.ReadAll(httpResp.Body)
		if err != nil {
			return fmt.Errorf("login failed with status code %d", httpResp.StatusCode)
		}
		return fmt.Errorf("login failed with status code %d, response %s", httpResp.StatusCode, string(bodyData))
	}
	resp := SpadeLoginResponse{}
	err = json.NewDecoder(httpResp.Body).Decode(&resp)
	if err != nil {
		return err
	}
	c.Token = resp.AccessToken
	return nil
}
