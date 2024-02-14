package config

import (
	"fmt"
	"sync"
	"time"

	"github.com/ian-kent/gofigure"
)

type config struct {
	gofigure         interface{} `order:"flag,env"`                                                                                                        //nolint:lll,gofmt,unused
	Env              string      `env:"ENV" flag:"env" flagDesc:"run mode (prod|dev|local)"`                                                               //nolint:lll
	ServerAddress    string      `env:"SERVER_ADDRESS" flag:"a" flagDesc:"http rest server address"`                                                       //nolint:lll
	BaseURL          string      `env:"BASE_URL" flag:"b" flagDesc:"base url address"`                                                                     //nolint:lll
	ShortURLLen      int         `env:"SHORT_URL_LEN" flag:"short-url-len" flagDesc:"short url length"`                                                    //nolint:lll
	StoragePath      string      `env:"FILE_STORAGE_PATH" flag:"f" flagDesc:"full path of the json repositories file"`                                     //nolint:lll
	RESTReadTimeout  int         `env:"REST_READ_TIMEOUT" flag:"rest-read-timeout" flagDesc:"wait timeout for reading request on the http rest server"`    //nolint:lll
	RESTWriteTimeout int         `env:"REST_WRITE_TIMEOUT" flag:"rest-write-timeout" flagDesc:"wait timeout for writing response on the http rest server"` //nolint:lll
	RESTEnableHTTPS  bool        `env:"ENABLE_HTTPS" flag:"s" flagDesc:"enable https for rest server"`                                                     //nolint:lll
	PsqlDSN          string      `env:"DATABASE_DSN" flag:"d" flagDesc:"string for connecting to database"`                                                //nolint:lll
	PsqlMaxConns     int         `env:"PSQL_MAX_CONNS" flag:"psql-max-conns" flagDesc:"max connections to database"`                                       //nolint:lll
	JWTSecret        string      `env:"JWT_SECRET" flag:"jwt" flagDesc:"jwt secret"`
}

var (
	cfg  config
	once sync.Once
	err  error
)

// Init singleton config initialization.
func Init() error {
	once.Do(func() {
		cfg = config{
			Env:              DefEnv,
			ServerAddress:    DefServerAddress,
			BaseURL:          DefBaseURL,
			ShortURLLen:      DefShortURLLen,
			RESTReadTimeout:  DefRESTReadTimeout,
			RESTWriteTimeout: DefRESTWriteTimeout,
			PsqlMaxConns:     DefPsqlMaxConns,
			JWTSecret:        DefJWTSecret,
		}

		err = gofigure.Gofigure(&cfg)
		if err != nil {
			err = fmt.Errorf("gofigure.Gofigure: %w", err)
			return
		}

		serverAddress, errParser := ParseServerAddress(cfg.ServerAddress)
		if errParser != nil {
			err = fmt.Errorf("parseServerAddress: %w", errParser)
			return
		}
		cfg.ServerAddress = serverAddress

		baseURL, errParser := ParseBaseURL(cfg.BaseURL)
		if errParser != nil {
			err = fmt.Errorf("parseBaseURL: %w", errParser)
			return
		}
		cfg.BaseURL = baseURL
	})

	return err
}

// Env returns run mode (dev|prod).
func Env() string {
	return cfg.Env
}

// ServerAddress returns http rest server address.
func ServerAddress() string {
	return cfg.ServerAddress
}

// BaseURL returns base url address.
func BaseURL() string {
	return cfg.BaseURL
}

// ShortURLLen returns short url length.
func ShortURLLen() int {
	return cfg.ShortURLLen
}

// StoragePath returns full path of the json repositories file.
func StoragePath() string {
	return cfg.StoragePath
}

// PsqlDSN returns string for connecting to database.
func PsqlDSN() string {
	return cfg.PsqlDSN
}

// RESTReadTimeout returns wait timeout for reading request on the http rest server.
func RESTReadTimeout() time.Duration {
	return time.Duration(cfg.RESTReadTimeout) * time.Millisecond
}

// RESTWriteTimeout returns wait timeout for writing response on the http rest server.
func RESTWriteTimeout() time.Duration {
	return time.Duration(cfg.RESTWriteTimeout) * time.Millisecond
}

// RESTEnableHTTPS returns true if https is enabled for rest server.
func RESTEnableHTTPS() bool {
	return cfg.RESTEnableHTTPS
}

// PsqlMaxConns returns max connections to database.
func PsqlMaxConns() int {
	return cfg.PsqlMaxConns
}

// JWTSecret returns jwt secret.
func JWTSecret() []byte {
	return []byte(cfg.JWTSecret)
}
