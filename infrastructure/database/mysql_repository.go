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
