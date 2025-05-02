package main

import (
	"log"
	"os"
)

func exit(err error) {
	if err != nil {
		log.Printf("server error: %v", err)
		os.Exit(1)
	}
	os.Exit(0)
}
