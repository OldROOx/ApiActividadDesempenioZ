package routes

import (
	"ActividadDesempenioAPIz/core/ports"
	"ActividadDesempenioAPIz/infrastructure/api/handlers"
	"ActividadDesempenioAPIz/infrastructure/websocket"

	"github.com/gin-gonic/gin"
)

// SetupRouter configura todas las rutas de la aplicación
func SetupRouter(
	engine *gin.Engine,
	productRepo ports.ProductoRepository,
	proveedorRepo ports.ProveedorRepository,
	pedidoRepo ports.PedidoRepository,
	detallesPedidoRepo ports.DetallesPedidoRepository,
	ventaRepo ports.VentaRepository,
	detallesVentaRepo ports.DetallesVentaRepository,
	ordenRepo ports.OrdenProveedorRepository,
	detallesOrdenRepo ports.DetallesOrdenRepository,
	notificationService ports.NotificationService,
	stockWS ports.WebSocketService,
	ordersWS ports.WebSocketService,
	cancellationsWS ports.WebSocketService,
) {
	// Inicializar controladores usando la fábrica
	controllerFactory := handlers.NewControllerFactory(
		productRepo,
		proveedorRepo,
		pedidoRepo,
		detallesPedidoRepo,
		ventaRepo,
		detallesVentaRepo,
		ordenRepo,
		detallesOrdenRepo,
		notificationService,
	)

	// Obtener controladores
	productoController := controllerFactory.GetProductoController()
	proveedorController := controllerFactory.GetProveedorController()
	pedidoController := controllerFactory.GetPedidoController()
	ventaController := controllerFactory.GetVentaController()
	ordenController := controllerFactory.GetOrdenProveedorController()

	// Inicializar manejadores de WebSocket
	productStockWSHandler := websocket.NewProductStockWebsocketHandler(stockWS)
	orderCreationWSHandler := websocket.NewOrderCreationWebsocketHandler(ordersWS)
	orderCancelWSHandler := websocket.NewOrderCancelWebsocketHandler(cancellationsWS)

	// WebSocket routes - Cada tipo de notificación tiene su propia ruta
	ws := engine.Group("ws")
	ws.GET("/stock", productStockWSHandler.Handle)
	ws.GET("/orders", orderCreationWSHandler.Handle)
	ws.GET("/cancellations", orderCancelWSHandler.Handle)

	// API routes
	api := engine.Group("api")

	// Rutas de productos
	productos := api.Group("productos")
	productos.GET("/", productoController.GetAll)
	productos.GET("/:id", productoController.GetByID)
	productos.POST("/", productoController.Create)
	productos.PUT("/:id", productoController.Update)
	productos.PATCH("/:id/stock", productoController.UpdateStock)
	productos.DELETE("/:id", productoController.Delete)

	// Rutas de proveedores
	proveedores := api.Group("proveedores")
	proveedores.GET("/", proveedorController.GetAll)
	proveedores.GET("/:id", proveedorController.GetByID)
	proveedores.POST("/", proveedorController.Create)
	proveedores.PUT("/:id", proveedorController.Update)
	proveedores.DELETE("/:id", proveedorController.Delete)

	// Rutas de pedidos
	pedidos := api.Group("pedidos")
	pedidos.GET("/", pedidoController.GetAll)
	pedidos.GET("/:id", pedidoController.GetByID)
	pedidos.POST("/", pedidoController.Create)
	pedidos.PUT("/:id", pedidoController.Update)
	pedidos.POST("/:id/cancelar", pedidoController.CancelPedido)
	pedidos.GET("/:id/productos", pedidoController.GetDetallesPedido)
	pedidos.POST("/:id/productos", pedidoController.AddDetallePedido)

	// Rutas de ventas
	ventas := api.Group("ventas")
	ventas.GET("/", ventaController.GetAll)
	ventas.GET("/:id", ventaController.GetByID)
	ventas.POST("/", ventaController.Create)
	ventas.PUT("/:id", ventaController.Update)
	ventas.POST("/:id/cancelar", ventaController.CancelVenta)
	ventas.GET("/:id/productos", ventaController.GetDetallesVenta)
	ventas.POST("/:id/productos", ventaController.AddDetalleVenta)

	// Rutas de órdenes de proveedor
	ordenes := api.Group("ordenes")
	ordenes.GET("/", ordenController.GetAll)
	ordenes.GET("/:id", ordenController.GetByID)
	ordenes.POST("/", ordenController.Create)
	ordenes.PUT("/:id", ordenController.Update)
	ordenes.POST("/:id/cancelar", ordenController.CancelOrden)
	ordenes.POST("/:id/recibir", ordenController.RecibirOrden)
	ordenes.GET("/:id/productos", ordenController.GetDetallesOrden)
	ordenes.POST("/:id/productos", ordenController.AddDetalleOrden)
}
