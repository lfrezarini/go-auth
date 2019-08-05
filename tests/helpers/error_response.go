package tests

// ErrorResponse represents an common error structure returned by api
type ErrorResponse []struct {
	Message    string   `json:"message"`
	Path       []string `json:"path"`
	Extensions struct {
		Code string `json:"code"`
	} `json:"extensions"`
}
