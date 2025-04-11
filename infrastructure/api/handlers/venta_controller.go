package handlers

import (
	"ActividadDesempenioAPIz/core/domain"
	"ActividadDesempenioAPIz/core/ports"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// VentaController controla las solicitudes relacionadas con ventas
type VentaController struct {
	repository          ports.VentaRepository
	detallesRepo        ports.DetallesVentaRepository
	productoRepo        ports.ProductoRepository
	notificationService ports.NotificationService
}

// NewVentaController crea un nuevo controlador de ventas
func NewVentaController(
	repository ports.VentaRepository,
	detallesRepo ports.DetallesVentaRepository,
	productoRepo ports.ProductoRepository,
	notificationService ports.NotificationService,
) *VentaController {
	return &VentaController{
		repository:          repository,
		detallesRepo:        detallesRepo,
		productoRepo:        productoRepo,
		notificationService: notificationService,
	}
}

// GetAll obtiene todas las ventas
func (vc *VentaController) GetAll(c *gin.Context) {
	ventas, err := vc.repository.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ventas)
}

// GetByID obtiene una venta por su ID
func (vc *VentaController) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	venta, err := vc.repository.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Venta no encontrada"})
		return
	}

	c.JSON(http.StatusOK, venta)
}

// Create crea una nueva venta
func (vc *VentaController) Create(c *gin.Context) {
	var venta domain.Venta
	if err := c.ShouldBindJSON(&venta); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Establecer fecha de venta
	venta.FechaVenta = time.Now()

	id, err := vc.repository.Create(&venta)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	venta.ID = id

	// Notificar la creación de la venta
	vc.notificationService.NotifyNewVenta(id, venta.Total)

	c.JSON(http.StatusCreated, venta)
}

// Update actualiza una venta existente
func (vc *VentaController) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var venta domain.Venta
	if err := c.ShouldBindJSON(&venta); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	venta.ID = id
	if err := vc.repository.Update(&venta); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, venta)
}

// CancelVenta cancela una venta
func (vc *VentaController) CancelVenta(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Obtener la venta actual
	venta, err := vc.repository.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Venta no encontrada"})
		return
	}

	// Actualizar el estado a cancelado
	err = vc.repository.UpdateEstado(id, "cancelada")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Notificar la cancelación de la venta
	vc.notificationService.NotifyCanceledVenta(id, venta.Total)

	c.JSON(http.StatusOK, gin.H{
		"message":  "Venta cancelada correctamente",
		"venta_id": id,
	})
}

// GetDetallesVenta obtiene los detalles de una venta
func (vc *VentaController) GetDetallesVenta(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	detalles, err := vc.detallesRepo.GetByVentaID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, detalles)
}

// AddDetalleVenta añade un detalle a una venta
func (vc *VentaController) AddDetalleVenta(c *gin.Context) {
	idParam := c.Param("id")
	ventaID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de venta inválido"})
		return
	}

	var detalle domain.DetallesVenta
	if err := c.ShouldBindJSON(&detalle); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	detalle.VentaID = ventaID
	detalle.Subtotal = float64(detalle.Cantidad) * detalle.PrecioUnitario

	id, err := vc.detallesRepo.Create(&detalle)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	detalle.ID = id

	// Actualizar el stock del producto
	producto, err := vc.productoRepo.GetByID(detalle.ProductoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener el producto"})
		return
	}

	nuevoStock := producto.Existencia - detalle.Cantidad
	if nuevoStock < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Stock insuficiente"})
		return
	}

	if err := vc.productoRepo.UpdateStock(detalle.ProductoID, nuevoStock); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar el stock"})
		return
	}

	// Verificar si el stock es bajo y enviar notificación
	if nuevoStock <= 5 {
		vc.notificationService.NotifyLowStock(detalle.ProductoID, nuevoStock)
	}

	c.JSON(http.StatusCreated, detalle)
}
