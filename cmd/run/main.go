package main

import (
	"log"

	_ "github.com/joho/godotenv/autoload"
)

func main() {

	// Try to setup bot and its dependencies
	b, err := bootstrap()
	if err != nil {
		log.Fatalf("failed to boostrap: %v", err)
	}

	// Run the bot. Blocks until error occurs
	err = b.run()
	if err != nil {
		log.Fatalf("app exited: %v", err)
	}
}
