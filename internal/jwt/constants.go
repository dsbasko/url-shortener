package jwt

type key string

// Constants.
var (
	// CookieKey a key for jwt token in cookie.
	CookieKey = "AccessToken"

	// ContextKey a key for jwt token in context.
	ContextKey key = "access-token"
)
