package websocket

import (
	"ActividadDesempenioAPIz/core/domain"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
)

// WebsocketService implementa el servicio de WebSocket
type WebsocketService struct {
	upgrader      websocket.Upgrader
	clients       map[*websocket.Conn]bool
	clientsMutex  sync.RWMutex
	sessions      map[string]*domain.Session
	sessionsMutex sync.RWMutex
	nextSessionID int
}

// NewWebsocketService crea un nuevo servicio de WebSocket
func NewWebsocketService() *WebsocketService {
	return &WebsocketService{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		clients:       make(map[*websocket.Conn]bool),
		clientsMutex:  sync.RWMutex{},
		sessions:      make(map[string]*domain.Session),
		sessionsMutex: sync.RWMutex{},
		nextSessionID: 1,
	}
}

// RegisterClient registra un nuevo cliente WebSocket
func (ws *WebsocketService) RegisterClient(conn interface{}) interface{} {
	wsConn, ok := conn.(*websocket.Conn)
	if !ok {
		log.Printf("Error: se esperaba una conexión websocket.Conn")
		return nil
	}

	ws.clientsMutex.Lock()
	defer ws.clientsMutex.Unlock()
	ws.clients[wsConn] = true
	return wsConn
}

// UnregisterClient elimina un cliente WebSocket
func (ws *WebsocketService) UnregisterClient(conn interface{}) {
	wsConn, ok := conn.(*websocket.Conn)
	if !ok {
		log.Printf("Error: se esperaba una conexión websocket.Conn")
		return
	}

	ws.clientsMutex.Lock()
	defer ws.clientsMutex.Unlock()
	delete(ws.clients, wsConn)
	wsConn.Close()
}

// Broadcast envía un mensaje a todos los clientes conectados
func (ws *WebsocketService) Broadcast(message []byte) {
	ws.clientsMutex.RLock()
	defer ws.clientsMutex.RUnlock()

	for client := range ws.clients {
		err := client.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Printf("Error broadcasting to client: %v", err)
		}
	}
}

// HandleConnection maneja una nueva conexión WebSocket
func (ws *WebsocketService) HandleConnection(
	w http.ResponseWriter, r *http.Request, sessionID string,
) error {
	conn, err := ws.upgrader.Upgrade(w, r, nil)

	if err != nil {
		return err
	}

	log.Println("Nueva conexión WebSocket con ID de sesión:", sessionID)

	if sessionID == "" {
		sessionID = ws.generateSessionID()
		log.Println("ID de sesión generado:", sessionID)
	}

	// Registramos el cliente
	ws.RegisterClient(conn)

	session := domain.NewSession(conn, sessionID, ws.sessions)

	ws.addSession(sessionID, session)

	// Inicia el manejo en una goroutine para no bloquear
	go session.StartHandling(ws.removeSession)

	return nil
}

// generateSessionID genera un nuevo ID de sesión
func (ws *WebsocketService) generateSessionID() string {
	ws.sessionsMutex.Lock()
	defer ws.sessionsMutex.Unlock()

	id := ws.nextSessionID
	ws.nextSessionID++
	return strconv.Itoa(id)
}

// addSession agrega una sesión
func (ws *WebsocketService) addSession(sessionID string, session *domain.Session) {
	ws.sessionsMutex.Lock()
	defer ws.sessionsMutex.Unlock()

	ws.sessions[sessionID] = session
	log.Printf("Sesión %s agregada, sesiones totales: %d", sessionID, len(ws.sessions))
}

// removeSession elimina una sesión
func (ws *WebsocketService) removeSession(sessionID string) {
	ws.sessionsMutex.Lock()
	defer ws.sessionsMutex.Unlock()

	delete(ws.sessions, sessionID)
	log.Printf("Sesión %s eliminada, sesiones restantes: %d", sessionID, len(ws.sessions))
}

// GetSessions retorna las sesiones activas
func (ws *WebsocketService) GetSessions() map[string]*domain.Session {
	ws.sessionsMutex.RLock()
	defer ws.sessionsMutex.RUnlock()
	return ws.sessions
}
