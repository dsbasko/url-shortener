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
	RestReadTimeout  int         `env:"REST_READ_TIMEOUT" flag:"rest-read-timeout" flagDesc:"wait timeout for reading request on the http rest server"`    //nolint:lll
	RestWriteTimeout int         `env:"REST_WRITE_TIMEOUT" flag:"rest-write-timeout" flagDesc:"wait timeout for writing response on the http rest server"` //nolint:lll
	PsqlDSN          string      `env:"DATABASE_DSN" flag:"d" flagDesc:"string for connecting to database"`                                                //nolint:lll
	PsqlMaxConns     int         `env:"PSQL_MAX_CONNS" flag:"psql-max-conns" flagDesc:"max connections to database"`                                       //nolint:lll
	JWTSecret        string      `env:"JWT_SECRET" flag:"jwt" flagDesc:"jwt secret"`
}

var (
	cfg  config
	once sync.Once
	err  error
)

func Init() error {
	once.Do(func() {
		cfg = config{
			Env:              DefEnv,
			ServerAddress:    DefServerAddress,
			BaseURL:          DefBaseURL,
			ShortURLLen:      DefShortURLLen,
			RestReadTimeout:  DefRestReadTimeout,
			RestWriteTimeout: DefRestWriteTimeout,
			PsqlMaxConns:     DefPsqlMaxConns,
			JWTSecret:        DefJWTSecret,
		}

		errGofigure := gofigure.Gofigure(&cfg)
		if errGofigure != nil {
			err = fmt.Errorf("gofigure.Gofigure: %w", errGofigure)
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

func GetEnv() string {
	return cfg.Env
}

func GetServerAddress() string {
	return cfg.ServerAddress
}

func GetBaseURL() string {
	return cfg.BaseURL
}

func GetShortURLLen() int {
	return cfg.ShortURLLen
}

func GetStoragePath() string {
	return cfg.StoragePath
}

func GetPsqlDSN() string {
	return cfg.PsqlDSN
}

func GetRestReadTimeout() time.Duration {
	return time.Duration(cfg.RestReadTimeout) * time.Millisecond
}

func GetRestWriteTimeout() time.Duration {
	return time.Duration(cfg.RestWriteTimeout) * time.Millisecond
}

func GetPsqlMaxConns() int {
	return cfg.PsqlMaxConns
}

func GetJWTSecret() []byte {
	return []byte(cfg.JWTSecret)
}