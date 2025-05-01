package types

type APISuccessStatus struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

type APIErrorStatus struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}
