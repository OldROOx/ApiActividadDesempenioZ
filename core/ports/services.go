package ports

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
