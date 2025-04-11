package infrastructure

import (
	"ActividadDesempenioAPIz/infrastructure/websocket"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ProductStockWebsocketHandler maneja las conexiones WebSocket para notificaciones de stock
type ProductStockWebsocketHandler struct {
	wsService *websocket.WebsocketService
}

// NewProductStockWebsocketHandler crea un nuevo manejador de WebSocket para stock
func NewProductStockWebsocketHandler(wsService *websocket.WebsocketService) *ProductStockWebsocketHandler {
	return &ProductStockWebsocketHandler{
		wsService: wsService,
	}
}

// Handle maneja una nueva conexión WebSocket
func (wh *ProductStockWebsocketHandler) Handle(c *gin.Context) {
	sessionID := c.Query("session_id")
	if sessionID == "" {
		sessionID = "stock-" + c.ClientIP()
	}

	err := wh.wsService.HandleConnection(c.Writer, c.Request, sessionID)
	if err != nil {
		log.Printf("Error al manejar conexión WebSocket de stock: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

// OrderCreationWebsocketHandler maneja las conexiones WebSocket para notificaciones de creación de órdenes
type OrderCreationWebsocketHandler struct {
	wsService *websocket.WebsocketService
}

// NewOrderCreationWebsocketHandler crea un nuevo manejador de WebSocket para creación de órdenes
func NewOrderCreationWebsocketHandler(wsService *websocket.WebsocketService) *OrderCreationWebsocketHandler {
	return &OrderCreationWebsocketHandler{
		wsService: wsService,
	}
}

// Handle maneja una nueva conexión WebSocket
func (wh *OrderCreationWebsocketHandler) Handle(c *gin.Context) {
	sessionID := c.Query("session_id")
	if sessionID == "" {
		sessionID = "orders-" + c.ClientIP()
	}

	err := wh.wsService.HandleConnection(c.Writer, c.Request, sessionID)
	if err != nil {
		log.Printf("Error al manejar conexión WebSocket de creación de órdenes: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

// OrderCancelWebsocketHandler maneja las conexiones WebSocket para notificaciones de cancelación
type OrderCancelWebsocketHandler struct {
	wsService *websocket.WebsocketService
}

// NewOrderCancelWebsocketHandler crea un nuevo manejador de WebSocket para cancelaciones
func NewOrderCancelWebsocketHandler(wsService *websocket.WebsocketService) *OrderCancelWebsocketHandler {
	return &OrderCancelWebsocketHandler{
		wsService: wsService,
	}
}

// Handle maneja una nueva conexión WebSocket
func (wh *OrderCancelWebsocketHandler) Handle(c *gin.Context) {
	sessionID := c.Query("session_id")
	if sessionID == "" {
		sessionID = "cancellations-" + c.ClientIP()
	}

	err := wh.wsService.HandleConnection(c.Writer, c.Request, sessionID)
	if err != nil {
		log.Printf("Error al manejar conexión WebSocket de cancelación: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}
