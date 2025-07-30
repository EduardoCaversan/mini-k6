package results

import "time"

type TestResult struct {
	Method     string        `json:"method"`
	URL        string        `json:"url"`
	StatusCode int           `json:"status_code"`
	Error      string        `json:"error,omitempty"`
	Duration   time.Duration `json:"duration"`
}