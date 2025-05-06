package uptrace

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	BaseURL = "https://api2.uptrace.dev"
)

type UptraceClient struct {
	BaseURL   string
	ProjectID string
	Token     string
	Client    *http.Client
}

func NewUptraceClient(projectID, token string) *UptraceClient {
	return &UptraceClient{
		BaseURL:   BaseURL,
		ProjectID: projectID,
		Token:     token,
		Client:    &http.Client{},
	}
}

func (u *UptraceClient) do(ctx context.Context, method, endpoint string, body []byte, out any) error {
	url := u.BaseURL + endpoint

	var reqBody io.Reader
	if body != nil {
		reqBody = bytes.NewBuffer(body)
		tflog.Debug(ctx, "Uptrace request", map[string]any{
			"method": method,
			"url":    url,
			"body":   string(body),
		})
	} else {
		tflog.Debug(ctx, "Uptrace request (no body)", map[string]any{
			"method": method,
			"url":    url,
		})
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+u.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := u.Client.Do(req)
	if err != nil {
		return fmt.Errorf("performing request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response: %w", err)
	}

	tflog.Debug(ctx, "Uptrace response", map[string]any{
		"status": resp.Status,
		"body":   string(respBody),
	})

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status %s: %s", resp.Status, respBody)
	}

	if out != nil {
		if err := json.Unmarshal(respBody, out); err != nil {
			return fmt.Errorf("decoding response JSON: %w", err)
		}
	}

	return nil
}

func (u *UptraceClient) GetMonitors(ctx context.Context) (*GetMonitorsResponse, error) {
	endpoint := fmt.Sprintf("/internal/v1/projects/%s/monitors", u.ProjectID)

	var result GetMonitorsResponse
	if err := u.do(ctx, "GET", endpoint, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (u *UptraceClient) GetMonitorById(ctx context.Context, id string) (*GetMonitorByIdResponse, error) {
	endpoint := fmt.Sprintf("/internal/v1/projects/%s/monitors/%s", u.ProjectID, id)

	var result GetMonitorByIdResponse
	if err := u.do(ctx, "GET", endpoint, nil, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (u *UptraceClient) CreateMonitor(ctx context.Context, req MonitorRequest) (*MonitorIdResponse, error) {
	endpoint := fmt.Sprintf("/internal/v1/projects/%s/monitors", u.ProjectID)

	if req.Type != "metric" {
		return nil, fmt.Errorf("error monitors must be of type \"metric\". You provided: %s", req.Type)
	}

	var result MonitorIdResponse
	if err := u.do(ctx, "POST", endpoint, nil, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (u *UptraceClient) UpdateMonitor(ctx context.Context, id string, req MonitorRequest) (*MonitorIdResponse, error) {
	endpoint := fmt.Sprintf("/internal/v1/projects/%s/monitors/%s", u.ProjectID, id)

	var result MonitorIdResponse
	if err := u.do(ctx, "PUT", endpoint, nil, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (u *UptraceClient) CreateErrorMonitor(ctx context.Context, req MonitorRequest) (*MonitorIdResponse, error) {
	endpoint := fmt.Sprintf("/internal/v1/projects/%s/monitors", u.ProjectID)

	if req.Type != "error" {
		return nil, fmt.Errorf("error monitors must be of type \"error\". You provided: %s", req.Type)
	}

	var result MonitorIdResponse
	if err := u.do(ctx, "POST", endpoint, nil, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (u *UptraceClient) UpdateErrorMonitor(ctx context.Context, id string, req MonitorRequest) (*MonitorIdResponse, error) {
	endpoint := fmt.Sprintf("/internal/v1/projects/%s/monitors/%s", u.ProjectID, id)

	if req.Type != "error" {
		return nil, fmt.Errorf("error monitors must be of type \"error\". You provided: %s", req.Type)
	}

	var result MonitorIdResponse
	if err := u.do(ctx, "PUT", endpoint, nil, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (u *UptraceClient) DeleteMonitor(ctx context.Context, id string) error {
	endpoint := fmt.Sprintf("/internal/v1/projects/%s/monitors/%s", u.ProjectID, id)

	var result any
	if err := u.do(ctx, "DELETE", endpoint, nil, &result); err != nil {
		return err
	}
	return nil
}
