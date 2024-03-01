package config

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	gocfg "github.com/dsbasko/go-cfg"
)

type config struct {
	Env                 string `env:"ENV" flag:"env" json:"env" description:"run mode (prod|dev|local)"`                                                                              //nolint:lll
	Controller          string `env:"CONTROLLER" flag:"controller" json:"controller" description:"controller (grpc|http)"`                                                            //nolint:lll
	ConfigPath          string `env:"CONFIG" s-flag:"c" flag:"config" description:"path to config file (json,yaml)"`                                                                  //nolint:lll
	ServerAddress       string `env:"SERVER_ADDRESS" s-flag:"a" flag:"server-address" json:"server_address" description:"http rest server address"`                                   //nolint:lll
	BaseURL             string `env:"BASE_URL" s-flag:"b" flag:"base-url" json:"base_url" description:"base url address"`                                                             //nolint:lll
	ShortURLLen         int    `env:"SHORT_URL_LEN" flag:"short-url-len" json:"short_url_len" description:"short url length"`                                                         //nolint:lll
	StoragePath         string `env:"FILE_STORAGE_PATH" s-flag:"f" flag:"file-storage-path" json:"file_storage_path" description:"full path of the json repositories file"`           //nolint:lll
	RESTReadTimeout     int    `env:"REST_READ_TIMEOUT" flag:"rest-read-timeout" json:"rest_read_timeout" description:"wait timeout for reading request on the http rest server"`     //nolint:lll
	RESTWriteTimeout    int    `env:"REST_WRITE_TIMEOUT" flag:"rest-write-timeout" json:"rest_write_timeout" description:"wait timeout for writing response on the http rest server"` //nolint:lll
	IsRESTEnableHTTPS   bool   `env:"ENABLE_HTTPS" s-flag:"s" flag:"enable-https" json:"enable_https" description:"enable https for rest server"`                                     //nolint:lll
	DatabaseDSN         string `env:"DATABASE_DSN" s-flag:"d" flag:"database-dsn" json:"database_dsn" description:"string for connecting to database"`                                //nolint:lll
	DatabaseMaxConnects int    `env:"DATABASE_MAX_CONNECTIONS" json:"database_max_connections" flag:"database-max-connections" description:"max connections to database"`             //nolint:lll
	JWTSecret           string `env:"JWT_SECRET" flag:"jwt-secret" json:"jwt_secret" description:"jwt secret"`                                                                        //nolint:lll
	TrustedSubnet       string `env:"TRUSTED_SUBNET" s-flag:"t" flag:"trusted-subnet" json:"trusted_subnet" description:"jwt secret"`                                                 //nolint:lll
	IsEnabledPPROF      bool   `env:"PPROF" flag:"pprof" json:"pprof" description:"enable pprof for rest server"`                                                                     //nolint:lll
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
			Env:                 DefEnv,
			Controller:          DefController,
			ServerAddress:       DefServerAddress,
			BaseURL:             DefBaseURL,
			ShortURLLen:         DefShortURLLen,
			RESTReadTimeout:     DefRESTReadTimeout,
			RESTWriteTimeout:    DefRESTWriteTimeout,
			DatabaseMaxConnects: DefDatabaseMaxConns,
			JWTSecret:           DefJWTSecret,
			TrustedSubnet:       DefTrustedSubnet,
		}

		// Use my own library to read the configuration
		// github.com/dsbasko/go-cfg
		cfgPath := CfgPath()
		if cfgPath != "" {
			gocfg.MustReadFile(cfgPath, &cfg)
		}

		gocfg.MustReadEnv(&cfg)
		gocfg.MustReadFlag(&cfg)
	})

	return err
}

// MustInit singleton config initialization with panic.
func MustInit() {
	if err = Init(); err != nil {
		panic(fmt.Errorf("the configuration could not be loaded: %w", err))
	}
}

