package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

// InitDB inicializa la conexión a la base de datos
func InitDB() {
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbPort := os.Getenv("DB_PORT")
	dbName := "ventas_online"

	if dbHost == "" {
		dbHost = "localhost"
	}

	if dbUser == "" {
		dbUser = "root"
	}

	if dbPort == "" {
		dbPort = "3306"
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	log.Printf("Intentando conectar a la base de datos: %s@%s:%s/%s", dbUser, dbHost, dbPort, dbName)

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error al abrir conexión a la base de datos: %v", err)
	}

	// Configuración del pool de conexiones
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)
	DB.SetConnMaxLifetime(5 * 60 * time.Second)

	// Intentamos varias veces conectar a la base de datos
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		err = DB.Ping()
		if err == nil {
			break
		}
		log.Printf("Intento %d: Error al hacer ping a la base de datos: %v", i+1, err)
		if i < maxRetries-1 {
			// Esperar un poco antes de reintentar
			time.Sleep(2 * time.Second)
		}
	}

	if err != nil {
		log.Fatalf("No se pudo establecer conexión con la base de datos después de %d intentos: %v", maxRetries, err)
	}

	log.Println("Conexión exitosa a la base de datos")

	// Intentamos crear la base de datos si no existe
	_, err = DB.Exec("CREATE DATABASE IF NOT EXISTS ventas_online")
	if err != nil {
		log.Printf("Error al crear la base de datos: %v", err)
	}

	// Seleccionamos la base de datos
	_, err = DB.Exec("USE ventas_online")
	if err != nil {
		log.Fatalf("Error al seleccionar la base de datos: %v", err)
	}

	// Creamos las tablas si no existen
	createTables()
}

// GetDB returns the database connection
func GetDB() *sql.DB {
	return DB
}

// createTables crea las tablas necesarias si no existen
func createTables() {
	// Tabla Proveedor
	_, err := DB.Exec(`
	CREATE TABLE IF NOT EXISTS Proveedor (
		id_proveedor INT AUTO_INCREMENT PRIMARY KEY,
		nombre VARCHAR(100) NOT NULL,
		direccion VARCHAR(200),
		telefono VARCHAR(20),
		email VARCHAR(100),
		fecha_registro DATETIME
	)`)

	if err != nil {
		log.Printf("Error al crear tabla Proveedor: %v", err)
	}

	// Tabla Producto
	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS Producto (
		id_producto INT AUTO_INCREMENT PRIMARY KEY,
		nombre VARCHAR(100) NOT NULL,
		descripcion TEXT,
		precio INT NOT NULL,
		existencia INT NOT NULL DEFAULT 0,
		id_proveedor INT,
		fecha_creacion DATETIME,
		FOREIGN KEY (id_proveedor) REFERENCES Proveedor(id_proveedor)
	)`)

	if err != nil {
		log.Printf("Error al crear tabla Producto: %v", err)
	}

	// Tabla Pedido
	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS Pedido (
		id_pedido INT AUTO_INCREMENT PRIMARY KEY,
		fecha_pedido DATETIME,
		estado VARCHAR(20) NOT NULL,
		total FLOAT NOT NULL
	)`)

	if err != nil {
		log.Printf("Error al crear tabla Pedido: %v", err)
	}

	// Tabla Detalles_Pedido
	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS Detalles_Pedido (
		id_detalle_pedido INT AUTO_INCREMENT PRIMARY KEY,
		id_pedido INT,
		id_producto INT,
		cantidad INT NOT NULL,
		precio_unitario FLOAT NOT NULL,
		subtotal FLOAT,
		FOREIGN KEY (id_pedido) REFERENCES Pedido(id_pedido),
		FOREIGN KEY (id_producto) REFERENCES Producto(id_producto)
	)`)

	if err != nil {
		log.Printf("Error al crear tabla Detalles_Pedido: %v", err)
	}

	// Tabla Venta
	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS Venta (
		id_venta INT AUTO_INCREMENT PRIMARY KEY,
		fecha_venta DATETIME,
		estado VARCHAR(20) NOT NULL,
		total FLOAT NOT NULL
	)`)

	if err != nil {
		log.Printf("Error al crear tabla Venta: %v", err)
	}

	// Tabla Detalles_Venta
	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS Detalles_Venta (
		id_detalle_venta INT AUTO_INCREMENT PRIMARY KEY,
		id_venta INT,
		id_producto INT,
		cantidad INT NOT NULL,
		precio_unitario FLOAT NOT NULL,
		subtotal FLOAT,
		FOREIGN KEY (id_venta) REFERENCES Venta(id_venta),
		FOREIGN KEY (id_producto) REFERENCES Producto(id_producto)
	)`)

	if err != nil {
		log.Printf("Error al crear tabla Detalles_Venta: %v", err)
	}

	// Tabla Orden_Proveedor
	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS Orden_Proveedor (
		id_orden_proveedor INT AUTO_INCREMENT PRIMARY KEY,
		id_proveedor INT,
		fecha_orden DATETIME,
		estado VARCHAR(20) NOT NULL,
		total INT NOT NULL,
		FOREIGN KEY (id_proveedor) REFERENCES Proveedor(id_proveedor)
	)`)

	if err != nil {
		log.Printf("Error al crear tabla Orden_Proveedor: %v", err)
	}

	// Tabla Detalles_Orden
	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS Detalles_Orden (
		id_detalle_orden INT AUTO_INCREMENT PRIMARY KEY,
		id_orden_proveedor INT,
		id_producto INT,
		cantidad INT NOT NULL,
		precio_unitario FLOAT NOT NULL,
		subtotal FLOAT,
		FOREIGN KEY (id_orden_proveedor) REFERENCES Orden_Proveedor(id_orden_proveedor),
		FOREIGN KEY (id_producto) REFERENCES Producto(id_producto)
	)`)

	if err != nil {
		log.Printf("Error al crear tabla Detalles_Orden: %v", err)
	}
}
