package ports

import (
	"ActividadDesempenioAPIz/core/domain"
	"net/http"
)

// Interfaces para servicios

// NotificationService define el servicio básico de notificaciones
type NotificationService interface {
	NotifyLowStock(productID int, stockLevel int)
	NotifyNewPedido(pedidoID int, amount float64)
	NotifyNewVenta(ventaID int, amount float64)
	NotifyNewOrdenProveedor(ordenID int, amount float64)
	NotifyCanceledPedido(pedidoID int, amount float64)
	NotifyCanceledVenta(ventaID int, amount float64)
	NotifyCanceledOrdenProveedor(ordenID int, amount float64, providerID int)
}

// WebSocketService define el servicio para conexiones WebSocket
type WebSocketService interface {
	HandleConnection(w http.ResponseWriter, r *http.Request, sessionID string) error
	RegisterClient(conn interface{}) interface{}
	UnregisterClient(conn interface{})
	Broadcast(message []byte)
	GetSessions() map[string]*domain.Session
}
