package e2etest

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/liushuangls/go-server-template/dto/response"
)

type Envelope struct {
	Code    int             `json:"code"`
	Message string          `json:"msg"`
	Data    json.RawMessage `json:"data"`
}

type Client struct {
	BaseURL         string
	HTTP            *http.Client
	explicitBaseURL bool
}

func NewClientFromEnv() *Client {
	baseURL, ok := os.LookupEnv("E2E_BASE_URL")
	if !ok || baseURL == "" {
		baseURL = "http://127.0.0.1:8082"
	}
	return &Client{
		BaseURL: baseURL,
		HTTP: &http.Client{
			Timeout: 10 * time.Second,
		},
		explicitBaseURL: ok && baseURL != "",
	}
}

func (c *Client) buildURL(path string, query url.Values) (string, error) {
	u, err := url.Parse(c.BaseURL)
	if err != nil {
		return "", fmt.Errorf("parse base url: %w", err)
	}
	ref, err := url.Parse(path)
	if err != nil {
		return "", fmt.Errorf("parse path: %w", err)
	}
	out := u.ResolveReference(ref)
	if query != nil {
		out.RawQuery = query.Encode()
	}
	return out.String(), nil
}

func (c *Client) Get(ctx context.Context, path string, query url.Values) (int, []byte, error) {
	fullURL, err := c.buildURL(path, query)
	if err != nil {
		return 0, nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return 0, nil, err
	}
	rsp, err := c.HTTP.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer rsp.Body.Close()

	body, err := io.ReadAll(rsp.Body)
	if err != nil {
		return rsp.StatusCode, nil, err
	}
	return rsp.StatusCode, body, nil
}

func decodeEnvelope[T any](body []byte) (Envelope, *T, error) {
	var env Envelope
	if err := json.Unmarshal(body, &env); err != nil {
		return Envelope{}, nil, err
	}
	if len(env.Data) == 0 || string(env.Data) == "null" {
		return env, nil, nil
	}
	var data T
	if err := json.Unmarshal(env.Data, &data); err != nil {
		return Envelope{}, nil, err
	}
	return env, &data, nil
}

type API struct {
	Health *HealthAPI
}

func NewAPI(c *Client) *API {
	return &API{
		Health: &HealthAPI{c: c},
	}
}

type HealthAPI struct {
	c *Client
}

func (h *HealthAPI) Get(ctx context.Context, message string) (int, Envelope, *response.HealthResp, error) {
	q := url.Values{}
	if message != "" {
		q.Set("message", message)
	}
	status, body, err := h.c.Get(ctx, "/health", q)
	if err != nil {
		return 0, Envelope{}, nil, err
	}
	env, data, err := decodeEnvelope[response.HealthResp](body)
	if err != nil {
		return status, Envelope{}, nil, err
	}
	return status, env, data, nil
}
