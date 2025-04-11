package ports

import (
	"ActividadDesempenioAPIz/core/domain"
)

// Interfaces para repositorios
type ProductoRepository interface {
	GetByID(id int) (*domain.Producto, error)
	GetAll() ([]*domain.Producto, error)
	Create(producto *domain.Producto) (int, error)
	Update(producto *domain.Producto) error
	UpdateStock(id int, cantidad int) error
	Delete(id int) error
}

type ProveedorRepository interface {
	GetByID(id int) (*domain.Proveedor, error)
	GetAll() ([]*domain.Proveedor, error)
	Create(proveedor *domain.Proveedor) (int, error)
	Update(proveedor *domain.Proveedor) error
	Delete(id int) error
}

type PedidoRepository interface {
	GetByID(id int) (*domain.Pedido, error)
	GetAll() ([]*domain.Pedido, error)
	Create(pedido *domain.Pedido) (int, error)
	Update(pedido *domain.Pedido) error
	UpdateEstado(id int, estado string) error
	Delete(id int) error
}

type DetallesPedidoRepository interface {
	GetByPedidoID(pedidoID int) ([]*domain.DetallesPedido, error)
	Create(detalle *domain.DetallesPedido) (int, error)
	Update(detalle *domain.DetallesPedido) error
	Delete(id int) error
}

type VentaRepository interface {
	GetByID(id int) (*domain.Venta, error)
	GetAll() ([]*domain.Venta, error)
	Create(venta *domain.Venta) (int, error)
	Update(venta *domain.Venta) error
	UpdateEstado(id int, estado string) error
	Delete(id int) error
}

type DetallesVentaRepository interface {
	GetByVentaID(ventaID int) ([]*domain.DetallesVenta, error)
	Create(detalle *domain.DetallesVenta) (int, error)
	Update(detalle *domain.DetallesVenta) error
	Delete(id int) error
}

type OrdenProveedorRepository interface {
	GetByID(id int) (*domain.OrdenProveedor, error)
	GetAll() ([]*domain.OrdenProveedor, error)
	Create(orden *domain.OrdenProveedor) (int, error)
	Update(orden *domain.OrdenProveedor) error
	UpdateEstado(id int, estado string) error
	Delete(id int) error
}

type DetallesOrdenRepository interface {
	GetByOrdenID(ordenID int) ([]*domain.DetallesOrden, error)
	Create(detalle *domain.DetallesOrden) (int, error)
	Update(detalle *domain.DetallesOrden) error
	Delete(id int) error
}
