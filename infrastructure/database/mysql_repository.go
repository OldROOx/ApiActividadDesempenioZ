package database

import (
	"ActividadDesempenioAPIz/core/domain"
	"ActividadDesempenioAPIz/core/ports"
	"database/sql"
	"time"
)

// SQLProductoRepository implementa la interfaz ProductoRepository usando MySQL
type SQLProductoRepository struct {
	db *sql.DB
}

// NewSQLProductoRepository crea un nuevo repositorio de productos SQL
func NewSQLProductoRepository(db *sql.DB) ports.ProductoRepository {
	return &SQLProductoRepository{
		db: db,
	}
}

// GetByID obtiene un producto por su ID
func (r *SQLProductoRepository) GetByID(id int) (*domain.Producto, error) {
	query := `SELECT id_producto, nombre, descripcion, precio, existencia, 
              id_proveedor, fecha_creacion FROM Producto WHERE id_producto = ?`

	producto := &domain.Producto{}
	err := r.db.QueryRow(query, id).Scan(
		&producto.ID, &producto.Nombre, &producto.Descripcion,
		&producto.Precio, &producto.Existencia, &producto.ProveedorID,
		&producto.FechaCreacion,
	)

	if err != nil {
		return nil, err
	}

	return producto, nil
}

// GetAll obtiene todos los productos
func (r *SQLProductoRepository) GetAll() ([]*domain.Producto, error) {
	query := `SELECT id_producto, nombre, descripcion, precio, existencia, 
              id_proveedor, fecha_creacion FROM Producto`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	productos := []*domain.Producto{}
	for rows.Next() {
		producto := &domain.Producto{}
		err := rows.Scan(
			&producto.ID, &producto.Nombre, &producto.Descripcion,
			&producto.Precio, &producto.Existencia, &producto.ProveedorID,
			&producto.FechaCreacion,
		)
		if err != nil {
			return nil, err
		}
		productos = append(productos, producto)
	}

	return productos, nil
}

