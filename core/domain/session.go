package domain

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Session representa una sesión de conexión WebSocket
type Session struct {
	Conn         *websocket.Conn
	SessionID    string
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

// StartHandling inicia el manejo de la sesión
func (s *Session) StartHandling(removeSession func(sessionID string)) {
	s.CloseHandler = removeSession
	s.readPump()
}

// readPump escucha mensajes entrantes
func (s *Session) readPump() {
	defer func() {
		s.Conn.Close()
		if s.CloseHandler != nil {
			s.CloseHandler(s.SessionID)
		}
	}()

	for {
		messageType, _, err := s.Conn.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(
				err,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
			) {
				log.Printf("Error en la sesión %s: %v", s.SessionID, err)
			}
			break
		}

		if messageType != -1 {
			// Simple eco para mensajes del cliente
			message := fmt.Sprintf("Recibido mensaje del cliente %s", s.SessionID)
			s.SendMessage(websocket.TextMessage, []byte(message))
		}

		time.Sleep(17 * time.Millisecond)
	}
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

// SendMessage envía un mensaje a la sesión
func (s *Session) SendMessage(messageType int, payloadByte []byte) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	err := s.Conn.WriteMessage(messageType, payloadByte)
	if err != nil {
		log.Printf("Error al enviar mensaje a la sesión %s: %v", s.SessionID, err)
	}
}
