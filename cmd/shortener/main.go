// Package main entry point of the application
package main

import (
	"log"

	"github.com/dsbasko/url-shortener/internal/app"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	if err := app.RunURLShortener(buildVersion, buildDate, buildCommit); err != nil {
		log.Panic(err.Error())
	}
}
