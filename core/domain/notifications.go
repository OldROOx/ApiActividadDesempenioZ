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

// NewLowStockNotification crea una nueva notificación de stock bajo
func NewLowStockNotification(productID string, stockLevel int) *Notification {
	return &Notification{
		Type:       LowStockNotification,
		Message:    "Alerta: Stock bajo",
		Timestamp:  time.Now(),
		EntityID:   productID,
		StockLevel: stockLevel,
	}
}

// OrderNotification crea una nueva notificación de creación de orden
func OrderNotification(orderID string, amount float64, productsURL string) *Notification {
	return &Notification{
		Type:        NewOrderNotification,
		Message:     "Nueva orden creada",
		Timestamp:   time.Now(),
		EntityID:    orderID,
		Amount:      amount,
		ProductsURL: productsURL,
	}
}

// NewCancelOrderNotification crea una nueva notificación de cancelación de orden
func NewCancelOrderNotification(orderID string, amount float64, provider string) *Notification {
	return &Notification{
		Type:      CancelOrderNotification,
		Message:   "Orden cancelada",
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
