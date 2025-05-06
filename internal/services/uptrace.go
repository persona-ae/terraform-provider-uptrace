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
		tflog.Error(ctx, "Error creating request", map[string]any{"error": err})
		return fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+u.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := u.Client.Do(req)
	if err != nil {
		tflog.Error(ctx, "Error performing request", map[string]any{"error": err})
		return fmt.Errorf("performing request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		tflog.Error(ctx, "Error reading response", map[string]any{"error": err})
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
		tflog.Debug(ctx, "Unmarshalling into", map[string]any{
			"type": fmt.Sprintf("%T", out),
		})
		if err := json.Unmarshal(respBody, out); err != nil {
			tflog.Error(ctx, "Error decoding JSON", map[string]any{"error": err})
			return fmt.Errorf("decoding response JSON: %w", err)
		}

		db, _ := json.Marshal(out)
		tflog.Debug(ctx, "debug out: "+string(db))
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
