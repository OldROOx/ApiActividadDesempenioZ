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
func (c *OrdenProveedorController) GetAll(ctx *gin.Context) {
	ordenes, err := c.repository.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, ordenes)
}

// GetByID obtiene una orden de proveedor por su ID
func (c *OrdenProveedorController) GetByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	orden, err := c.repository.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Orden no encontrada"})
		return
	}

	ctx.JSON(http.StatusOK, orden)
}

// Create crea una nueva orden de proveedor
func (c *OrdenProveedorController) Create(ctx *gin.Context) {
	var createOrdenRequest struct {
		ProveedorID int `json:"id_proveedor" binding:"required"`
		Detalles    []struct {
			ProductoID     int     `json:"id_producto" binding:"required"`
			Cantidad       int     `json:"cantidad" binding:"required"`
			PrecioUnitario float64 `json:"precio_unitario" binding:"required"`
		} `json:"detalles" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&createOrdenRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verificar que el proveedor existe
	_, err := c.proveedorRepo.GetByID(createOrdenRequest.ProveedorID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Proveedor no encontrado"})
		return
	}

	// Crear la orden
	orden := &domain.OrdenProveedor{
		ProveedorID: createOrdenRequest.ProveedorID,
		Estado:      "pendiente",
		FechaOrden:  time.Now().Format("2006-01-02 15:04:05"),
		Total:       0,
	}

	ordenID, err := c.repository.Create(orden)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Crear los detalles de la orden y calcular el total
	var totalFloat float64
	for _, detalleRequest := range createOrdenRequest.Detalles {
		detalle := &domain.DetallesOrden{
			OrdenProveedorID: ordenID,
			ProductoID:       detalleRequest.ProductoID,
			Cantidad:         detalleRequest.Cantidad,
			PrecioUnitario:   detalleRequest.PrecioUnitario,
			Subtotal:         float64(detalleRequest.Cantidad) * detalleRequest.PrecioUnitario,
		}

		_, err := c.detallesRepo.Create(detalle)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		totalFloat += detalle.Subtotal
	}

	// Actualizar el total de la orden
	orden.ID = ordenID
	orden.Total = int(totalFloat)
	c.repository.Update(orden)

	// Enviar notificación de nueva orden
	c.notificationService.NotifyNewOrdenProveedor(ordenID, totalFloat)

	ctx.JSON(http.StatusCreated, gin.H{
		"id_orden": ordenID,
		"total":    totalFloat,
		"mensaje":  "Orden creada correctamente",
	})
}

// Update actualiza una orden de proveedor existente
func (c *OrdenProveedorController) Update(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var orden domain.OrdenProveedor
	if err := ctx.ShouldBindJSON(&orden); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orden.ID = id
	if err := c.repository.Update(&orden); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, orden)
}

// CancelOrden cancela una orden de proveedor
func (c *OrdenProveedorController) CancelOrden(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Obtener la orden actual
	orden, err := c.repository.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Orden no encontrada"})
		return
	}

	// Actualizar el estado a cancelado
	err = c.repository.UpdateEstado(id, "cancelada")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Notificar la cancelación de la orden
	c.notificationService.NotifyCanceledOrdenProveedor(id, float64(orden.Total), orden.ProveedorID)

	ctx.JSON(http.StatusOK, gin.H{
		"message":  "Orden cancelada correctamente",
		"orden_id": id,
	})
}

// GetDetallesOrden obtiene los detalles de una orden
func (c *OrdenProveedorController) GetDetallesOrden(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	detalles, err := c.detallesRepo.GetByOrdenID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, detalles)
}

// AddDetalleOrden añade un detalle a una orden
func (c *OrdenProveedorController) AddDetalleOrden(ctx *gin.Context) {
	idParam := ctx.Param("id")
	ordenID, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID de orden inválido"})
		return
	}

	var detalle domain.DetallesOrden
	if err := ctx.ShouldBindJSON(&detalle); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	detalle.OrdenProveedorID = ordenID
	detalle.Subtotal = float64(detalle.Cantidad) * detalle.PrecioUnitario

	id, err := c.detallesRepo.Create(&detalle)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	detalle.ID = id

	// Actualizar el total de la orden
	orden, err := c.repository.GetByID(ordenID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener la orden"})
		return
	}

	orden.Total += int(detalle.Subtotal)
	c.repository.Update(orden)

	ctx.JSON(http.StatusCreated, detalle)
}

// RecibirOrden marca una orden como recibida y actualiza el inventario
func (c *OrdenProveedorController) RecibirOrden(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Obtener la orden actual
	orden, err := c.repository.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Orden no encontrada"})
		return
	}

	// Verificar que la orden está pendiente
	if orden.Estado != "pendiente" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "La orden no está en estado pendiente"})
		return
	}

	// Actualizar el estado a recibida
	err = c.repository.UpdateEstado(id, "recibida")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Obtener detalles de la orden para actualizar el stock
	detalles, err := c.detallesRepo.GetByOrdenID(id)
	if err == nil {
		for _, detalle := range detalles {
			producto, err := c.productoRepo.GetByID(detalle.ProductoID)
			if err == nil {
				nuevoStock := producto.Existencia + detalle.Cantidad
				c.productoRepo.UpdateStock(detalle.ProductoID, nuevoStock)

				// Verificar si aún hay stock bajo después de recibir
				if nuevoStock <= 5 {
					c.notificationService.NotifyLowStock(detalle.ProductoID, nuevoStock)
				}
			}
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":  "Orden recibida correctamente",
		"orden_id": id,
	})
}
