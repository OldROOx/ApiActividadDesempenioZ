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
	// Inicializar controladores
	productoController := handlers.NewProductoController(productRepo, notificationService)
	// Aquí se inicializarían los demás controladores

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

	// Aquí se definirían las rutas para las demás entidades
	// siguiendo el mismo patrón
}
