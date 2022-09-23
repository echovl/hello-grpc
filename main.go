package main

import (
	"log"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	server := NewServer(":8080")
	if err := server.ListenAndServe(); err != nil {
		return err
	}
	return nil
}
