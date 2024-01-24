// Package main entry point of the application
package main

import (
	"log"

	"github.com/dsbasko/yandex-go-shortener/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Panic(err.Error())
	}
}
