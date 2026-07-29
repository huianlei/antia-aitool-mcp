package jenkins

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"github.com/huianlei/antia-aitool-mcp/pkg/utils"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/net/publicsuffix"
)

// Client is a Jenkins REST API client
type Client struct {
	baseURL    string
	username   string
	password   string
	httpClient *http.Client
	logger     *zap.Logger
	maxRetries int
	retryDelay time.Duration
}

// NewClient creates a new Jenkins client
func NewClient(config Config, logger *zap.Logger) *Client {
	// Create cookie jar for session management
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		logger.Warn("failed to create cookie jar, session may not persist", zap.Error(err))
	}

	httpClient := utils.NewHTTPClient(config.Timeout, !config.VerifySSL)
	httpClient.Jar = jar

	return &Client{
		baseURL:    strings.TrimSuffix(config.URL, "/"),
		username:   config.Username,
		password:   config.Password,
		httpClient: httpClient,
		logger:     logger,
		maxRetries: config.MaxRetries,
		retryDelay: config.RetryDelay,
	}
}

// GetJobs retrieves all jobs from Jenkins
func (c *Client) GetJobs(ctx context.Context) ([]Job, error) {
	endpoint := "/api/json?tree=jobs[name,url,color,description]"

	var response JobsResponse
	if err := c.doRequest(ctx, "GET", endpoint, nil, &response); err != nil {
		return nil, errors.Wrap(err, "failed to get jobs")
	}

	return response.Jobs, nil
}

// GetJob retrieves detailed information about a specific job
func (c *Client) GetJob(ctx context.Context, jobName string) (*JobDetail, error) {
	endpoint := fmt.Sprintf("/job/%s/api/json", url.PathEscape(jobName))

	var job JobDetail
	if err := c.doRequest(ctx, "GET", endpoint, nil, &job); err != nil {
		return nil, errors.Wrapf(err, "failed to get job: %s", jobName)
	}

	return &job, nil
}

// GetCrumb retrieves the CSRF protection crumb from Jenkins
func (c *Client) GetCrumb(ctx context.Context) (*Crumb, error) {
	endpoint := "/crumbIssuer/api/json"

	var crumb Crumb
	if err := c.doRequest(ctx, "GET", endpoint, nil, &crumb); err != nil {
		// Some Jenkins instances don't have CSRF protection enabled
		c.logger.Debug("failed to get crumb, CSRF protection may be disabled", zap.Error(err))
		return nil, nil
	}

	return &crumb, nil
}

// TriggerBuild triggers a build for a job
func (c *Client) TriggerBuild(ctx context.Context, jobName string, params map[string]string) error {
	// Get crumb for CSRF protection
	crumb, err := c.GetCrumb(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get crumb")
	}

	var endpoint string
	if len(params) > 0 {
		// Build with parameters
		endpoint = fmt.Sprintf("/job/%s/buildWithParameters", url.PathEscape(jobName))

		// Build query string
		values := url.Values{}
		for key, value := range params {
			values.Add(key, value)
		}
		endpoint = endpoint + "?" + values.Encode()
	} else {
		// Build without parameters
		endpoint = fmt.Sprintf("/job/%s/build", url.PathEscape(jobName))
	}

	if err := c.doRequestWithCrumb(ctx, "POST", endpoint, nil, nil, crumb); err != nil {
		return errors.Wrapf(err, "failed to trigger build: %s", jobName)
	}

	c.logger.Info("build triggered", zap.String("job", jobName))
	return nil
}

// GetBuild retrieves information about a specific build
func (c *Client) GetBuild(ctx context.Context, jobName string, buildNumber int) (*BuildDetail, error) {
	endpoint := fmt.Sprintf("/job/%s/%d/api/json", url.PathEscape(jobName), buildNumber)

	var build BuildDetail
	if err := c.doRequest(ctx, "GET", endpoint, nil, &build); err != nil {
		return nil, errors.Wrapf(err, "failed to get build: %s #%d", jobName, buildNumber)
	}

	return &build, nil
}

// GetBuildLog retrieves the console log for a build
func (c *Client) GetBuildLog(ctx context.Context, jobName string, buildNumber int, start int) (string, error) {
	endpoint := fmt.Sprintf("/job/%s/%d/consoleText", url.PathEscape(jobName), buildNumber)

	// Add start parameter if specified
	if start > 0 {
		endpoint = fmt.Sprintf("%s?start=%d", endpoint, start)
	}

	body, err := c.doRawRequest(ctx, "GET", endpoint, nil, nil)
	if err != nil {
		return "", errors.Wrapf(err, "failed to get build log: %s #%d", jobName, buildNumber)
	}

	return string(body), nil
}

