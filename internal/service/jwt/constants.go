package jwt

type key string

// Constants.
var (
	// CookieKey a key for jwt token in cookie.
	CookieKey = "AccessToken"

	// ContextString a key for jwt token in context.
	ContextString string = "access-token"

	// ContextKey a key for jwt token in context.
	ContextKey key = "access-token"
)
