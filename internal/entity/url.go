package entity

// URL a url entity.
type URL struct {
	ID          string `json:"id,omitempty" db:"id"`
	UserID      string `json:"user_id" db:"user_id"`
	ShortURL    string `json:"short_url" db:"short_url"`
	OriginalURL string `json:"original_url" db:"original_url"`
	DeletedFlag bool   `json:"is_deleted" db:"is_deleted"`
}
