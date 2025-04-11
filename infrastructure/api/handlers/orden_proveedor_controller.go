package handlers

import (
	"ActividadDesempenioAPIz/core/domain"
	"ActividadDesempenioAPIz/core/ports"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// OrdenProveedorController controla las solicitudes relacionadas con órdenes de proveedor
type OrdenProveedorController struct {
	repository          ports.OrdenProveedorRepository
	detallesRepo        ports.DetallesOrdenRepository
	proveedorRepo       ports.ProveedorRepository
	productoRepo        ports.ProductoRepository
	notificationService ports.NotificationService
}

// NewOrdenProveedorController crea un nuevo controlador de órdenes de proveedor
func NewOrdenProveedorController(
	repository ports.OrdenProveedorRepository,
	detallesRepo ports.DetallesOrdenRepository,
	proveedorRepo ports.ProveedorRepository,
	productoRepo ports.ProductoRepository,
	notificationService ports.NotificationService,
) *OrdenProveedorController {
	return &OrdenProveedorController{
		repository:          repository,
		detallesRepo:        detallesRepo,
		proveedorRepo:       proveedorRepo,
		productoRepo:        productoRepo,
		notificationService: notificationService,
	}
}

// GetAll obtiene todas las órdenes de proveedor
func (oc *OrdenProveedorController) GetAll(c *gin.Context) {
	ordenes, err := oc.repository.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ordenes)
}

// GetByID obtiene una orden de proveedor por su ID
func (oc *OrdenProveedorController) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	orden, err := oc.repository.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Orden no encontrada"})
		return
	}

	c.JSON(http.StatusOK, orden)
}

// Create crea una nueva orden de proveedor
func (oc *OrdenProveedorController) Create(c *gin.Context) {
	var orden domain.OrdenProveedor
	if err := c.ShouldBindJSON(&orden); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Establecer fecha de orden
	orden.FechaOrden = time.Now().Format("2006-01-02 15:04:05")

	// Verificar que el proveedor exista
	_, err := oc.proveedorRepo.GetByID(orden.ProveedorID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Proveedor no encontrado"})
		return
	}

	id, err := oc.repository.Create(&orden)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	orden.ID = id

	// Notificar la creación de la orden
	oc.notificationService.NotifyNewOrdenProveedor(id, float64(orden.Total))

	c.JSON(http.StatusCreated, orden)
}

// Update actualiza una orden de proveedor existente
func (oc *OrdenProveedorController) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var orden domain.OrdenProveedor
	if err := c.ShouldBindJSON(&orden); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orden.ID = id
	if err := oc.repository.Update(&orden); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, orden)
}

// CancelOrden cancela una orden de proveedor
func (oc *OrdenProveedorController) CancelOrden(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Obtener la orden actual
	orden, err := oc.repository.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Orden no encontrada"})
		return
	}

	// Actualizar el estado a cancelado
	err = oc.repository.UpdateEstado(id, "cancelada")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Notificar la cancelación de la orden
	oc.notificationService.NotifyCanceledOrdenProveedor(id, float64(orden.Total), orden.ProveedorID)

	c.JSON(http.StatusOK, gin.H{
		"message":  "Orden cancelada correctamente",
		"orden_id": id,
	})
}

// GetDetallesOrden obtiene los detalles de una orden
func (oc *OrdenProveedorController) GetDetallesOrden(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	detalles, err := oc.detallesRepo.GetByOrdenID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, detalles)
}

// AddDetalleOrden añade un detalle a una orden
func (oc *OrdenProveedorController) AddDetalleOrden(c *gin.Context) {
	idParam := c.Param("id")
	ordenID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de orden inválido"})
		return
	}

	var detalle domain.DetallesOrden
	if err := c.ShouldBindJSON(&detalle); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	detalle.OrdenProveedorID = ordenID
	detalle.Subtotal = float64(detalle.Cantidad) * detalle.PrecioUnitario

	id, err := oc.detallesRepo.Create(&detalle)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	detalle.ID = id

	// Actualizar el total de la orden
	orden, err := oc.repository.GetByID(ordenID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener la orden"})
		return
	}

	// Al recibir los productos, actualizamos el stock
	if orden.Estado == "recibida" {
		producto, err := oc.productoRepo.GetByID(detalle.ProductoID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener el producto"})
			return
		}

		nuevoStock := producto.Existencia + detalle.Cantidad
		if err := oc.productoRepo.UpdateStock(detalle.ProductoID, nuevoStock); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar el stock"})
			return
		}
	}

	c.JSON(http.StatusCreated, detalle)
}

// RecibirOrden marca una orden como recibida y actualiza el stock
func (oc *OrdenProveedorController) RecibirOrden(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Obtener la orden actual
	orden, err := oc.repository.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Orden no encontrada"})
		return
	}

	// Actualizar el estado a recibida
	err = oc.repository.UpdateEstado(id, "recibida")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Obtener detalles de la orden para actualizar el stock
	detalles, err := oc.detallesRepo.GetByOrdenID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Actualizar el stock de cada producto en la orden
	for _, detalle := range detalles {
		producto, err := oc.productoRepo.GetByID(detalle.ProductoID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener el producto"})
			return
		}

		nuevoStock := producto.Existencia + detalle.Cantidad
		if err := oc.productoRepo.UpdateStock(detalle.ProductoID, nuevoStock); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar el stock"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Orden recibida correctamente",
		"orden_id": id,
	})
}
