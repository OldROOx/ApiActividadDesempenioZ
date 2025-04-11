package ports

import (
	"ActividadDesempenioAPIz/core/domain"
	"net/http"
)

// WebSocketService define el servicio para conexiones WebSocket
type WebSocketService interface {
	HandleConnection(w http.ResponseWriter, r *http.Request, sessionID string) error
	RegisterClient(conn interface{}) interface{}
	UnregisterClient(conn interface{})
	Broadcast(message []byte)
	GetSessions() map[string]*domain.Session
}
