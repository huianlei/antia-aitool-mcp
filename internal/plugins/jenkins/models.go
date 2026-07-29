package jenkins

import (
	"time"
)

// Job represents a Jenkins job
type Job struct {
	Name        string `json:"name"`
	URL         string `json:"url"`
	Color       string `json:"color"` // blue, red, yellow, etc.
	Description string `json:"description,omitempty"`
}

// JobDetail represents detailed information about a Jenkins job
type JobDetail struct {
	Name            string  `json:"name"`
	URL             string  `json:"url"`
	Description     string  `json:"description"`
	Buildable       bool    `json:"buildable"`
	Color           string  `json:"color"`
	LastBuild       *Build  `json:"lastBuild,omitempty"`
	LastSuccessBuild *Build `json:"lastSuccessfulBuild,omitempty"`
	LastFailedBuild *Build  `json:"lastFailedBuild,omitempty"`
	HealthReport    []struct {
		Description string `json:"description"`
		Score       int    `json:"score"`
	} `json:"healthReport"`
}

// Build represents a Jenkins build
type Build struct {
	Number    int    `json:"number"`
	URL       string `json:"url"`
	Result    string `json:"result,omitempty"` // SUCCESS, FAILURE, UNSTABLE, ABORTED, null (in progress)
	Building  bool   `json:"building"`
	Duration  int64  `json:"duration"` // milliseconds
	Timestamp int64  `json:"timestamp"` // Unix timestamp in milliseconds
}

// BuildDetail represents detailed information about a build
type BuildDetail struct {
	Number      int                    `json:"number"`
	URL         string                 `json:"url"`
	Result      string                 `json:"result,omitempty"`
	Building    bool                   `json:"building"`
	Duration    int64                  `json:"duration"`
	Timestamp   int64                  `json:"timestamp"`
	Description string                 `json:"description,omitempty"`
	FullDisplayName string             `json:"fullDisplayName"`
	Actions     []map[string]interface{} `json:"actions,omitempty"`
	ChangeSet   struct {
		Items []struct {
			Author struct {
				FullName string `json:"fullName"`
			} `json:"author"`
			Comment string `json:"comment"`
			Date    string `json:"date"`
		} `json:"items"`
	} `json:"changeSet"`
}

// JobsResponse represents the response from /api/json
type JobsResponse struct {
	Jobs []Job `json:"jobs"`
}

// BuildsResponse represents the response for build history
type BuildsResponse struct {
	Builds []Build `json:"builds"`
}

// QueueItem represents a queued build
type QueueItem struct {
	ID          int    `json:"id"`
	URL         string `json:"url"`
	Why         string `json:"why"`
	Blocked     bool   `json:"blocked"`
	Buildable   bool   `json:"buildable"`
	BuildableStartMilliseconds int64 `json:"buildableStartMilliseconds"`
}

// Config represents Jenkins configuration
type Config struct {
	URL         string
	Username    string
	Password    string
	Timeout     time.Duration
	VerifySSL   bool
	MaxRetries  int
	RetryDelay  time.Duration
}

// Crumb represents Jenkins CSRF protection token
type Crumb struct {
	Crumb             string `json:"crumb"`
	CrumbRequestField string `json:"crumbRequestField"`
}
