package api

type CreateURLRequest struct {
	URL string `json:"url"`
}

type CreateURLResponse struct {
	Result string `json:"result"`
}
