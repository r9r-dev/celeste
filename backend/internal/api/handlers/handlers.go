package handlers

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"

	"aperture-science-network/internal/compose"
	"aperture-science-network/internal/docker"
	"aperture-science-network/internal/stats"
)

// System Stats
func GetSystemStats(c *gin.Context) {
	sysStats, err := stats.GetSystemStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sysStats)
}

// StackInfo represents a compose stack with its status
type StackInfo struct {
	Name            string `json:"name"`
	Path            string `json:"path"`
	Status          string `json:"status"`
	Services        int    `json:"services"`
	RunningServices int    `json:"runningServices"`
}

// Stacks
func ListStacks(stacksPath string, dockerClient *docker.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		entries, err := os.ReadDir(stacksPath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Get all containers to determine stack status
		containers, err := dockerClient.ListContainers(c.Request.Context(), true)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Group containers by compose project
		containersByProject := make(map[string][]docker.ContainerInfo)
		for _, ctr := range containers {
			if project, ok := ctr.Labels["com.docker.compose.project"]; ok {
				containersByProject[project] = append(containersByProject[project], ctr)
			}
		}

		stacks := make([]StackInfo, 0)
		for _, entry := range entries {
			if entry.IsDir() {
				composePath := filepath.Join(stacksPath, entry.Name(), "docker-compose.yml")
				if _, err := os.Stat(composePath); err == nil {
					stackName := entry.Name()
					projectContainers := containersByProject[stackName]

					// Count running services
					runningCount := 0
					for _, ctr := range projectContainers {
						if ctr.State == "running" {
							runningCount++
						}
					}

					// Determine stack status
					status := "stopped"
					totalServices := len(projectContainers)
					if totalServices > 0 {
						if runningCount == totalServices {
							status = "running"
						} else if runningCount > 0 {
							status = "partial"
						}
					}

					stacks = append(stacks, StackInfo{
						Name:            stackName,
						Path:            filepath.Join(stacksPath, stackName),
						Status:          status,
						Services:        totalServices,
						RunningServices: runningCount,
					})
				}
			}
		}

		c.JSON(http.StatusOK, stacks)
	}
}

func GetStack(stacksPath string, dockerClient *docker.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")
		stackPath := filepath.Join(stacksPath, name)

		if _, err := os.Stat(stackPath); os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Stack not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"name": name,
			"path": stackPath,
		})
	}
}

func StartStack(stacksPath string, composeManager *compose.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")
		stackPath := filepath.Join(stacksPath, name)

		if _, err := os.Stat(stackPath); os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Stack not found"})
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Minute)
		defer cancel()

		if err := composeManager.Up(ctx, stackPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "started"})
	}
}

func StopStack(stacksPath string, composeManager *compose.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")
		stackPath := filepath.Join(stacksPath, name)

		if _, err := os.Stat(stackPath); os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Stack not found"})
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Minute)
		defer cancel()

		if err := composeManager.Down(ctx, stackPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "stopped"})
	}
}

func RestartStack(stacksPath string, composeManager *compose.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")
		stackPath := filepath.Join(stacksPath, name)

		if _, err := os.Stat(stackPath); os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Stack not found"})
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Minute)
		defer cancel()

		if err := composeManager.Restart(ctx, stackPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "restarted"})
	}
}

func PullStack(stacksPath string, composeManager *compose.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")
		stackPath := filepath.Join(stacksPath, name)

		if _, err := os.Stat(stackPath); os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Stack not found"})
			return
		}

		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Minute)
		defer cancel()

		if err := composeManager.Pull(ctx, stackPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "pulled"})
	}
}

func GetComposeFile(stacksPath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")
		composePath := filepath.Join(stacksPath, name, "docker-compose.yml")

		content, err := os.ReadFile(composePath)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Compose file not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"content": string(content),
		})
	}
}

func UpdateComposeFile(stacksPath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")
		composePath := filepath.Join(stacksPath, name, "docker-compose.yml")

		var body struct {
			Content string `json:"content"`
		}

		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := os.WriteFile(composePath, []byte(body.Content), 0644); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "updated"})
	}
}

// Containers
func ListContainers(dockerClient *docker.Client) gin.HandlerFunc {
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

func GetContainer(dockerClient *docker.Client) gin.HandlerFunc {
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

func StartContainer(dockerClient *docker.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if err := dockerClient.StartContainer(c.Request.Context(), id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "started"})
	}
}

func StopContainer(dockerClient *docker.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if err := dockerClient.StopContainer(c.Request.Context(), id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "stopped"})
	}
}

func RestartContainer(dockerClient *docker.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if err := dockerClient.RestartContainer(c.Request.Context(), id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "restarted"})
	}
}

func GetContainerLogs(dockerClient *docker.Client) gin.HandlerFunc {
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

func GetContainerStats(dockerClient *docker.Client) gin.HandlerFunc {
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
func ListVolumes(dockerClient *docker.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		volumes, err := dockerClient.ListVolumes(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, volumes)
	}
}

func CreateVolume(dockerClient *docker.Client) gin.HandlerFunc {
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

func DeleteVolume(dockerClient *docker.Client) gin.HandlerFunc {
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
func ListNetworks(dockerClient *docker.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		networks, err := dockerClient.ListNetworks(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, networks)
	}
}

func CreateNetwork(dockerClient *docker.Client) gin.HandlerFunc {
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

func DeleteNetwork(dockerClient *docker.Client) gin.HandlerFunc {
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
func ListImages(dockerClient *docker.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		images, err := dockerClient.ListImages(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, images)
	}
}
