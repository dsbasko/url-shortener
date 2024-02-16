package config

// Def* constants are default values for config.
const (
	// DefEnv is default environment.
	DefEnv = "prod"

	// DefServerAddress is default http rest server address.
	DefServerAddress = "localhost:8080"

	// DefBaseURL is default base url address for short urls.
	DefBaseURL = DefServerAddress

	// DefShortURLLen is default length of short url.
	DefShortURLLen = 4

	// DefRESTReadTimeout is default http rest server read timeout.
	DefRESTReadTimeout = 3000

	// DefRESTWriteTimeout is default http rest server write timeout.
	DefRESTWriteTimeout = 3000

	// DefDatabaseMaxConns is default max connections to postgresql.
	DefDatabaseMaxConns = 8

	// DefJWTSecret is default secret for JWT auth.
	DefJWTSecret = "allons-y"
)
