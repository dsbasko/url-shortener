package config

// Def* constants are default values for config.
const (
	DefEnv              = "prod"
	DefPort             = "8080"
	DefServerAddress    = "localhost:" + DefPort
	DefBaseURL          = "http://" + DefServerAddress + "/"
	DefShortURLLen      = 4
	DefRestReadTimeout  = 3000
	DefRestWriteTimeout = 3000
	DefPsqlMaxConns     = 8
	DefJWTSecret        = "allons-y"
)
