package application

import (
	"ActividadDesempenioAPIz/core/domain"
	"ActividadDesempenioAPIz/infrastructure/websocket"
	"log"
	"strconv"
	"sync"
)

// NotificationServiceExtended extiende el servicio de notificaciones
type NotificationServiceExtended struct {
	productStockWS  *websocket.WebsocketService
	orderCreationWS *websocket.WebsocketService
	orderCancelWS   *websocket.WebsocketService
	proveedorRepo   ProveedorRepository
	mutex           sync.RWMutex
}

// NewNotificationServiceExtended crea un nuevo servicio de notificaciones extendido
func NewNotificationServiceExtended(
	productStockWS *websocket.WebsocketService,
	orderCreationWS *websocket.WebsocketService,
	orderCancelWS *websocket.WebsocketService,
	proveedorRepo ProveedorRepository,
) *NotificationServiceExtended {
	return &NotificationServiceExtended{
		productStockWS:  productStockWS,
		orderCreationWS: orderCreationWS,
		orderCancelWS:   orderCancelWS,
		proveedorRepo:   proveedorRepo,
		mutex:           sync.RWMutex{},
	}
}

// ProveedorRepository define la interfaz para acceder a proveedores
type ProveedorRepository interface {
	GetByID(id int) (*domain.Proveedor, error)
}

// GetProductStockWS retorna el servicio WebSocket para notificaciones de stock
func (ns *NotificationServiceExtended) GetProductStockWS() *websocket.WebsocketService {
	return ns.productStockWS
}

// GetOrderCreationWS retorna el servicio WebSocket para notificaciones de creación de órdenes
func (ns *NotificationServiceExtended) GetOrderCreationWS() *websocket.WebsocketService {
	return ns.orderCreationWS
}

// GetOrderCancelWS retorna el servicio WebSocket para notificaciones de cancelación de órdenes
func (ns *NotificationServiceExtended) GetOrderCancelWS() *websocket.WebsocketService {
	return ns.orderCancelWS
}

// NotifyLowStock envía una notificación cuando un producto tiene poco stock
func (ns *NotificationServiceExtended) NotifyLowStock(productID int, stockLevel int) {
	ns.mutex.RLock()
	defer ns.mutex.RUnlock()

	if stockLevel <= 5 {
		productIDStr := strconv.Itoa(productID)
		notification := domain.NewLowStockNotification(productIDStr, stockLevel)
		payload, err := notification.ToJSON()
		if err != nil {
			log.Printf("Error al serializar notificación de stock: %v", err)
			return
		}

		ns.productStockWS.Broadcast(payload)
		log.Printf("Notificación de stock bajo para producto %d con nivel de stock %d",
			productID, stockLevel)
	}
}

// NotifyNewPedido envía una notificación cuando se crea un nuevo pedido
func (ns *NotificationServiceExtended) NotifyNewPedido(pedidoID int, amount float64) {
	ns.mutex.RLock()
	defer ns.mutex.RUnlock()

	pedidoIDStr := strconv.Itoa(pedidoID)
	productsURL := "/api/pedidos/" + pedidoIDStr + "/productos"
	notification := domain.OrderNotification(pedidoIDStr, amount, productsURL)
	payload, err := notification.ToJSON()
	if err != nil {
		log.Printf("Error al serializar notificación de pedido: %v", err)
		return
	}

	ns.orderCreationWS.Broadcast(payload)
	log.Printf("Notificación de nuevo pedido para pedido %d con monto %.2f",
		pedidoID, amount)
}

// NotifyNewVenta envía una notificación cuando se crea una nueva venta
func (ns *NotificationServiceExtended) NotifyNewVenta(ventaID int, amount float64) {
	ns.mutex.RLock()
	defer ns.mutex.RUnlock()

	ventaIDStr := strconv.Itoa(ventaID)
	productsURL := "/api/ventas/" + ventaIDStr + "/productos"
	notification := domain.OrderNotification(ventaIDStr, amount, productsURL)
	payload, err := notification.ToJSON()
	if err != nil {
		log.Printf("Error al serializar notificación de venta: %v", err)
		return
	}

	ns.orderCreationWS.Broadcast(payload)
	log.Printf("Notificación de nueva venta para venta %d con monto %.2f",
		ventaID, amount)
}

// NotifyNewOrdenProveedor envía una notificación cuando se crea una nueva orden de proveedor
func (ns *NotificationServiceExtended) NotifyNewOrdenProveedor(ordenID int, amount float64) {
	ns.mutex.RLock()
	defer ns.mutex.RUnlock()

	ordenIDStr := strconv.Itoa(ordenID)
	productsURL := "/api/ordenes/" + ordenIDStr + "/productos"
	notification := domain.OrderNotification(ordenIDStr, amount, productsURL)
	payload, err := notification.ToJSON()
	if err != nil {
		log.Printf("Error al serializar notificación de orden: %v", err)
		return
	}

	ns.orderCreationWS.Broadcast(payload)
	log.Printf("Notificación de nueva orden de proveedor para orden %d con monto %.2f",
		ordenID, amount)
}

// NotifyCanceledPedido envía una notificación cuando se cancela un pedido
func (ns *NotificationServiceExtended) NotifyCanceledPedido(pedidoID int, amount float64) {
	ns.mutex.RLock()
	defer ns.mutex.RUnlock()

	pedidoIDStr := strconv.Itoa(pedidoID)
	notification := domain.NewCancelOrderNotification(pedidoIDStr, amount, "")
	payload, err := notification.ToJSON()
	if err != nil {
		log.Printf("Error al serializar notificación de cancelación: %v", err)
		return
	}

	ns.orderCancelWS.Broadcast(payload)
	log.Printf("Notificación de pedido cancelado para pedido %d con monto %.2f",
		pedidoID, amount)
}

// NotifyCanceledVenta envía una notificación cuando se cancela una venta
func (ns *NotificationServiceExtended) NotifyCanceledVenta(ventaID int, amount float64) {
	ns.mutex.RLock()
	defer ns.mutex.RUnlock()

	ventaIDStr := strconv.Itoa(ventaID)
	notification := domain.NewCancelOrderNotification(ventaIDStr, amount, "")
	payload, err := notification.ToJSON()
	if err != nil {
		log.Printf("Error al serializar notificación de cancelación: %v", err)
		return
	}

	ns.orderCancelWS.Broadcast(payload)
	log.Printf("Notificación de venta cancelada para venta %d con monto %.2f",
		ventaID, amount)
}

// NotifyCanceledOrdenProveedor envía una notificación cuando se cancela una orden de proveedor
func (ns *NotificationServiceExtended) NotifyCanceledOrdenProveedor(ordenID int, amount float64, providerID int) {
	ns.mutex.RLock()
	defer ns.mutex.RUnlock()

	// Obtener información del proveedor
	provider, err := ns.proveedorRepo.GetByID(providerID)
	providerName := ""
	if err == nil && provider != nil {
		providerName = provider.Nombre
	}

	ordenIDStr := strconv.Itoa(ordenID)
	notification := domain.NewCancelOrderNotification(ordenIDStr, amount, providerName)
	payload, err := notification.ToJSON()
	if err != nil {
		log.Printf("Error al serializar notificación de cancelación: %v", err)
		return
	}

	ns.orderCancelWS.Broadcast(payload)
	log.Printf("Notificación de orden cancelada para orden %d con monto %.2f y proveedor %s",
		ordenID, amount, providerName)
}
