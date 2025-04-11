package handlers

import (
	"ActividadDesempenioAPIz/core/ports"
)

// ControllerFactory crea y retorna todos los controladores necesarios para la API
type ControllerFactory struct {
	productoController       *ProductoController
	proveedorController      *ProveedorController
	pedidoController         *PedidoController
	ventaController          *VentaController
	ordenProveedorController *OrdenProveedorController
}

// NewControllerFactory crea una nueva fábrica de controladores
func NewControllerFactory(
	productoRepo ports.ProductoRepository,
	proveedorRepo ports.ProveedorRepository,
	pedidoRepo ports.PedidoRepository,
	detallesPedidoRepo ports.DetallesPedidoRepository,
	ventaRepo ports.VentaRepository,
	detallesVentaRepo ports.DetallesVentaRepository,
	ordenRepo ports.OrdenProveedorRepository,
	detallesOrdenRepo ports.DetallesOrdenRepository,
	notificationService ports.NotificationService,
) *ControllerFactory {
	productoController := NewProductoController(productoRepo, notificationService)
	proveedorController := NewProveedorController(proveedorRepo)
	pedidoController := NewPedidoController(pedidoRepo, detallesPedidoRepo, productoRepo, notificationService)
	ventaController := NewVentaController(ventaRepo, detallesVentaRepo, productoRepo, notificationService)
	ordenProveedorController := NewOrdenProveedorController(ordenRepo, detallesOrdenRepo, proveedorRepo, productoRepo, notificationService)

	return &ControllerFactory{
		productoController:       productoController,
		proveedorController:      proveedorController,
		pedidoController:         pedidoController,
		ventaController:          ventaController,
		ordenProveedorController: ordenProveedorController,
	}
}

// GetProductoController retorna el controlador de productos
func (cf *ControllerFactory) GetProductoController() *ProductoController {
	return cf.productoController
}

// GetProveedorController retorna el controlador de proveedores
func (cf *ControllerFactory) GetProveedorController() *ProveedorController {
	return cf.proveedorController
}

// GetPedidoController retorna el controlador de pedidos
func (cf *ControllerFactory) GetPedidoController() *PedidoController {
	return cf.pedidoController
}

// GetVentaController retorna el controlador de ventas
func (cf *ControllerFactory) GetVentaController() *VentaController {
	return cf.ventaController
}

// GetOrdenProveedorController retorna el controlador de órdenes de proveedor
func (cf *ControllerFactory) GetOrdenProveedorController() *OrdenProveedorController {
	return cf.ordenProveedorController
}
