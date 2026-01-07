package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"aperture-science-network/internal/compose"
	"aperture-science-network/internal/docker"
	"aperture-science-network/internal/stack"
	"aperture-science-network/internal/stats"
)

// System Stats
func GetSystemStats(statsProvider stats.Provider) gin.HandlerFunc {
	return func(c *gin.Context) {
		sysStats, err := statsProvider.GetSystemStats()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, sysStats)
	}
}

// Stacks
func ListStacks(stackProvider stack.Provider) gin.HandlerFunc {
	return func(c *gin.Context) {
		stacks, err := stackProvider.ListStacks()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, stacks)
	}
}

func GetStack(stackProvider stack.Provider) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")
		s, err := stackProvider.GetStack(name)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Stack not found"})
			return
		}
		c.JSON(http.StatusOK, s)
	}
}

func StartStack(stackProvider stack.Provider, composeManager *compose.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")

		if !stackProvider.StackExists(name) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Stack not found"})
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Minute)
		defer cancel()

		stackPath := stackProvider.GetStackPath(name)
		if err := composeManager.Up(ctx, stackPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "started"})
	}
}

func StopStack(stackProvider stack.Provider, composeManager *compose.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")

		if !stackProvider.StackExists(name) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Stack not found"})
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Minute)
		defer cancel()

		stackPath := stackProvider.GetStackPath(name)
		if err := composeManager.Down(ctx, stackPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "stopped"})
	}
}

func RestartStack(stackProvider stack.Provider, composeManager *compose.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")

		if !stackProvider.StackExists(name) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Stack not found"})
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Minute)
		defer cancel()

		stackPath := stackProvider.GetStackPath(name)
		if err := composeManager.Restart(ctx, stackPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "restarted"})
	}
}

func PullStack(stackProvider stack.Provider, composeManager *compose.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")

		if !stackProvider.StackExists(name) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Stack not found"})
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Minute)
		defer cancel()

		stackPath := stackProvider.GetStackPath(name)
		if err := composeManager.Pull(ctx, stackPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "pulled"})
	}
}

func GetComposeFile(stackProvider stack.Provider) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")
		content, err := stackProvider.GetComposeFile(name)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Compose file not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"content": content,
		})
	}
}

func UpdateComposeFile(stackProvider stack.Provider) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")

		var body struct {
			Content string `json:"content"`
		}

		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := stackProvider.UpdateComposeFile(name, body.Content); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "updated"})
	}
}

// Containers
func ListContainers(dockerClient docker.DockerClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		all := c.Query("all") == "true"
		containers, err := dockerClient.ListContainers(c.Request.Context(), all)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, containers)
	}
}

func GetContainer(dockerClient docker.DockerClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		container, err := dockerClient.GetContainer(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, container)
	}
}

func StartContainer(dockerClient docker.DockerClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if err := dockerClient.StartContainer(c.Request.Context(), id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "started"})
	}
}

func StopContainer(dockerClient docker.DockerClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if err := dockerClient.StopContainer(c.Request.Context(), id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "stopped"})
	}
}

func RestartContainer(dockerClient docker.DockerClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if err := dockerClient.RestartContainer(c.Request.Context(), id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "restarted"})
	}
}

func GetContainerLogs(dockerClient docker.DockerClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		tail := c.DefaultQuery("tail", "100")

		logs, err := dockerClient.GetContainerLogs(c.Request.Context(), id, tail)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"logs": logs})
	}
}

func GetContainerStats(dockerClient docker.DockerClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		stats, err := dockerClient.GetContainerStats(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, stats)
	}
}

// Volumes
func ListVolumes(dockerClient docker.DockerClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		volumes, err := dockerClient.ListVolumes(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, volumes)
	}
}

func CreateVolume(dockerClient docker.DockerClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body struct {
			Name   string            `json:"name"`
			Driver string            `json:"driver"`
			Labels map[string]string `json:"labels"`
		}

		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if body.Driver == "" {
			body.Driver = "local"
		}

		volume, err := dockerClient.CreateVolume(c.Request.Context(), body.Name, body.Driver, body.Labels)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, volume)
	}
}

func DeleteVolume(dockerClient docker.DockerClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")
		force := c.Query("force") == "true"

		if err := dockerClient.DeleteVolume(c.Request.Context(), name, force); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "deleted"})
	}
}

// Networks
func ListNetworks(dockerClient docker.DockerClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		networks, err := dockerClient.ListNetworks(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, networks)
	}
}

func CreateNetwork(dockerClient docker.DockerClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body struct {
			Name   string `json:"name"`
			Driver string `json:"driver"`
		}

		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if body.Driver == "" {
			body.Driver = "bridge"
		}

		network, err := dockerClient.CreateNetwork(c.Request.Context(), body.Name, body.Driver)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, network)
	}
}

func DeleteNetwork(dockerClient docker.DockerClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		if err := dockerClient.DeleteNetwork(c.Request.Context(), id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "deleted"})
	}
}

// Images
func ListImages(dockerClient docker.DockerClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		images, err := dockerClient.ListImages(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, images)
	}
}
