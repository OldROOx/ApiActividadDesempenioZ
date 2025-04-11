package websocket

import (
	"ActividadDesempenioAPIz/core/ports"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// WebsocketHandler es la base para los manejadores de WebSocket
type WebsocketHandler struct {
	wsService ports.WebSocketService
}

// ProductStockWebsocketHandler maneja las conexiones WebSocket para notificaciones de stock
type ProductStockWebsocketHandler struct {
	WebsocketHandler
}

// NewProductStockWebsocketHandler crea un nuevo manejador de WebSocket para stock
func NewProductStockWebsocketHandler(wsService ports.WebSocketService) *ProductStockWebsocketHandler {
	return &ProductStockWebsocketHandler{
		WebsocketHandler: WebsocketHandler{
			wsService: wsService,
		},
	}
}

// Handle maneja una nueva conexión WebSocket
func (wh *ProductStockWebsocketHandler) Handle(c *gin.Context) {
	sessionID := c.Query("session_id")
	err := wh.wsService.HandleConnection(c.Writer, c.Request, sessionID)
	if err != nil {
		log.Printf("Error al manejar conexión WebSocket de stock: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

// OrderCreationWebsocketHandler maneja las conexiones WebSocket para notificaciones de creación de órdenes
type OrderCreationWebsocketHandler struct {
	WebsocketHandler
}

// NewOrderCreationWebsocketHandler crea un nuevo manejador de WebSocket para creación de órdenes
func NewOrderCreationWebsocketHandler(wsService ports.WebSocketService) *OrderCreationWebsocketHandler {
	return &OrderCreationWebsocketHandler{
		WebsocketHandler: WebsocketHandler{
			wsService: wsService,
		},
	}
}

// Handle maneja una nueva conexión WebSocket
func (wh *OrderCreationWebsocketHandler) Handle(c *gin.Context) {
	sessionID := c.Query("session_id")
	err := wh.wsService.HandleConnection(c.Writer, c.Request, sessionID)
	if err != nil {
		log.Printf("Error al manejar conexión WebSocket de creación de órdenes: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

// OrderCancelWebsocketHandler maneja las conexiones WebSocket para notificaciones de cancelación
type OrderCancelWebsocketHandler struct {
	WebsocketHandler
}

// NewOrderCancelWebsocketHandler crea un nuevo manejador de WebSocket para cancelaciones
func NewOrderCancelWebsocketHandler(wsService ports.WebSocketService) *OrderCancelWebsocketHandler {
	return &OrderCancelWebsocketHandler{
		WebsocketHandler: WebsocketHandler{
			wsService: wsService,
		},
	}
}

// Handle maneja una nueva conexión WebSocket
func (wh *OrderCancelWebsocketHandler) Handle(c *gin.Context) {
	sessionID := c.Query("session_id")
	err := wh.wsService.HandleConnection(c.Writer, c.Request, sessionID)
	if err != nil {
		log.Printf("Error al manejar conexión WebSocket de cancelación: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}
