package database

import (
	"log"
	"time"
)

// SetupTestData inserta datos de prueba en la base de datos
func SetupTestData() {
	log.Println("Inicializando datos de prueba...")

	// Insertar proveedores
	insertProveedores()

	// Insertar productos
	insertProductos()

	// Insertar pedidos y sus detalles
	insertPedidos()

	// Insertar ventas y sus detalles
	insertVentas()

	// Insertar órdenes de proveedor y sus detalles
	insertOrdenesProveedor()

	log.Println("Datos de prueba inicializados correctamente")
}

// insertProveedores inserta proveedores de prueba
func insertProveedores() {
	db := GetDB()

	proveedores := []struct {
		nombre    string
		direccion string
		telefono  string
		email     string
	}{
		{"Electrónica Global", "Av. Tecnología 123", "555-1234", "contacto@electronicaglobal.com"},
		{"Muebles Premium", "Calle Decoración 456", "555-5678", "ventas@mueblespremium.com"},
		{"Alimentos Frescos", "Blvd. Nutrición 789", "555-9012", "info@alimentosfrescos.com"},
		{"Textiles Modernos", "Plaza Fashion 321", "555-3456", "ventas@textilesmodernos.com"},
	}

	for _, p := range proveedores {
		_, err := db.Exec(
			"INSERT INTO Proveedor (nombre, direccion, telefono, email, fecha_registro) VALUES (?, ?, ?, ?, ?)",
			p.nombre, p.direccion, p.telefono, p.email, time.Now().Format("2006-01-02 15:04:05"),
		)
		if err != nil {
			log.Printf("Error al insertar proveedor %s: %v", p.nombre, err)
		}
	}
}

// insertProductos inserta productos de prueba
func insertProductos() {
	db := GetDB()

	productos := []struct {
		nombre      string
		descripcion string
		precio      int
		existencia  int
		proveedorID int
	}{
		{"Laptop Pro", "Laptop de última generación", 15000, 20, 1},
		{"Smartphone X", "Teléfono inteligente", 8000, 30, 1},
		{"Mesa de Centro", "Mesa de centro de madera", 3500, 10, 2},
		{"Silla Ergonómica", "Silla para oficina", 2500, 15, 2},
		{"Frutas Mixtas", "Pack de frutas variadas", 150, 50, 3},
		{"Verduras Orgánicas", "Pack de verduras orgánicas", 200, 40, 3},
		{"Camisa Casual", "Camisa de algodón", 600, 25, 4},
		{"Pantalón Formal", "Pantalón de vestir", 800, 20, 4},
		{"Tablet Mini", "Tablet compacta", 3000, 4, 1}, // Producto con stock bajo
	}

	for _, p := range productos {
		_, err := db.Exec(
			"INSERT INTO Producto (nombre, descripcion, precio, existencia, id_proveedor, fecha_creacion) VALUES (?, ?, ?, ?, ?, ?)",
			p.nombre, p.descripcion, p.precio, p.existencia, p.proveedorID, time.Now().Format("2006-01-02 15:04:05"),
		)
		if err != nil {
			log.Printf("Error al insertar producto %s: %v", p.nombre, err)
		}
	}
}

// insertPedidos inserta pedidos de prueba y sus detalles
func insertPedidos() {
	db := GetDB()

	// Insertar pedidos
	pedidos := []struct {
		estado string
		total  float64
	}{
		{"pendiente", 18000.0},
		{"completado", 3500.0},
		{"cancelado", 5000.0},
	}

	for _, p := range pedidos {
		result, err := db.Exec(
			"INSERT INTO Pedido (fecha_pedido, estado, total) VALUES (?, ?, ?)",
			time.Now().Format("2006-01-02 15:04:05"), p.estado, p.total,
		)
		if err != nil {
			log.Printf("Error al insertar pedido: %v", err)
			continue
		}

		pedidoID, _ := result.LastInsertId()

		// Insertar detalles para este pedido
		if pedidoID == 1 {
			// Detalles para el primer pedido
			_, err = db.Exec(
				"INSERT INTO Detalles_Pedido (id_pedido, id_producto, cantidad, precio_unitario) VALUES (?, ?, ?, ?)",
				pedidoID, 1, 1, 15000.0,
			)
			if err != nil {
				log.Printf("Error al insertar detalle de pedido: %v", err)
			}

			_, err = db.Exec(
				"INSERT INTO Detalles_Pedido (id_pedido, id_producto, cantidad, precio_unitario) VALUES (?, ?, ?, ?)",
				pedidoID, 2, 1, 3000.0,
			)
			if err != nil {
				log.Printf("Error al insertar detalle de pedido: %v", err)
			}
		} else if pedidoID == 2 {
			// Detalles para el segundo pedido
			_, err = db.Exec(
				"INSERT INTO Detalles_Pedido (id_pedido, id_producto, cantidad, precio_unitario) VALUES (?, ?, ?, ?)",
				pedidoID, 3, 1, 3500.0,
			)
			if err != nil {
				log.Printf("Error al insertar detalle de pedido: %v", err)
			}
		} else if pedidoID == 3 {
			// Detalles para el tercer pedido
			_, err = db.Exec(
				"INSERT INTO Detalles_Pedido (id_pedido, id_producto, cantidad, precio_unitario) VALUES (?, ?, ?, ?)",
				pedidoID, 4, 2, 2500.0,
			)
			if err != nil {
				log.Printf("Error al insertar detalle de pedido: %v", err)
			}
		}
	}
}

