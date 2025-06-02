package wsserver

import (
	"github.com/DrusGalkin/forum-client/internal/entity"
	"github.com/DrusGalkin/forum-client/internal/usecase"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
	"sync"
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	conn     *websocket.Conn
	send     chan entity.Post
	threadID int
}

type Hub struct {
	clients    map[*Client]bool
	chat       chan entity.Post
	register   chan *Client
	unregister chan *Client
	mu         sync.Mutex
	UseCase    usecase.PostUseCase
	logger     *zap.Logger
}

func NewHub(UseCase usecase.PostUseCase, logger *zap.Logger) *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		chat:       make(chan entity.Post),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		mu:         sync.Mutex{},
		UseCase:    UseCase,
		logger:     logger,
	}
}

func (h *Hub) Run() {
	h.logger.Info("Запуск хаба WebSocket")
	for {
		select {
		case client := <-h.register:
			h.logger.Debug("Регистрация нового клиента",
				zap.Int("threadID", client.threadID))
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()

		case client := <-h.unregister:
			h.logger.Debug("Отключение клиента",
				zap.Int("threadID", client.threadID))
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				close(client.send)
				delete(h.clients, client)
			}
			h.mu.Unlock()

		case message := <-h.chat:
			h.logger.Debug("Рассылка сообщения всем клиентам",
				zap.Int("threadID", message.ThreadID),
				zap.Int("userID", message.UserID),
				zap.String("content", message.Content))
			h.mu.Lock()
			for client := range h.clients {
				if client.threadID == message.ThreadID {
					select {
					case client.send <- message:
					default:
						close(client.send)
						delete(h.clients, client)
						h.logger.Warn("Канал клиента переполнен, отключение",
							zap.Int("threadID", client.threadID))
					}
				}
			}
			h.mu.Unlock()
		}
	}
}
