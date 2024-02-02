// Package main entry point of the application
package main

import (
	"fmt"
	"log"

	"github.com/dsbasko/yandex-go-shortener/internal/app"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func printGreeting() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n\n", buildCommit)
}

func main() {
	printGreeting()

	if err := app.Run(); err != nil {
		log.Panic(err.Error())
	}
}
