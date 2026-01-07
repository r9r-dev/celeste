package api

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"aperture-science-network/internal/api/handlers"
	"aperture-science-network/internal/compose"
	"aperture-science-network/internal/docker"
	"aperture-science-network/internal/stack"
	"aperture-science-network/internal/stats"
	"aperture-science-network/internal/version"
	"aperture-science-network/internal/ws"
)

type Server struct {
	router         *gin.Engine
	dockerClient   docker.DockerClient
	statsProvider  stats.Provider
	stackProvider  stack.Provider
	composeManager *compose.Manager
	wsHub          *ws.Hub
	staticPath     string
}

// ServerOptions contains the dependencies for creating a new server
type ServerOptions struct {
	StaticPath    string
	DockerClient  docker.DockerClient
	StatsProvider stats.Provider
	StackProvider stack.Provider
}

func NewServer(opts ServerOptions) *Server {
	gin.SetMode(gin.ReleaseMode)

	wsHub := ws.NewHub(opts.DockerClient, opts.StatsProvider)
	go wsHub.Run()

	composeManager := compose.NewManager()

	s := &Server{
		router:         gin.New(),
		dockerClient:   opts.DockerClient,
		statsProvider:  opts.StatsProvider,
		stackProvider:  opts.StackProvider,
		composeManager: composeManager,
		wsHub:          wsHub,
		staticPath:     opts.StaticPath,
	}

	s.setupMiddleware()
	s.setupRoutes()

	return s
}

func (s *Server) setupMiddleware() {
	// Custom logger that only logs errors (4xx and 5xx responses)
	s.router.Use(func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		c.Next()

		// Only log client/server errors (4xx, 5xx)
		// Skip 2xx (success) and 3xx (redirects, cache)
		status := c.Writer.Status()
		if status >= 400 {
			latency := time.Since(start)
			log.Printf("[ERROR] %d | %v | %s %s | %s",
				status,
				latency,
				c.Request.Method,
				path,
				c.Errors.String(),
			)
		}
	})
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
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"version": version.Version,
		})
	})

	// API routes
	api := s.router.Group("/api")
	{
		// System stats
		api.GET("/stats", handlers.GetSystemStats(s.statsProvider))

		// Stacks
		stacks := api.Group("/stacks")
		{
			stacks.GET("", handlers.ListStacks(s.stackProvider))
			stacks.GET("/:name", handlers.GetStack(s.stackProvider))
			stacks.POST("/:name/start", handlers.StartStack(s.stackProvider, s.composeManager))
			stacks.POST("/:name/stop", handlers.StopStack(s.stackProvider, s.composeManager))
			stacks.POST("/:name/restart", handlers.RestartStack(s.stackProvider, s.composeManager))
			stacks.POST("/:name/pull", handlers.PullStack(s.stackProvider, s.composeManager))
			stacks.GET("/:name/compose", handlers.GetComposeFile(s.stackProvider))
			stacks.PUT("/:name/compose", handlers.UpdateComposeFile(s.stackProvider))
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

	// Serve static files (SvelteKit build output)
	s.router.Static("/_app", s.staticPath+"/_app")
	s.router.StaticFile("/favicon.ico", s.staticPath+"/favicon.ico")
	s.router.StaticFile("/robots.txt", s.staticPath+"/robots.txt")

	// SPA fallback - serve index.html for all unmatched routes
	s.router.NoRoute(func(c *gin.Context) {
		c.File(s.staticPath + "/index.html")
	})
}

func (s *Server) Run(addr string) error {
	return s.router.Run(addr)
}
