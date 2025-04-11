package domain

import (
	"time"
)

// Modelos base
type Proveedor struct {
	ID            int    `json:"id_proveedor"`
	Nombre        string `json:"nombre"`
	Direccion     string `json:"direccion"`
	Telefono      string `json:"telefono"`
	Email         string `json:"email"`
	FechaRegistro string `json:"fecha_registro"`
}

type Producto struct {
	ID            int    `json:"id_producto"`
	Nombre        string `json:"nombre"`
	Descripcion   string `json:"descripcion"`
	Precio        int    `json:"precio"`
	Existencia    int    `json:"existencia"`
	ProveedorID   int    `json:"id_proveedor"`
	FechaCreacion string `json:"fecha_creacion"`
}

type Pedido struct {
	ID          int     `json:"id_pedido"`
	FechaPedido string  `json:"fecha_pedido"`
	Estado      string  `json:"estado"`
	Total       float64 `json:"total"`
}

type DetallesPedido struct {
	ID             int     `json:"id_detalle_pedido"`
	PedidoID       int     `json:"id_pedido"`
	ProductoID     int     `json:"id_producto"`
	Cantidad       int     `json:"cantidad"`
	PrecioUnitario float64 `json:"precio_unitario"`
	Subtotal       float64 `json:"subtotal"`
}

type Venta struct {
	ID         int       `json:"id_venta"`
	FechaVenta time.Time `json:"fecha_venta"`
	Estado     string    `json:"estado"`
	Total      float64   `json:"total"`
}

type DetallesVenta struct {
	ID             int     `json:"id_detalle_venta"`
	VentaID        int     `json:"id_venta"`
	ProductoID     int     `json:"id_producto"`
	Cantidad       int     `json:"cantidad"`
	PrecioUnitario float64 `json:"precio_unitario"`
	Subtotal       float64 `json:"subtotal"`
}

type OrdenProveedor struct {
	ID          int    `json:"id_orden_proveedor"`
	ProveedorID int    `json:"id_proveedor"`
	FechaOrden  string `json:"fecha_orden"`
	Estado      string `json:"estado"`
	Total       int    `json:"total"`
}

type DetallesOrden struct {
	ID               int     `json:"id_detalle_orden"`
	OrdenProveedorID int     `json:"id_orden_proveedor"`
	ProductoID       int     `json:"id_producto"`
	Cantidad         int     `json:"cantidad"`
	PrecioUnitario   float64 `json:"precio_unitario"`
	Subtotal         float64 `json:"subtotal"`
}