// insertVentas inserta ventas de prueba y sus detalles
func insertVentas() {
	db := GetDB()

	// Insertar ventas
	ventas := []struct {
		estado string
		total  float64
	}{
		{"completada", 8000.0},
		{"completada", 5000.0},
		{"cancelada", 3000.0},
	}

	for _, v := range ventas {
		result, err := db.Exec(
			"INSERT INTO Venta (fecha_venta, estado, total) VALUES (?, ?, ?)",
			time.Now(), v.estado, v.total,
		)
		if err != nil {
			log.Printf("Error al insertar venta: %v", err)
			continue
		}

		ventaID, _ := result.LastInsertId()

		// Insertar detalles para esta venta
		if ventaID == 1 {
			// Detalles para la primera venta
			_, err = db.Exec(
				"INSERT INTO Detalles_Venta (id_venta, id_producto, cantidad, precio_unitario) VALUES (?, ?, ?, ?)",
				ventaID, 2, 1, 8000.0,
			)
			if err != nil {
				log.Printf("Error al insertar detalle de venta: %v", err)
			}
		} else if ventaID == 2 {
			// Detalles para la segunda venta
			_, err = db.Exec(
				"INSERT INTO Detalles_Venta (id_venta, id_producto, cantidad, precio_unitario) VALUES (?, ?, ?, ?)",
				ventaID, 4, 2, 2500.0,
			)
			if err != nil {
				log.Printf("Error al insertar detalle de venta: %v", err)
			}
		} else if ventaID == 3 {
			// Detalles para la tercera venta
			_, err = db.Exec(
				"INSERT INTO Detalles_Venta (id_venta, id_producto, cantidad, precio_unitario) VALUES (?, ?, ?, ?)",
				ventaID, 9, 1, 3000.0,
			)
			if err != nil {
				log.Printf("Error al insertar detalle de venta: %v", err)
			}
		}
	}
}

// insertOrdenesProveedor inserta órdenes de proveedor de prueba y sus detalles
func insertOrdenesProveedor() {
	db := GetDB()

	// Insertar órdenes de proveedor
	ordenes := []struct {
		proveedorID int
		estado      string
		total       int
	}{
		{1, "pendiente", 30000},
		{2, "recibida", 12500},
		{3, "cancelada", 3500},
	}

	for _, o := range ordenes {
		result, err := db.Exec(
			"INSERT INTO Orden_Proveedor (id_proveedor, fecha_orden, estado, total) VALUES (?, ?, ?, ?)",
			o.proveedorID, time.Now().Format("2006-01-02 15:04:05"), o.estado, o.total,
		)
		if err != nil {
			log.Printf("Error al insertar orden de proveedor: %v", err)
			continue
		}

		ordenID, _ := result.LastInsertId()

		// Insertar detalles para esta orden
		if ordenID == 1 {
			// Detalles para la primera orden
			_, err = db.Exec(
				"INSERT INTO Detalles_Orden (id_orden_proveedor, id_producto, cantidad, precio_unitario) VALUES (?, ?, ?, ?)",
				ordenID, 1, 2, 15000.0,
			)
			if err != nil {
				log.Printf("Error al insertar detalle de orden: %v", err)
			}
		} else if ordenID == 2 {
			// Detalles para la segunda orden
			_, err = db.Exec(
				"INSERT INTO Detalles_Orden (id_orden_proveedor, id_producto, cantidad, precio_unitario) VALUES (?, ?, ?, ?)",
				ordenID, 4, 5, 2500.0,
			)
			if err != nil {
				log.Printf("Error al insertar detalle de orden: %v", err)
			}
		} else if ordenID == 3 {
			// Detalles para la tercera orden
			_, err = db.Exec(
				"INSERT INTO Detalles_Orden (id_orden_proveedor, id_producto, cantidad, precio_unitario) VALUES (?, ?, ?, ?)",
				ordenID, 5, 20, 175.0,
			)
			if err != nil {
				log.Printf("Error al insertar detalle de orden: %v", err)
			}
		}
	}
}
