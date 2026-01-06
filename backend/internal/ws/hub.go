package ws

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"aperture-science-network/internal/docker"
	"aperture-science-network/internal/stats"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // In production, check origin properly
	},
}

type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

type Hub struct {
	clients      map[*Client]bool
	broadcast    chan []byte
	register     chan *Client
	unregister   chan *Client
	mutex        sync.RWMutex
	dockerClient *docker.Client
}

func NewHub(dockerClient *docker.Client) *Hub {
	return &Hub{
		clients:      make(map[*Client]bool),
		broadcast:    make(chan []byte),
		register:     make(chan *Client),
		unregister:   make(chan *Client),
		dockerClient: dockerClient,
	}
}

func (h *Hub) Run() {
	// Start stats broadcasters
	go h.broadcastSystemStats()
	go h.broadcastContainerStats()

	for {
		select {
		case client := <-h.register:
			h.mutex.Lock()
			h.clients[client] = true
			h.mutex.Unlock()
			log.Printf("Client connected. Total clients: %d", len(h.clients))

		case client := <-h.unregister:
			h.mutex.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mutex.Unlock()
			log.Printf("Client disconnected. Total clients: %d", len(h.clients))

		case message := <-h.broadcast:
			h.mutex.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mutex.RUnlock()
		}
	}
}

func (h *Hub) broadcastSystemStats() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		h.mutex.RLock()
		clientCount := len(h.clients)
		h.mutex.RUnlock()

		if clientCount == 0 {
			continue
		}

		sysStats, err := stats.GetSystemStats()
		if err != nil {
			log.Printf("Error getting system stats: %v", err)
			continue
		}

		msg := Message{
			Type:    "stats",
			Payload: sysStats,
		}

		data, err := json.Marshal(msg)
		if err != nil {
			log.Printf("Error marshaling stats: %v", err)
			continue
		}

		h.broadcast <- data
	}
}

// ContainerStatsPayload holds stats for all containers
type ContainerStatsPayload struct {
	Containers map[string]*docker.ContainerStats `json:"containers"`
	Timestamp  int64                              `json:"timestamp"`
}

func (h *Hub) broadcastContainerStats() {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		h.mutex.RLock()
		clientCount := len(h.clients)
		h.mutex.RUnlock()

		if clientCount == 0 || h.dockerClient == nil {
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		// Get running containers
		containers, err := h.dockerClient.ListContainers(ctx, false)
		if err != nil {
			cancel()
			log.Printf("Error listing containers: %v", err)
			continue
		}

		// Collect stats for each container
		containerStats := make(map[string]*docker.ContainerStats)
		for _, ctr := range containers {
			ctrStats, err := h.dockerClient.GetContainerStats(ctx, ctr.ID)
			if err != nil {
				continue // Skip containers we can't get stats for
			}
			containerStats[ctr.ID] = ctrStats
		}
		cancel()

		if len(containerStats) == 0 {
			continue
		}

		msg := Message{
			Type: "container_stats",
			Payload: ContainerStatsPayload{
				Containers: containerStats,
				Timestamp:  time.Now().Unix(),
			},
		}

		data, err := json.Marshal(msg)
		if err != nil {
			log.Printf("Error marshaling container stats: %v", err)
			continue
		}

		h.broadcast <- data
	}
}

func HandleWebSocket(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := &Client{
		hub:  hub,
		conn: conn,
		send: make(chan []byte, 256),
	}

	hub.register <- client

	go client.writePump()
	go client.readPump()
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(512)
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
