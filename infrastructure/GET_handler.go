package infrastructure

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// GetHandler proporciona acceso a los recursos HTTP del sistema
type GetHandler struct {
	baseURL string
}

// NewGetHandler crea un nuevo manejador de recursos HTTP
func NewGetHandler(baseURL string) *GetHandler {
	return &GetHandler{
		baseURL: baseURL,
	}
}

// HandleStaticResource maneja la solicitud de recursos estáticos
func (h *GetHandler) HandleStaticResource(c *gin.Context) {
	resourcePath := c.Param("resource")

	// Limpiar y validar el path del recurso
	resourcePath = strings.TrimPrefix(resourcePath, "/")
	resourcePath = strings.TrimSuffix(resourcePath, "/")

	// Determinar el tipo de contenido
	contentType := "text/plain"
	if strings.HasSuffix(resourcePath, ".html") {
		contentType = "text/html"
	} else if strings.HasSuffix(resourcePath, ".js") {
		contentType = "application/javascript"
	} else if strings.HasSuffix(resourcePath, ".css") {
		contentType = "text/css"
	} else if strings.HasSuffix(resourcePath, ".json") {
		contentType = "application/json"
	}

	// Establecer encabezados
	c.Header("Content-Type", contentType)
	c.File("./static/" + resourcePath)
}

// HandleProductImage maneja la solicitud de imágenes de productos
func (h *GetHandler) HandleProductImage(c *gin.Context) {
	idStr := c.Param("id")
	_, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Aquí podrías obtener la imagen real de una base de datos o sistema de archivos
	// Por ahora, simplemente servimos una imagen de marcador de posición
	c.File("./static/images/product_placeholder.png")
}

// NotFoundHandler maneja las solicitudes a rutas no encontradas
func (h *GetHandler) NotFoundHandler(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{
		"error": "Recurso no encontrado",
		"path":  c.Request.URL.Path,
	})
}

// HealthCheck proporciona un punto de verificación de salud para el sistema
func (h *GetHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "UP",
		"version": "1.0.0",
	})
}
