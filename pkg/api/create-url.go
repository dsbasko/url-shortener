package api

// CreateURLRequest the request body for the create url endpoints.
type CreateURLRequest struct {
	URL string `json:"url"`
}

// CreateURLResponse the response body for the create url endpoints.
type CreateURLResponse struct {
	Result string `json:"result"`
}
