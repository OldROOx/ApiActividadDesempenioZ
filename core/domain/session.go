package domain

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// Session representa una sesión de WebSocket
type Session struct {
	SessionID    string
	Conn         *websocket.Conn
	CloseHandler func(sessionID string)
	Sessions     map[string]*Session
	Mutex        sync.Mutex
}

// NewSession crea una nueva sesión
func NewSession(conn *websocket.Conn, sessionID string, sessions map[string]*Session) *Session {
	return &Session{
		Conn:      conn,
		SessionID: sessionID,
		Sessions:  sessions,
		Mutex:     sync.Mutex{},
	}
}

// SetCloseHandler establece el manejador de cierre
func (s *Session) SetCloseHandler(handler func(sessionID string)) {
	s.CloseHandler = handler
}

// SendMessage envía un mensaje a la sesión
func (s *Session) SendMessage(messageType int, payloadByte []byte) error {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	err := s.Conn.WriteMessage(messageType, payloadByte)
	if err != nil {
		log.Printf("Error al enviar mensaje a la sesión %s: %v", s.SessionID, err)
		return err
	}
	return nil
}

// BroadcastToAll envía un mensaje a todas las sesiones conectadas
func (s *Session) BroadcastToAll(messageType int, payload []byte) {
	for _, session := range s.Sessions {
		session.SendMessage(messageType, payload)
	}
}

// BroadcastNotification envía una notificación a todos los clientes conectados
func (s *Session) BroadcastNotification(notification *Notification) {
	payload, err := notification.ToJSON()
	if err != nil {
		log.Printf("Error al serializar la notificación: %v", err)
		return
	}

	s.BroadcastToAll(websocket.TextMessage, payload)
	log.Printf("Notificación enviada: %s", string(payload))
}
