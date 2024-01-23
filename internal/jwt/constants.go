package jwt

type key string

var (
	CookieKey      = "AccessToken"
	ContextKey key = "access-token"
)
