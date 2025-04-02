package corebanking

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"go.bankyaya.org/app/backend/internal/pkg/config"
)

const (
	tokenEndpoint       = "/token"
	eodEndpoint         = "/api/ref/core-status"
	transactionEndpoint = "/api/transaction"
	overbookEndpoint    = "/api/transaction"
)

// Client represents a core banking client for interacting with the API.
type Client struct {
	httpClient *http.Client
	url        string
	username   string
	password   string
}

// NewClient initializes and returns a new Client instance with the given configuration.
func NewClient(cfg *config.Config, httpClient *http.Client) *Client {
	return &Client{
		httpClient: httpClient,
		url:        cfg.CoreBanking.URL,
		username:   cfg.CoreBanking.Username,
		password:   cfg.CoreBanking.Password,
	}
}

// EOD retrieves end-of-day status.
func (c *Client) EOD(ctx context.Context) (*EODResponse, error) {
	resp := new(EODResponse)
	err := c.executeRequest(ctx, http.MethodGet, eodEndpoint, nil, &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Inquiry retrieves bank account details.
func (c *Client) Inquiry(ctx context.Context, accountNumber string) (*InquiryResponse, error) {
	req := map[string]string{
		"noRekening":    accountNumber,
		"tipeTransaksi": "inquiry",
	}
	resp := new(InquiryResponse)
	err := c.executeRequest(ctx, http.MethodPost, transactionEndpoint, req, resp)
	if err != nil {
		return nil, err
	}
	return resp, err
}

// Overbook initiates an overbooking transaction.
func (c *Client) Overbook(ctx context.Context, req OverbookRequest) (*OverbookResponse, error) {
	resp := new(OverbookResponse)
	err := c.executeRequest(ctx, http.MethodPost, overbookEndpoint, req, &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// token gets the authentication token required for API calls.
func (c *Client) token(ctx context.Context) (string, error) {
	authURL := c.url + tokenEndpoint

	data := url.Values{}
	data.Set("username", c.username)
	data.Set("password", c.password)
	data.Set("grant_type", "password")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, authURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	res, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", errors.New(res.Status)
	}

	resp := new(AuthResponse)
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
