package domain

import (
	"encoding/json"
	"time"
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

// NewLowStockNotification crea una notificación para stock bajo
func NewLowStockNotification(productID string, stockLevel int) *Notification {
	return &Notification{
		Type:       LowStockNotification,
		Message:    "El producto está alcanzando niveles bajos de stock",
		Timestamp:  time.Now(),
		EntityID:   productID,
		StockLevel: stockLevel,
	}
}

// OrderNotification crea una notificación para un nuevo pedido
func OrderNotification(orderID string, amount float64, productsURL string) *Notification {
	return &Notification{
		Type:        NewOrderNotification,
		Message:     "Se ha creado una nueva orden",
		Timestamp:   time.Now(),
		EntityID:    orderID,
		Amount:      amount,
		ProductsURL: productsURL,
	}
}

// NewCancelOrderNotification crea una notificación para una orden cancelada
func NewCancelOrderNotification(orderID string, amount float64, provider string) *Notification {
	return &Notification{
		Type:      CancelOrderNotification,
		Message:   "La orden ha sido cancelada",
		Timestamp: time.Now(),
		EntityID:  orderID,
		Amount:    amount,
		Provider:  provider,
	}
}

// ToJSON convierte la notificación a JSON
func (n *Notification) ToJSON() ([]byte, error) {
	return json.Marshal(n)
}