// GetBuilds retrieves build history for a job
func (c *Client) GetBuilds(ctx context.Context, jobName string, limit int) ([]Build, error) {
	// Use tree parameter to fetch specific fields
	tree := "builds[number,url,result,building,duration,timestamp]"
	endpoint := fmt.Sprintf("/job/%s/api/json?tree=%s", url.PathEscape(jobName), tree)

	var response struct {
		Builds []Build `json:"builds"`
	}

	if err := c.doRequest(ctx, "GET", endpoint, nil, &response); err != nil {
		return nil, errors.Wrapf(err, "failed to get builds: %s", jobName)
	}

	// Apply limit if specified
	if limit > 0 && len(response.Builds) > limit {
		return response.Builds[:limit], nil
	}

	return response.Builds, nil
}

// doRequest performs an HTTP request with JSON response
func (c *Client) doRequest(ctx context.Context, method, endpoint string, body io.Reader, result interface{}) error {
	respBody, err := c.doRawRequest(ctx, method, endpoint, body, nil)
	if err != nil {
		return err
	}

	// Parse JSON response if result is provided
	if result != nil {
		if err := json.Unmarshal(respBody, result); err != nil {
			return errors.Wrap(err, "failed to parse JSON response")
		}
	}

	return nil
}

// doRequestWithCrumb performs an HTTP request with CSRF crumb
func (c *Client) doRequestWithCrumb(ctx context.Context, method, endpoint string, body io.Reader, result interface{}, crumb *Crumb) error {
	respBody, err := c.doRawRequest(ctx, method, endpoint, body, crumb)
	if err != nil {
		return err
	}

	// Parse JSON response if result is provided
	if result != nil {
		if err := json.Unmarshal(respBody, result); err != nil {
			return errors.Wrap(err, "failed to parse JSON response")
		}
	}

	return nil
}

// doRawRequest performs an HTTP request with retry logic
func (c *Client) doRawRequest(ctx context.Context, method, endpoint string, body io.Reader, crumb *Crumb) ([]byte, error) {
	fullURL := c.baseURL + endpoint

	var lastErr error
	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		if attempt > 0 {
			c.logger.Debug("retrying request",
				zap.String("method", method),
				zap.String("url", fullURL),
				zap.Int("attempt", attempt),
			)
			time.Sleep(c.retryDelay)
		}

		// Create request
		req, err := http.NewRequestWithContext(ctx, method, fullURL, body)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create request")
		}

		// Add Basic Auth
		req.SetBasicAuth(c.username, c.password)

		// Add headers
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		// Add CSRF crumb if provided
		if crumb != nil {
			req.Header.Set(crumb.CrumbRequestField, crumb.Crumb)
			c.logger.Debug("added crumb to request", zap.String("field", crumb.CrumbRequestField))
		}

		// Send request
		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = errors.Wrap(err, "HTTP request failed")
			continue
		}

		// Read response body
		defer resp.Body.Close()
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			lastErr = errors.Wrap(err, "failed to read response body")
			continue
		}

		// Check status code
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return respBody, nil
		}

		// Handle error status codes
		switch resp.StatusCode {
		case 401, 403:
			return nil, errors.Errorf("authentication failed: %s (status: %d)", string(respBody), resp.StatusCode)
		case 404:
			return nil, errors.Errorf("resource not found: %s (status: %d)", endpoint, resp.StatusCode)
		default:
			lastErr = errors.Errorf("request failed with status %d: %s", resp.StatusCode, string(respBody))
			// Retry on 5xx errors
			if resp.StatusCode >= 500 {
				continue
			}
			return nil, lastErr
		}
	}

	return nil, errors.Wrapf(lastErr, "request failed after %d retries", c.maxRetries)
}

// Ping checks if Jenkins is accessible
func (c *Client) Ping(ctx context.Context) error {
	endpoint := "/api/json"

	if err := c.doRequest(ctx, "GET", endpoint, nil, nil); err != nil {
		return errors.Wrap(err, "failed to ping Jenkins")
	}

	return nil
}
