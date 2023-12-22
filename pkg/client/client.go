package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ws/livekitcloudanalytics/pkg/models"
)

type RateLimiter interface {
	Wait(ctx context.Context) error
}

type Client struct {
	token      string // still not totally sure if this should be a string or a *string
	baseUrl    string
	limiter    RateLimiter
	httpClient *http.Client
}

const (
	defaultBaseUrl = "https://cloud-api.livekit.io/api/project"
	userAgent      = "@ws/livekitcloudanalytics" // need to figure out how to include version/system info
)

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 5 * time.Second},
		baseUrl:    defaultBaseUrl,
	}
}

func (c *Client) WithToken(token string) *Client {
	c.token = token
	return c
}

func NewClientWithToken(token string) *Client {
	return NewClient().WithToken(token)
}

// at some point I could potentially see allowing passing apiKey/apiSecret
// and generating a token but that would require LK SDK as a dependency

func (c *Client) WithRateLimiter(limiter RateLimiter) *Client {
	c.limiter = limiter
	return c
}

// really can't imagine when/why you'd need to do this
// (since afaik there's no way to get the LK analytics API outside of their cloud offering)
// but figured I may as well leave the carveout if anyone needs it
func (c *Client) WithBaseURL(baseUrl string) *Client {
	c.baseUrl = baseUrl
	return c
}

func (c *Client) doRequest(method, url string, token string) (*http.Response, error) {
	if c.token == "" {
		return nil, errors.New("must pass a token (use .WithToken(token) or NewClientWithToken(token))")
	}

	if c.limiter != nil {
		if err := c.limiter.Wait(context.Background()); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error building request: %v", err)
	}

	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"Accept":        "application/json", // need to see if I can get protobufs from this API
		"User-Agent":    userAgent,
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close() // close the body if there's an error
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return resp, nil
}

func (c *Client) ListSessions(projectId string) (*models.ListSessionsResponse, error) {
	url := c.baseUrl + "/" + projectId + "/sessions"

	resp, err := c.doRequest("GET", url, c.token)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response models.ListSessionsResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &response, nil
}

func (c *Client) ListSessionDetails(projectId string, sessionId string) (*models.SessionDetailsResponse, error) {
	url := c.baseUrl + "/" + projectId + "/sessions/" + sessionId

	resp, err := c.doRequest("GET", url, c.token)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response models.SessionDetailsResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &response, nil
}