// Create crea un nuevo producto
func (r *SQLProductoRepository) Create(producto *domain.Producto) (int, error) {
	query := `INSERT INTO Producto (nombre, descripcion, precio, existencia, 
              id_proveedor, fecha_creacion) VALUES (?, ?, ?, ?, ?, ?)`

	result, err := r.db.Exec(query,
		producto.Nombre, producto.Descripcion, producto.Precio,
		producto.Existencia, producto.ProveedorID, time.Now().Format("2006-01-02 15:04:05"),
	)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// Update actualiza un producto existente
func (r *SQLProductoRepository) Update(producto *domain.Producto) error {
	query := `UPDATE Producto SET nombre = ?, descripcion = ?, precio = ?, 
              existencia = ?, id_proveedor = ? WHERE id_producto = ?`

	_, err := r.db.Exec(query,
		producto.Nombre, producto.Descripcion, producto.Precio,
		producto.Existencia, producto.ProveedorID, producto.ID,
	)

	return err
}

// UpdateStock actualiza el stock de un producto
func (r *SQLProductoRepository) UpdateStock(id int, cantidad int) error {
	query := `UPDATE Producto SET existencia = ? WHERE id_producto = ?`

	_, err := r.db.Exec(query, cantidad, id)

	return err
}

// Delete elimina un producto
func (r *SQLProductoRepository) Delete(id int) error {
	query := `DELETE FROM Producto WHERE id_producto = ?`

	_, err := r.db.Exec(query, id)

	return err
}

// SQLProveedorRepository implementa la interfaz ProveedorRepository usando MySQL
type SQLProveedorRepository struct {
	db *sql.DB
}

// NewSQLProveedorRepository crea un nuevo repositorio de proveedores SQL
func NewSQLProveedorRepository(db *sql.DB) ports.ProveedorRepository {
	return &SQLProveedorRepository{
		db: db,
	}
}

// GetByID obtiene un proveedor por su ID
func (r *SQLProveedorRepository) GetByID(id int) (*domain.Proveedor, error) {
	query := `SELECT id_proveedor, nombre, direccion, telefono, email, fecha_registro 
              FROM Proveedor WHERE id_proveedor = ?`

	proveedor := &domain.Proveedor{}
	err := r.db.QueryRow(query, id).Scan(
		&proveedor.ID, &proveedor.Nombre, &proveedor.Direccion,
		&proveedor.Telefono, &proveedor.Email, &proveedor.FechaRegistro,
	)

	if err != nil {
		return nil, err
	}

	return proveedor, nil
}

// GetAll obtiene todos los proveedores
func (r *SQLProveedorRepository) GetAll() ([]*domain.Proveedor, error) {
	query := `SELECT id_proveedor, nombre, direccion, telefono, email, fecha_registro 
              FROM Proveedor`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	proveedores := []*domain.Proveedor{}
	for rows.Next() {
		proveedor := &domain.Proveedor{}
		err := rows.Scan(
			&proveedor.ID, &proveedor.Nombre, &proveedor.Direccion,
			&proveedor.Telefono, &proveedor.Email, &proveedor.FechaRegistro,
		)
		if err != nil {
			return nil, err
		}
		proveedores = append(proveedores, proveedor)
	}

	return proveedores, nil
}

// Create crea un nuevo proveedor
func (r *SQLProveedorRepository) Create(proveedor *domain.Proveedor) (int, error) {
	query := `INSERT INTO Proveedor (nombre, direccion, telefono, email, fecha_registro) 
              VALUES (?, ?, ?, ?, ?)`

	result, err := r.db.Exec(query,
		proveedor.Nombre, proveedor.Direccion, proveedor.Telefono,
		proveedor.Email, time.Now().Format("2006-01-02 15:04:05"),
	)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// Update actualiza un proveedor existente
func (r *SQLProveedorRepository) Update(proveedor *domain.Proveedor) error {
	query := `UPDATE Proveedor SET nombre = ?, direccion = ?, telefono = ?, 
              email = ? WHERE id_proveedor = ?`

	_, err := r.db.Exec(query,
		proveedor.Nombre, proveedor.Direccion, proveedor.Telefono,
		proveedor.Email, proveedor.ID,
	)

	return err
}

// Delete elimina un proveedor
func (r *SQLProveedorRepository) Delete(id int) error {
	query := `DELETE FROM Proveedor WHERE id_proveedor = ?`

	_, err := r.db.Exec(query, id)

	return err
}

// SQLPedidoRepository implementa la interfaz PedidoRepository usando MySQL
type SQLPedidoRepository struct {
	db *sql.DB
}

// NewSQLPedidoRepository crea un nuevo repositorio de pedidos SQL
func NewSQLPedidoRepository(db *sql.DB) ports.PedidoRepository {
	return &SQLPedidoRepository{
		db: db,
	}
}

// GetByID obtiene un pedido por su ID
func (r *SQLPedidoRepository) GetByID(id int) (*domain.Pedido, error) {
	query := `SELECT id_pedido, fecha_pedido, estado, total 
              FROM Pedido WHERE id_pedido = ?`

	pedido := &domain.Pedido{}
	err := r.db.QueryRow(query, id).Scan(
		&pedido.ID, &pedido.FechaPedido, &pedido.Estado, &pedido.Total,
	)

	if err != nil {
		return nil, err
	}

	return pedido, nil
}

// GetAll obtiene todos los pedidos
func (r *SQLPedidoRepository) GetAll() ([]*domain.Pedido, error) {
	query := `SELECT id_pedido, fecha_pedido, estado, total FROM Pedido`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pedidos := []*domain.Pedido{}
	for rows.Next() {
		pedido := &domain.Pedido{}
		err := rows.Scan(
			&pedido.ID, &pedido.FechaPedido, &pedido.Estado, &pedido.Total,
		)
		if err != nil {
			return nil, err
		}
		pedidos = append(pedidos, pedido)
	}

	return pedidos, nil
}

// Create crea un nuevo pedido
func (r *SQLPedidoRepository) Create(pedido *domain.Pedido) (int, error) {
	query := `INSERT INTO Pedido (fecha_pedido, estado, total) VALUES (?, ?, ?)`

	result, err := r.db.Exec(query,
		time.Now().Format("2006-01-02 15:04:05"), pedido.Estado, pedido.Total,
	)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// Update actualiza un pedido existente
func (r *SQLPedidoRepository) Update(pedido *domain.Pedido) error {
	query := `UPDATE Pedido SET fecha_pedido = ?, estado = ?, total = ? 
              WHERE id_pedido = ?`

	_, err := r.db.Exec(query,
		pedido.FechaPedido, pedido.Estado, pedido.Total, pedido.ID,
	)

	return err
}

// UpdateEstado actualiza el estado de un pedido
func (r *SQLPedidoRepository) UpdateEstado(id int, estado string) error {
	query := `UPDATE Pedido SET estado = ? WHERE id_pedido = ?`

	_, err := r.db.Exec(query, estado, id)

	return err
}

// Delete elimina un pedido
func (r *SQLPedidoRepository) Delete(id int) error {
	query := `DELETE FROM Pedido WHERE id_pedido = ?`

	_, err := r.db.Exec(query, id)

	return err
}

// SQLDetallesPedidoRepository implementa la interfaz DetallesPedidoRepository usando MySQL
type SQLDetallesPedidoRepository struct {
	db *sql.DB
}

// NewSQLDetallesPedidoRepository crea un nuevo repositorio de detalles de pedido SQL
func NewSQLDetallesPedidoRepository(db *sql.DB) ports.DetallesPedidoRepository {
	return &SQLDetallesPedidoRepository{
		db: db,
	}
}

// GetByPedidoID obtiene los detalles de un pedido por ID del pedido
func (r *SQLDetallesPedidoRepository) GetByPedidoID(pedidoID int) ([]*domain.DetallesPedido, error) {
	query := `SELECT id_detalle_pedido, id_pedido, id_producto, cantidad, 
              precio_unitario, subtotal FROM Detalles_Pedido WHERE id_pedido = ?`

	rows, err := r.db.Query(query, pedidoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	detalles := []*domain.DetallesPedido{}
	for rows.Next() {
		detalle := &domain.DetallesPedido{}
		err := rows.Scan(
			&detalle.ID, &detalle.PedidoID, &detalle.ProductoID,
			&detalle.Cantidad, &detalle.PrecioUnitario, &detalle.Subtotal,
		)
		if err != nil {
			return nil, err
		}
		detalles = append(detalles, detalle)
	}

	return detalles, nil
}

// Create crea un nuevo detalle de pedido
func (r *SQLDetallesPedidoRepository) Create(detalle *domain.DetallesPedido) (int, error) {
	query := `INSERT INTO Detalles_Pedido (id_pedido, id_producto, cantidad, precio_unitario) 
              VALUES (?, ?, ?, ?)`

	result, err := r.db.Exec(query,
		detalle.PedidoID, detalle.ProductoID, detalle.Cantidad, detalle.PrecioUnitario,
	)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// Update actualiza un detalle de pedido existente
func (r *SQLDetallesPedidoRepository) Update(detalle *domain.DetallesPedido) error {
	query := `UPDATE Detalles_Pedido SET id_pedido = ?, id_producto = ?, 
              cantidad = ?, precio_unitario = ? WHERE id_detalle_pedido = ?`

	_, err := r.db.Exec(query,
		detalle.PedidoID, detalle.ProductoID, detalle.Cantidad,
		detalle.PrecioUnitario, detalle.ID,
	)

	return err
}

// Delete elimina un detalle de pedido
func (r *SQLDetallesPedidoRepository) Delete(id int) error {
	query := `DELETE FROM Detalles_Pedido WHERE id_detalle_pedido = ?`

	_, err := r.db.Exec(query, id)

	return err
}

// SQLVentaRepository implementa la interfaz VentaRepository usando MySQL
type SQLVentaRepository struct {
	db *sql.DB
}

// NewSQLVentaRepository crea un nuevo repositorio de ventas SQL
func NewSQLVentaRepository(db *sql.DB) ports.VentaRepository {
	return &SQLVentaRepository{
		db: db,
	}
}

// GetByID obtiene una venta por su ID
func (r *SQLVentaRepository) GetByID(id int) (*domain.Venta, error) {
	query := `SELECT id_venta, fecha_venta, estado, total 
              FROM Venta WHERE id_venta = ?`

	venta := &domain.Venta{}
	err := r.db.QueryRow(query, id).Scan(
		&venta.ID, &venta.FechaVenta, &venta.Estado, &venta.Total,
	)

	if err != nil {
		return nil, err
	}

	return venta, nil
}

// GetAll obtiene todas las ventas
func (r *SQLVentaRepository) GetAll() ([]*domain.Venta, error) {
	query := `SELECT id_venta, fecha_venta, estado, total FROM Venta`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ventas := []*domain.Venta{}
	for rows.Next() {
		venta := &domain.Venta{}
		err := rows.Scan(
			&venta.ID, &venta.FechaVenta, &venta.Estado, &venta.Total,
		)
		if err != nil {
			return nil, err
		}
		ventas = append(ventas, venta)
	}

	return ventas, nil
}

// Create crea una nueva venta
func (r *SQLVentaRepository) Create(venta *domain.Venta) (int, error) {
	query := `INSERT INTO Venta (fecha_venta, estado, total) VALUES (?, ?, ?)`

	result, err := r.db.Exec(query,
		venta.FechaVenta, venta.Estado, venta.Total,
	)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// Update actualiza una venta existente
func (r *SQLVentaRepository) Update(venta *domain.Venta) error {
	query := `UPDATE Venta SET fecha_venta = ?, estado = ?, total = ? 
              WHERE id_venta = ?`

	_, err := r.db.Exec(query,
		venta.FechaVenta, venta.Estado, venta.Total, venta.ID,
	)

	return err
}

// UpdateEstado actualiza el estado de una venta
func (r *SQLVentaRepository) UpdateEstado(id int, estado string) error {
	query := `UPDATE Venta SET estado = ? WHERE id_venta = ?`

	_, err := r.db.Exec(query, estado, id)

	return err
}

// Delete elimina una venta
func (r *SQLVentaRepository) Delete(id int) error {
	query := `DELETE FROM Venta WHERE id_venta = ?`

	_, err := r.db.Exec(query, id)

	return err
}

// SQLDetallesVentaRepository implementa la interfaz DetallesVentaRepository usando MySQL
type SQLDetallesVentaRepository struct {
	db *sql.DB
}

// NewSQLDetallesVentaRepository crea un nuevo repositorio de detalles de venta SQL
func NewSQLDetallesVentaRepository(db *sql.DB) ports.DetallesVentaRepository {
	return &SQLDetallesVentaRepository{
		db: db,
	}
}

// GetByVentaID obtiene los detalles de una venta por ID de la venta
func (r *SQLDetallesVentaRepository) GetByVentaID(ventaID int) ([]*domain.DetallesVenta, error) {
	query := `SELECT id_detalle_venta, id_venta, id_producto, cantidad, 
              precio_unitario, subtotal FROM Detalles_Venta WHERE id_venta = ?`

	rows, err := r.db.Query(query, ventaID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	detalles := []*domain.DetallesVenta{}
	for rows.Next() {
		detalle := &domain.DetallesVenta{}
		err := rows.Scan(
			&detalle.ID, &detalle.VentaID, &detalle.ProductoID,
			&detalle.Cantidad, &detalle.PrecioUnitario, &detalle.Subtotal,
		)
		if err != nil {
			return nil, err
		}
		detalles = append(detalles, detalle)
	}

	return detalles, nil
}

// Create crea un nuevo detalle de venta
func (r *SQLDetallesVentaRepository) Create(detalle *domain.DetallesVenta) (int, error) {
	query := `INSERT INTO Detalles_Venta (id_venta, id_producto, cantidad, precio_unitario) 
              VALUES (?, ?, ?, ?)`

	result, err := r.db.Exec(query,
		detalle.VentaID, detalle.ProductoID, detalle.Cantidad, detalle.PrecioUnitario,
	)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// Update actualiza un detalle de venta existente
func (r *SQLDetallesVentaRepository) Update(detalle *domain.DetallesVenta) error {
	query := `UPDATE Detalles_Venta SET id_venta = ?, id_producto = ?, 
              cantidad = ?, precio_unitario = ? WHERE id_detalle_venta = ?`

	_, err := r.db.Exec(query,
		detalle.VentaID, detalle.ProductoID, detalle.Cantidad,
		detalle.PrecioUnitario, detalle.ID,
	)

	return err
}

// Delete elimina un detalle de venta
func (r *SQLDetallesVentaRepository) Delete(id int) error {
	query := `DELETE FROM Detalles_Venta WHERE id_detalle_venta = ?`

	_, err := r.db.Exec(query, id)

	return err
}

// SQLOrdenProveedorRepository implementa la interfaz OrdenProveedorRepository usando MySQL
type SQLOrdenProveedorRepository struct {
	db *sql.DB
}

// NewSQLOrdenProveedorRepository crea un nuevo repositorio de órdenes de proveedor SQL
func NewSQLOrdenProveedorRepository(db *sql.DB) ports.OrdenProveedorRepository {
	return &SQLOrdenProveedorRepository{
		db: db,
	}
}

// GetByID obtiene una orden de proveedor por su ID
func (r *SQLOrdenProveedorRepository) GetByID(id int) (*domain.OrdenProveedor, error) {
	query := `SELECT id_orden_proveedor, id_proveedor, fecha_orden, estado, total 
              FROM Orden_Proveedor WHERE id_orden_proveedor = ?`

	orden := &domain.OrdenProveedor{}
	err := r.db.QueryRow(query, id).Scan(
		&orden.ID, &orden.ProveedorID, &orden.FechaOrden, &orden.Estado, &orden.Total,
	)

	if err != nil {
		return nil, err
	}

	return orden, nil
}

// GetAll obtiene todas las órdenes de proveedor
func (r *SQLOrdenProveedorRepository) GetAll() ([]*domain.OrdenProveedor, error) {
	query := `SELECT id_orden_proveedor, id_proveedor, fecha_orden, estado, total 
              FROM Orden_Proveedor`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ordenes := []*domain.OrdenProveedor{}
	for rows.Next() {
		orden := &domain.OrdenProveedor{}
		err := rows.Scan(
			&orden.ID, &orden.ProveedorID, &orden.FechaOrden, &orden.Estado, &orden.Total,
		)
		if err != nil {
			return nil, err
		}
		ordenes = append(ordenes, orden)
	}

	return ordenes, nil
}

// Create crea una nueva orden de proveedor
func (r *SQLOrdenProveedorRepository) Create(orden *domain.OrdenProveedor) (int, error) {
	query := `INSERT INTO Orden_Proveedor (id_proveedor, fecha_orden, estado, total) 
              VALUES (?, ?, ?, ?)`

	result, err := r.db.Exec(query,
		orden.ProveedorID, time.Now().Format("2006-01-02 15:04:05"), orden.Estado, orden.Total,
	)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// Update actualiza una orden de proveedor existente
func (r *SQLOrdenProveedorRepository) Update(orden *domain.OrdenProveedor) error {
	query := `UPDATE Orden_Proveedor SET id_proveedor = ?, fecha_orden = ?, 
              estado = ?, total = ? WHERE id_orden_proveedor = ?`

	_, err := r.db.Exec(query,
		orden.ProveedorID, orden.FechaOrden, orden.Estado, orden.Total, orden.ID,
	)

	return err
}

// UpdateEstado actualiza el estado de una orden de proveedor
func (r *SQLOrdenProveedorRepository) UpdateEstado(id int, estado string) error {
	query := `UPDATE Orden_Proveedor SET estado = ? WHERE id_orden_proveedor = ?`

	_, err := r.db.Exec(query, estado, id)

	return err
}

// Delete elimina una orden de proveedor
func (r *SQLOrdenProveedorRepository) Delete(id int) error {
	query := `DELETE FROM Orden_Proveedor WHERE id_orden_proveedor = ?`

	_, err := r.db.Exec(query, id)

	return err
}

// SQLDetallesOrdenRepository implementa la interfaz DetallesOrdenRepository usando MySQL
type SQLDetallesOrdenRepository struct {
	db *sql.DB
}

// NewSQLDetallesOrdenRepository crea un nuevo repositorio de detalles de orden SQL
func NewSQLDetallesOrdenRepository(db *sql.DB) ports.DetallesOrdenRepository {
	return &SQLDetallesOrdenRepository{
		db: db,
	}
}

// GetByOrdenID obtiene los detalles de una orden por ID de la orden
func (r *SQLDetallesOrdenRepository) GetByOrdenID(ordenID int) ([]*domain.DetallesOrden, error) {
	query := `SELECT id_detalle_orden, id_orden_proveedor, id_producto, cantidad, 
              precio_unitario, subtotal FROM Detalles_Orden WHERE id_orden_proveedor = ?`

	rows, err := r.db.Query(query, ordenID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	detalles := []*domain.DetallesOrden{}
	for rows.Next() {
		detalle := &domain.DetallesOrden{}
		err := rows.Scan(
			&detalle.ID, &detalle.OrdenProveedorID, &detalle.ProductoID,
			&detalle.Cantidad, &detalle.PrecioUnitario, &detalle.Subtotal,
		)
		if err != nil {
			return nil, err
		}
		detalles = append(detalles, detalle)
	}

	return detalles, nil
}

// Create crea un nuevo detalle de orden
func (r *SQLDetallesOrdenRepository) Create(detalle *domain.DetallesOrden) (int, error) {
	query := `INSERT INTO Detalles_Orden (id_orden_proveedor, id_producto, cantidad, precio_unitario) 
              VALUES (?, ?, ?, ?)`

	result, err := r.db.Exec(query,
		detalle.OrdenProveedorID, detalle.ProductoID, detalle.Cantidad, detalle.PrecioUnitario,
	)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// Update actualiza un detalle de orden existente
func (r *SQLDetallesOrdenRepository) Update(detalle *domain.DetallesOrden) error {
	query := `UPDATE Detalles_Orden SET id_orden_proveedor = ?, id_producto = ?, 
              cantidad = ?, precio_unitario = ? WHERE id_detalle_orden = ?`

	_, err := r.db.Exec(query,
		detalle.OrdenProveedorID, detalle.ProductoID, detalle.Cantidad,
		detalle.PrecioUnitario, detalle.ID,
	)

	return err
}

// Delete elimina un detalle de orden
func (r *SQLDetallesOrdenRepository) Delete(id int) error {
	query := `DELETE FROM Detalles_Orden WHERE id_detalle_orden = ?`

	_, err := r.db.Exec(query, id)

	return err
}
