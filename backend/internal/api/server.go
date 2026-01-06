package api

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"aperture-science-network/internal/api/handlers"
	"aperture-science-network/internal/compose"
	"aperture-science-network/internal/docker"
	"aperture-science-network/internal/ws"
)

type Server struct {
	router         *gin.Engine
	dockerClient   *docker.Client
	composeManager *compose.Manager
	wsHub          *ws.Hub
	stacksPath     string
}

func NewServer(stacksPath string) *Server {
	gin.SetMode(gin.ReleaseMode)

	dockerClient, err := docker.NewClient()
	if err != nil {
		panic(err)
	}

	wsHub := ws.NewHub(dockerClient)
	go wsHub.Run()

	composeManager := compose.NewManager()

	s := &Server{
		router:         gin.New(),
		dockerClient:   dockerClient,
		composeManager: composeManager,
		wsHub:          wsHub,
		stacksPath:     stacksPath,
	}

	s.setupMiddleware()
	s.setupRoutes()

	return s
}

func (s *Server) setupMiddleware() {
	s.router.Use(gin.Logger())
	s.router.Use(gin.Recovery())

	s.router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3001", "http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
}

func (s *Server) setupRoutes() {
	// Health check
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// API routes
	api := s.router.Group("/api")
	{
		// System stats
		api.GET("/stats", handlers.GetSystemStats)

		// Stacks
		stacks := api.Group("/stacks")
		{
			stacks.GET("", handlers.ListStacks(s.stacksPath, s.dockerClient))
			stacks.GET("/:name", handlers.GetStack(s.stacksPath, s.dockerClient))
			stacks.POST("/:name/start", handlers.StartStack(s.stacksPath, s.composeManager))
			stacks.POST("/:name/stop", handlers.StopStack(s.stacksPath, s.composeManager))
			stacks.POST("/:name/restart", handlers.RestartStack(s.stacksPath, s.composeManager))
			stacks.POST("/:name/pull", handlers.PullStack(s.stacksPath, s.composeManager))
			stacks.GET("/:name/compose", handlers.GetComposeFile(s.stacksPath))
			stacks.PUT("/:name/compose", handlers.UpdateComposeFile(s.stacksPath))
		}

		// Containers
		containers := api.Group("/containers")
		{
			containers.GET("", handlers.ListContainers(s.dockerClient))
			containers.GET("/:id", handlers.GetContainer(s.dockerClient))
			containers.POST("/:id/start", handlers.StartContainer(s.dockerClient))
			containers.POST("/:id/stop", handlers.StopContainer(s.dockerClient))
			containers.POST("/:id/restart", handlers.RestartContainer(s.dockerClient))
			containers.GET("/:id/logs", handlers.GetContainerLogs(s.dockerClient))
			containers.GET("/:id/stats", handlers.GetContainerStats(s.dockerClient))
		}

		// Volumes
		volumes := api.Group("/volumes")
		{
			volumes.GET("", handlers.ListVolumes(s.dockerClient))
			volumes.POST("", handlers.CreateVolume(s.dockerClient))
			volumes.DELETE("/:name", handlers.DeleteVolume(s.dockerClient))
		}

		// Networks
		networks := api.Group("/networks")
		{
			networks.GET("", handlers.ListNetworks(s.dockerClient))
			networks.POST("", handlers.CreateNetwork(s.dockerClient))
			networks.DELETE("/:id", handlers.DeleteNetwork(s.dockerClient))
		}

		// Images
		api.GET("/images", handlers.ListImages(s.dockerClient))
	}

	// WebSocket
	s.router.GET("/ws", func(c *gin.Context) {
		ws.HandleWebSocket(s.wsHub, c.Writer, c.Request)
	})
}

func (s *Server) Run(addr string) error {
	return s.router.Run(addr)
}