// InitMock initializes the configuration with default values.
func InitMock() {
	cfg = config{
		Env:                 DefEnv,
		Controller:          DefController,
		ServerAddress:       DefServerAddress,
		BaseURL:             DefBaseURL,
		ShortURLLen:         DefShortURLLen,
		RESTReadTimeout:     DefRESTReadTimeout,
		RESTWriteTimeout:    DefRESTWriteTimeout,
		DatabaseMaxConnects: DefDatabaseMaxConns,
		JWTSecret:           DefJWTSecret,
		TrustedSubnet:       DefTrustedSubnet,
	}
}

// Env returns run mode (dev|prod).
func Env() string {
	return cfg.Env
}

// Controller returns controller (grpc|http).
func Controller() string {
	return cfg.Controller
}

// CfgPath returns path to config file.
func CfgPath() string {
	pathCfg := config{}
	gocfg.MustReadEnv(&pathCfg)
	gocfg.MustReadFlag(&pathCfg)
	return pathCfg.ConfigPath
}

// ServerAddress returns http rest server address.
func ServerAddress() string {
	return cfg.ServerAddress
}

// BaseURL returns base url address.
func BaseURL() string {
	if strings.HasPrefix(cfg.BaseURL, "http") {
		if strings.HasSuffix(cfg.BaseURL, "/") {
			return cfg.BaseURL
		}

		return cfg.BaseURL + "/"
	}

	if cfg.IsRESTEnableHTTPS {
		return fmt.Sprintf("https://%s/", cfg.BaseURL)
	}

	return fmt.Sprintf("http://%s/", cfg.BaseURL)
}

// ShortURLLen returns short url length.
func ShortURLLen() int {
	return cfg.ShortURLLen
}

// StoragePath returns full path of the json repositories file.
func StoragePath() string {
	return cfg.StoragePath
}

// DatabaseDSN returns string for connecting to database.
func DatabaseDSN() string {
	return cfg.DatabaseDSN
}

// DatabaseMaxConnects returns max connections to database.
func DatabaseMaxConnects() int {
	return cfg.DatabaseMaxConnects
}

// RESTReadTimeout returns wait timeout for reading request on the http rest server.
func RESTReadTimeout() time.Duration {
	return time.Duration(cfg.RESTReadTimeout) * time.Millisecond
}

// RESTWriteTimeout returns wait timeout for writing response on the http rest server.
func RESTWriteTimeout() time.Duration {
	return time.Duration(cfg.RESTWriteTimeout) * time.Millisecond
}

// IsRESTEnableHTTPS returns true if https is enabled for rest server.
func IsRESTEnableHTTPS() bool {
	return cfg.IsRESTEnableHTTPS
}

// JWTSecret returns jwt secret.
func JWTSecret() []byte {
	return []byte(cfg.JWTSecret)
}

// IsTrustedSubnet checks if the given IP address belongs to a trusted subnet.
// It takes an IP address and port in the format "ip:port" and returns a boolean
// indicating whether the IP address is within the trusted subnet or not. If an
// error occurs during the process, it returns an error.
func IsTrustedSubnet(ipAndPort string) (bool, error) {
	ip, _, errSHP := net.SplitHostPort(ipAndPort)
	if errSHP != nil {
		return false, fmt.Errorf("failed to split host and port: %w", errSHP)
	}

	trustedSubnet := cfg.TrustedSubnet
	if !strings.Contains(trustedSubnet, "/") {
		trustedSubnet = fmt.Sprintf("%s/32", trustedSubnet)
	}

	_, ipNet, errCIDR := net.ParseCIDR(trustedSubnet)
	if errCIDR != nil {
		return false, fmt.Errorf("failed to parse trusted subnet: %w", errCIDR)
	}

	return ipNet.Contains(net.ParseIP(ip)), nil
}

// SetTrustedSubnet sets the trusted subnet for the configuration.
//
// ⚠️ This function is used for mocking purposes!
func SetTrustedSubnet(trustedSubnet string) {
	cfg.TrustedSubnet = trustedSubnet
}

// IsEnabledPPROF returns true if pprof is enabled for rest server.
func IsEnabledPPROF() bool {
	return cfg.IsEnabledPPROF
}
