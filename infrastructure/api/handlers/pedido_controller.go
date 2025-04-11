package handlers

import (
	"ActividadDesempenioAPIz/core/domain"
	"ActividadDesempenioAPIz/core/ports"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// PedidoController controla las solicitudes relacionadas con pedidos
type PedidoController struct {
	repository          ports.PedidoRepository
	detallesRepo        ports.DetallesPedidoRepository
	productoRepo        ports.ProductoRepository
	notificationService ports.NotificationService
}

// NewPedidoController crea un nuevo controlador de pedidos
func NewPedidoController(
	repository ports.PedidoRepository,
	detallesRepo ports.DetallesPedidoRepository,
	productoRepo ports.ProductoRepository,
	notificationService ports.NotificationService,
) *PedidoController {
	return &PedidoController{
		repository:          repository,
		detallesRepo:        detallesRepo,
		productoRepo:        productoRepo,
		notificationService: notificationService,
	}
}

// GetAll obtiene todos los pedidos
func (pc *PedidoController) GetAll(c *gin.Context) {
	pedidos, err := pc.repository.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pedidos)
}

// GetByID obtiene un pedido por su ID
func (pc *PedidoController) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	pedido, err := pc.repository.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pedido no encontrado"})
		return
	}

	c.JSON(http.StatusOK, pedido)
}

// Create crea un nuevo pedido
func (pc *PedidoController) Create(c *gin.Context) {
	var pedido domain.Pedido
	if err := c.ShouldBindJSON(&pedido); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Establecer fecha del pedido
	pedido.FechaPedido = time.Now().Format("2006-01-02 15:04:05")

	id, err := pc.repository.Create(&pedido)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	pedido.ID = id

	// Notificar la creación del pedido
	pc.notificationService.NotifyNewPedido(id, pedido.Total)

	c.JSON(http.StatusCreated, pedido)
}

// Update actualiza un pedido existente
func (pc *PedidoController) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var pedido domain.Pedido
	if err := c.ShouldBindJSON(&pedido); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pedido.ID = id
	if err := pc.repository.Update(&pedido); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pedido)
}

// CancelPedido cancela un pedido
func (pc *PedidoController) CancelPedido(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Obtener el pedido actual
	pedido, err := pc.repository.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pedido no encontrado"})
		return
	}

	// Actualizar el estado a cancelado
	err = pc.repository.UpdateEstado(id, "cancelado")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Notificar la cancelación del pedido
	pc.notificationService.NotifyCanceledPedido(id, pedido.Total)

	c.JSON(http.StatusOK, gin.H{
		"message":   "Pedido cancelado correctamente",
		"pedido_id": id,
	})
}

// GetDetallesPedido obtiene los detalles de un pedido
func (pc *PedidoController) GetDetallesPedido(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	detalles, err := pc.detallesRepo.GetByPedidoID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, detalles)
}

// AddDetallePedido añade un detalle a un pedido
func (pc *PedidoController) AddDetallePedido(c *gin.Context) {
	idParam := c.Param("id")
	pedidoID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de pedido inválido"})
		return
	}

	var detalle domain.DetallesPedido
	if err := c.ShouldBindJSON(&detalle); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	detalle.PedidoID = pedidoID
	detalle.Subtotal = float64(detalle.Cantidad) * detalle.PrecioUnitario

	id, err := pc.detallesRepo.Create(&detalle)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	detalle.ID = id

	// Actualizar el stock del producto
	producto, err := pc.productoRepo.GetByID(detalle.ProductoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener el producto"})
		return
	}

	nuevoStock := producto.Existencia - detalle.Cantidad
	if nuevoStock < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Stock insuficiente"})
		return
	}

	if err := pc.productoRepo.UpdateStock(detalle.ProductoID, nuevoStock); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar el stock"})
		return
	}

	// Verificar si el stock es bajo y enviar notificación
	if nuevoStock <= 5 {
		pc.notificationService.NotifyLowStock(detalle.ProductoID, nuevoStock)
	}

	c.JSON(http.StatusCreated, detalle)
}
