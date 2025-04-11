package main

import (
	"ActividadDesempenioAPIz/application"
	"ActividadDesempenioAPIz/infrastructure/api/routes"
	"ActividadDesempenioAPIz/infrastructure/database"
	"ActividadDesempenioAPIz/infrastructure/websocket"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Cargar variables de entorno
	err := godotenv.Load()
	if err != nil {
		log.Println("No se encontró archivo .env, usando variables de entorno del sistema")
	}

	// Inicializar base de datos
	database.InitDB()
	db := database.GetDB()

	// Inicializar repositorios
	productoRepo := database.NewSQLProductoRepository(db)
	proveedorRepo := database.NewSQLProveedorRepository(db)
	pedidoRepo := database.NewSQLPedidoRepository(db)
	detallesPedidoRepo := database.NewSQLDetallesPedidoRepository(db)
	ventaRepo := database.NewSQLVentaRepository(db)
	detallesVentaRepo := database.NewSQLDetallesVentaRepository(db)
	ordenRepo := database.NewSQLOrdenProveedorRepository(db)
	detallesOrdenRepo := database.NewSQLDetallesOrdenRepository(db)

	// Inicializar servicios WebSocket
	stockWS := websocket.NewWebsocketService()
	ordersWS := websocket.NewWebsocketService()
	cancellationsWS := websocket.NewWebsocketService()

	// Inicializar servicio de notificaciones
	notificationService := application.NewNotificationServiceExtended(
		stockWS,
		ordersWS,
		cancellationsWS,
		proveedorRepo,
	)

	// Configurar Gin
	r := gin.Default()

	// Configurar CORS
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	r.Use(cors.New(config))

	// Configurar rutas
	routes.SetupRouter(
		r,
		productoRepo,
		proveedorRepo,
		pedidoRepo,
		detallesPedidoRepo,
		ventaRepo,
		detallesVentaRepo,
		ordenRepo,
		detallesOrdenRepo,
		notificationService,
		stockWS,
		ordersWS,
		cancellationsWS,
	)

	// Inicializar datos de prueba
	go func() {
		// Esperamos un poco para que la base de datos esté lista
		time.Sleep(2 * time.Second)
		// Asegúrate de implementar esta función en database
		database.SetupTestData()
	}()

	// Obtener puerto
	port := os.Getenv("PORT")
	if port == "" {
		port = "4000"
	}

	// Iniciar servidor
	log.Printf("Servidor iniciado en el puerto %s", port)
	err = r.Run(":" + port)
	if err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
