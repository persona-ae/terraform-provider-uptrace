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
	APIKey    string
	Client    *http.Client
}

func NewUptraceClient(projectID, apiKey string) *UptraceClient {
	return &UptraceClient{
		BaseURL:   BaseURL,
		ProjectID: projectID,
		APIKey:    apiKey,
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

	req.Header.Set("Authorization", "Bearer "+u.APIKey)
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

func (u *UptraceClient) GetMonitors(ctx context.Context, out *GetMonitorsResponse) error {
	endpoint := fmt.Sprintf("/internal/v1/projects/%s/monitors", u.ProjectID)
	return u.do(ctx, "GET", endpoint, nil, out)
}

func (u *UptraceClient) GetMonitorById(ctx context.Context, id int32, out *GetMonitorByIdResponse) error {
	endpoint := fmt.Sprintf("/internal/v1/projects/%s/monitors/%d", u.ProjectID, id)
	return u.do(ctx, "GET", endpoint, nil, out)
}

func (u *UptraceClient) GetMonitorByIdStr(ctx context.Context, id string, out *GetMonitorByIdResponse) error {
	endpoint := fmt.Sprintf("/internal/v1/projects/%s/monitors/%s", u.ProjectID, id)
	return u.do(ctx, "GET", endpoint, nil, out)
}

func (u *UptraceClient) CreateMonitor(ctx context.Context, req Monitor, out *MonitorIdResponse) error {
	endpoint := fmt.Sprintf("/internal/v1/projects/%s/monitors", u.ProjectID)
	return u.do(ctx, "POST", endpoint, nil, out)
}

func (u *UptraceClient) UpdateMonitor(ctx context.Context, id int32, req Monitor, out *MonitorIdResponse) error {
	endpoint := fmt.Sprintf("/internal/v1/projects/%s/monitors/%d", u.ProjectID, id)
	return u.do(ctx, "PUT", endpoint, nil, out)
}

func (u *UptraceClient) DeleteMonitor(ctx context.Context, id int32) error {
	endpoint := fmt.Sprintf("/internal/v1/projects/%s/monitors/%d", u.ProjectID, id)

	var result any
	return u.do(ctx, "DELETE", endpoint, nil, &result)
}
