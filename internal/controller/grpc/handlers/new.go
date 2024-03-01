package handlers

import (
	pb "github.com/dsbasko/url-shortener/api/proto"
	"github.com/dsbasko/url-shortener/internal/service/urls"
	"github.com/dsbasko/url-shortener/pkg/logger"
)

// URLShortenerServer is the server that provides the URLShortener service.
type URLShortenerServer struct {
	pb.UnimplementedURLShortenerV1Server

	log        *logger.Logger
	pinger     Pinger
	urlService urls.URLs
}

// New creates a new URLShortenerServer
func New(log *logger.Logger, pinger Pinger, urlService urls.URLs) URLShortenerServer {
	return URLShortenerServer{
		log:        log,
		pinger:     pinger,
		urlService: urlService,
	}
}
