package main

import (
	"log"
	"os"

	"aperture-science-network/internal/api"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	stacksPath := os.Getenv("STACKS_PATH")
	if stacksPath == "" {
		stacksPath = "/home/share/docker/dockge/stacks"
	}

	server := api.NewServer(stacksPath)

	log.Printf("Aperture Science Network starting on port %s", port)
	log.Printf("Stacks path: %s", stacksPath)

	if err := server.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
