package models

type TestScenario struct {
	ConcurrentUsers int          `json:"concurrent_users"`
	DurationSeconds int          `json:"duration_seconds"`
	Requests        []APIRequest `json:"requests"`
	MaxRequests     int          `json:"max_requests,omitempty"`
}

type APIRequest struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers,omitempty"`
	Body    any               `json:"body,omitempty"`
}
