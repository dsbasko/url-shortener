package config

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

func ParseServerAddress(serverAddress string) (string, error) {
	if !strings.Contains(serverAddress, "://") {
		serverAddress = "https://" + serverAddress
	}

	parsedURL, err := url.ParseRequestURI(serverAddress)
	if err != nil {
		return "", fmt.Errorf("url.ParseRequestURI: %w", err)
	}

	return parsedURL.Host, nil
}

func ParseBaseURL(baseURL string) (string, error) {
	if baseURL[0] == ':' {
		return "", errors.New("need host before the port")
	}

	if !strings.Contains(baseURL, "://") {
		baseURL = "https://" + baseURL
	}

	_, err := url.ParseRequestURI(baseURL)
	if err != nil {
		return "", fmt.Errorf("url.ParseRequestURI: %w", err)
	}

	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}

	return baseURL, nil
}
