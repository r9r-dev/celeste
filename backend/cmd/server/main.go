package main

import (
	"log"
	"os"

	"aperture-science-network/internal/api"
	"aperture-science-network/internal/docker"
	"aperture-science-network/internal/mock"
	"aperture-science-network/internal/stack"
	"aperture-science-network/internal/stats"
	"aperture-science-network/internal/version"
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

	staticPath := os.Getenv("STATIC_PATH")
	if staticPath == "" {
		staticPath = "./static"
	}

	debugMode := os.Getenv("DEBUG_MODE") == "true"

	var dockerClient docker.DockerClient
	var statsProvider stats.Provider
	var stackProvider stack.Provider

	if debugMode {
		log.Println("[DEBUG MODE] Using mock data providers")
		dockerClient = mock.NewDockerClient()
		statsProvider = mock.NewStatsProvider()
		stackProvider = mock.NewStackProvider()
	} else {
		var err error
		dockerClient, err = docker.NewClient()
		if err != nil {
			log.Fatalf("Failed to create Docker client: %v", err)
		}
		statsProvider = stats.NewDefaultProvider()
		stackProvider = stack.NewFilesystemProvider(stacksPath, dockerClient)
	}

	server := api.NewServer(api.ServerOptions{
		StaticPath:    staticPath,
		DockerClient:  dockerClient,
		StatsProvider: statsProvider,
		StackProvider: stackProvider,
	})

	log.Printf("Aperture Science Network v%s starting on port %s", version.Version, port)
	log.Printf("Stacks path: %s", stacksPath)
	if debugMode {
		log.Println("[DEBUG MODE] Mock data active - Docker not required")
	}

	if err := server.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
