// Package main entry point of the application
package main

import (
	"log"
)

func main() {
	if err := run(); err != nil {
		log.Panic(err.Error())
	}
}
