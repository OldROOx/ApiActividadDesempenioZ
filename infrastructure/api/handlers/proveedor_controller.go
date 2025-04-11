package handlers

import (
	"ActividadDesempenioAPIz/core/domain"
	"ActividadDesempenioAPIz/core/ports"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// ProveedorController controla las solicitudes relacionadas con proveedores
type ProveedorController struct {
	repository ports.ProveedorRepository
}

// NewProveedorController crea un nuevo controlador de proveedores
func NewProveedorController(repository ports.ProveedorRepository) *ProveedorController {
	return &ProveedorController{
		repository: repository,
	}
}

// GetAll obtiene todos los proveedores
func (pc *ProveedorController) GetAll(c *gin.Context) {
	proveedores, err := pc.repository.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, proveedores)
}

// GetByID obtiene un proveedor por su ID
func (pc *ProveedorController) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	proveedor, err := pc.repository.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Proveedor no encontrado"})
		return
	}

	c.JSON(http.StatusOK, proveedor)
}

// Create crea un nuevo proveedor
func (pc *ProveedorController) Create(c *gin.Context) {
	var proveedor domain.Proveedor
	if err := c.ShouldBindJSON(&proveedor); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Establecer fecha de registro
	proveedor.FechaRegistro = time.Now().Format("2006-01-02 15:04:05")

	id, err := pc.repository.Create(&proveedor)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	proveedor.ID = id
	c.JSON(http.StatusCreated, proveedor)
}

// Update actualiza un proveedor existente
func (pc *ProveedorController) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var proveedor domain.Proveedor
	if err := c.ShouldBindJSON(&proveedor); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	proveedor.ID = id
	if err := pc.repository.Update(&proveedor); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, proveedor)
}

// Delete elimina un proveedor
func (pc *ProveedorController) Delete(c *gin.Context) {
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

	c.JSON(http.StatusOK, gin.H{"message": "Proveedor eliminado correctamente"})
}
