package config

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

// ParseServerAddress parses server address.
func ParseServerAddress(serverAddress string) (string, error) {
	if !strings.Contains(serverAddress, "://") {
		serverAddress = "https://" + serverAddress
	}

	parsedURL, errParse := url.ParseRequestURI(serverAddress)
	if errParse != nil {
		return "", fmt.Errorf("url.ParseRequestURI: %w", errParse)
	}

	return parsedURL.Host, nil
}

// ParseBaseURL parses base url.
func ParseBaseURL(baseURL string) (string, error) {
	if baseURL[0] == ':' {
		return "", errors.New("need host before the port")
	}

	if !strings.Contains(baseURL, "://") {
		baseURL = "https://" + baseURL
	}

	_, errParse := url.ParseRequestURI(baseURL)
	if errParse != nil {
		return "", fmt.Errorf("url.ParseRequestURI: %w", errParse)
	}

	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}

	return baseURL, nil
}
