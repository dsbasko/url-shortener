package config

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

// ParseServerAddress parses server address.
//
// If serverAddress does not contain "://", it will be prefixed with "https://".
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
//
// If baseURL does not contain "://", it will be prefixed with "https://".
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
