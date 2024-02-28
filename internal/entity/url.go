package entity

// URL a url entity.
type URL struct {
	// ID is the unique identifier of the URL.
	ID string `json:"id,omitempty" db:"id"`

	// UserID is the unique identifier of the user.
	UserID string `json:"user_id" db:"user_id"`

	// ShortURL is the short version of the URL.
	ShortURL string `json:"short_url" db:"short_url"`

	// OriginalURL is the original URL.
	OriginalURL string `json:"original_url" db:"original_url"`

	// DeletedFlag is the flag to indicate if the URL is deleted.
	DeletedFlag bool `json:"is_deleted" db:"is_deleted"`
}

// URLStats represents the stats of the URL.
type URLStats struct {
	// URLs is the number of unique URLs.
	URLs string `json:"urls"`

	// Users is the number of unique users.
	Users string `json:"users"`
}
