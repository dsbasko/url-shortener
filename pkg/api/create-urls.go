package api

// CreateURLsRequest the request body for the create urls endpoints.
type CreateURLsRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// CreateURLsResponse the response body for the create urls endpoints.
type CreateURLsResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}
