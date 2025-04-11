package handlers

import (
	"ActividadDesempenioAPIz/core/domain"
	"ActividadDesempenioAPIz/core/ports"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// ProductoController controla las solicitudes relacionadas con productos
type ProductoController struct {
	repository          ports.ProductoRepository
	notificationService ports.NotificationService
}

// NewProductoController crea un nuevo controlador de productos
func NewProductoController(
	repository ports.ProductoRepository,
	notificationService ports.NotificationService,
) *ProductoController {
	return &ProductoController{
		repository:          repository,
		notificationService: notificationService,
	}
}

// GetAll obtiene todos los productos
func (pc *ProductoController) GetAll(c *gin.Context) {
	productos, err := pc.repository.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, productos)
}

// GetByID obtiene un producto por su ID
func (pc *ProductoController) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	producto, err := pc.repository.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Producto no encontrado"})
		return
	}

	c.JSON(http.StatusOK, producto)
}

// Create crea un nuevo producto
func (pc *ProductoController) Create(c *gin.Context) {
	var producto domain.Producto
	if err := c.ShouldBindJSON(&producto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Establecer fecha de creación
	producto.FechaCreacion = time.Now().Format("2006-01-02 15:04:05")

	id, err := pc.repository.Create(&producto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	producto.ID = id
	c.JSON(http.StatusCreated, producto)
}

// Update actualiza un producto existente
func (pc *ProductoController) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var producto domain.Producto
	if err := c.ShouldBindJSON(&producto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	producto.ID = id
	if err := pc.repository.Update(&producto); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, producto)
}

// UpdateStock actualiza el stock de un producto
func (pc *ProductoController) UpdateStock(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var stockData struct {
		Stock int `json:"stock" binding:"required"`
	}

	if err := c.ShouldBindJSON(&stockData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := pc.repository.UpdateStock(id, stockData.Stock); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Verificar si el stock es bajo y enviar notificación
	if stockData.Stock <= 5 {
		pc.notificationService.NotifyLowStock(id, stockData.Stock)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Stock actualizado correctamente",
		"producto_id": id,
		"stock":       stockData.Stock,
	})
}

// Delete elimina un producto
func (pc *ProductoController) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	if err := pc.repository.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Producto eliminado correctamente"})
}

// Aquí seguirían los demás controladores como ProveedorController, PedidoController, etc.
// con estructura similar pero adaptada a cada entidad
