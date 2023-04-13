package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Config struct {
	M2MToken     string
	Environment  string
	ClientID     string
	ClientSecret string
}

type M2MTokenPayload struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Audience     string `json:"audience"`
	GrantType    string `json:"grant_type"`
}

type M2MTokenResponse struct {
	AccessToken string `json:"access_token"`
}

var environments = map[string]Environment{
	"dev": {
		Host:        "http://api-spec-service.dev.int.revelup.io",
		M2MTokenURL: "https://authentication.dev.revelup.com/oauth/token",
	},
	"qa": {
		Host:        "http://api-spec-service.qa.int.revelup.io",
		M2MTokenURL: "https://authentication.qa.revelup.com/oauth/token",
	},
	"prod": {
		Host:        "http://api-spec-service.int.revelup.io",
		M2MTokenURL: "https://authentication.revelup.com/oauth/token",
	},
}

func (c *Config) Client() (*Client, error) {
	httpClient := &http.Client{}

	client := &Client{
		M2MToken:    c.M2MToken,
		Environment: environments[c.Environment],
		HTTPClient:  httpClient,
	}

	if c.M2MToken == "" {
		if c.ClientID == "" || c.ClientSecret == "" {
			return nil, fmt.Errorf("client_id and client_secret must be provided if m2m_token was not")
		}

		m2mToken, err := c.getM2MToken(httpClient)
		if err != nil {
			return nil, err
		}

		client.M2MToken = *m2mToken
	}

	return client, nil
}

func (c *Config) getM2MToken(client *http.Client) (*string, error) {
	url := environments[c.Environment].M2MTokenURL
	method := "POST"

	payload := M2MTokenPayload{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		Audience:     "https://api-spec-service.revelup.io",
		GrantType:    "client_credentials",
	}
	reqBody, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Add("content-type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		return nil, fmt.Errorf("%s", resp.Body)
	}

	var m2mTokenResp M2MTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&m2mTokenResp); err != nil {
		return nil, err
	}

	return &m2mTokenResp.AccessToken, nil
}
