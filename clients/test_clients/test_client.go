package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

// NotificationType define el tipo de notificación
type NotificationType string

const (
	LowStockNotification    NotificationType = "low_stock"
	NewOrderNotification    NotificationType = "new_order"
	CancelOrderNotification NotificationType = "cancel_order"
)

// Notification representa una notificación del sistema
type Notification struct {
	Type        NotificationType `json:"type"`
	Message     string           `json:"message"`
	Timestamp   time.Time        `json:"timestamp"`
	EntityID    string           `json:"entity_id"`
	Amount      float64          `json:"amount,omitempty"`
	StockLevel  int              `json:"stock_level,omitempty"`
	Provider    string           `json:"provider,omitempty"`
	ProductsURL string           `json:"products_url,omitempty"`
}

func main() {
	// Verificar argumentos
	if len(os.Args) < 2 {
		log.Println("Uso: go run test_client.go [stock|orders|cancellations]")
		os.Exit(1)
	}

	notificationType := os.Args[1]
	var endpoint string

	switch notificationType {
	case "stock":
		endpoint = "/ws/stock"
	case "orders":
		endpoint = "/ws/orders"
	case "cancellations":
		endpoint = "/ws/cancellations"
	default:
		log.Fatalf("Tipo de notificación no válido: %s. Debe ser 'stock', 'orders' o 'cancellations'", notificationType)
	}

	// Crear un canal para manejar señales
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Conectar al servidor WebSocket
	u := url.URL{
		Scheme:   "ws",
		Host:     "localhost:4000",
		Path:     endpoint,
		RawQuery: "session_id=test-client-" + notificationType,
	}

	log.Printf("Conectando a %s\n", u.String())
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatalf("Error al conectar: %v", err)
	}
	defer c.Close()

	// Canal para recibir mensajes del servidor
	done := make(chan struct{})

	// Iniciar una goroutine para leer mensajes del servidor
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Printf("Error al leer mensaje: %v", err)
				return
			}

			// Intentar parsear el mensaje como una notificación
			var notification Notification
			if err := json.Unmarshal(message, &notification); err != nil {
				// No es una notificación, imprimir el mensaje crudo
				log.Printf("Mensaje recibido: %s", message)
				continue
			}

			// Manejar diferentes tipos de notificación
			switch notification.Type {
			case LowStockNotification:
				fmt.Printf("\n⚠️ ALERTA DE STOCK BAJO ⚠️\n")
				fmt.Printf("ID del Producto: %s\n", notification.EntityID)
				fmt.Printf("Stock Actual: %d unidades\n", notification.StockLevel)
				fmt.Printf("Hora: %s\n\n", notification.Timestamp.Format(time.RFC1123))

			case NewOrderNotification:
				fmt.Printf("\n🛒 NUEVA ORDEN CREADA 🛒\n")
				fmt.Printf("ID de la Orden: %s\n", notification.EntityID)
				fmt.Printf("Monto Total: $%.2f\n", notification.Amount)
				fmt.Printf("Productos: %s\n", notification.ProductsURL)
				fmt.Printf("Hora: %s\n\n", notification.Timestamp.Format(time.RFC1123))

			case CancelOrderNotification:
				fmt.Printf("\n❌ ORDEN CANCELADA ❌\n")
				fmt.Printf("ID de la Orden: %s\n", notification.EntityID)
				fmt.Printf("Monto: $%.2f\n", notification.Amount)
				if notification.Provider != "" {
					fmt.Printf("Proveedor: %s\n", notification.Provider)
				}
				fmt.Printf("Hora: %s\n\n", notification.Timestamp.Format(time.RFC1123))

			default:
				fmt.Printf("Tipo de notificación desconocido: %s\n", notification.Type)
			}
		}
	}()

	// Enviar un latido periódicamente
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			// Enviar un mensaje de latido
			err := c.WriteMessage(websocket.TextMessage, []byte("heartbeat"))
			if err != nil {
				log.Println("Error al enviar latido:", err)
				return
			}
		case <-interrupt:
			// Cerrar la conexión correctamente
			log.Println("Interrupción recibida, cerrando conexión...")
			err := c.WriteMessage(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("Error al cerrar WebSocket:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
