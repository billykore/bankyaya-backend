package qris

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"go.bankyaya.org/app/backend/pkg/config"
)

const (
	tokenEndpoint   = "/oauth/accesstoken?grant_type=client_credentials"
	inquiryEndpoint = "/v2/qris/inquiry-v2"
)

type Client struct {
	httpClient   *http.Client
	url          string
	clientId     string
	clientSecret string
}

func NewClient(cfg *config.Config, httpClient *http.Client) *Client {
	return &Client{
		httpClient:   httpClient,
		url:          cfg.QRIS.URL,
		clientId:     cfg.QRIS.ClientId,
		clientSecret: cfg.QRIS.ClientSecret,
	}
}

func (c *Client) Inquiry(ctx context.Context, req InquiryRequest) (*InquiryResponse, error) {
	resp := new(InquiryResponse)
	err := c.executeRequest(ctx, http.MethodPost, inquiryEndpoint, req, resp)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (c *Client) token(ctx context.Context) (string, error) {
	fullUrl := c.url + tokenEndpoint
	credential := fmt.Sprintf("%s:%s", c.clientId, c.clientSecret)
	basicAuth := base64.StdEncoding.EncodeToString([]byte(credential))

	r, err := http.NewRequestWithContext(ctx, http.MethodPost, fullUrl, nil)
	if err != nil {
		return "", err
	}

	r.Header.Add("Authorization", "Basic "+basicAuth)

	res, err := c.httpClient.Do(r)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", errors.New(res.Status)
	}

	resp := new(Token)
	if err = json.NewDecoder(res.Body).Decode(resp); err != nil {
		return "", err
	}

	return resp.AccessToken, nil
}

// executeRequest is a common helper for sending requests and handling responses.
func (c *Client) executeRequest(ctx context.Context, method, endpoint string, body any, response any) error {
	fullUrl := c.url + endpoint

	// Marshal body if provided
	var reqBody []byte
	var err error
	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			return err
		}
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, method, fullUrl, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	// Add token to the request
	token, err := c.token(ctx)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Add("Content-Type", "application/json")

	// Perform HTTP request
	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.New(res.Status)
	}

	// Decode response
	return json.NewDecoder(res.Body).Decode(response)
}
