package config

// Def* constants are default values for config.
const (
	// DefEnv is default environment.
	DefEnv = "prod"

	// DefPort is default http rest server port.
	DefPort = "8080"

	// DefServerAddress is default http rest server address.
	DefServerAddress = "localhost:" + DefPort

	// DefBaseURL is default base url address for short urls.
	DefBaseURL = "http://" + DefServerAddress + "/"

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
