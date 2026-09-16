package sample

import "github.com/gin-gonic/gin"

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// @Summary		Sample handler
// @Description	Returns a sample message
// @Tags			Sample
// @Accept			json
// @Produce		json
// @Success		200	{object}	map[string]string
// @Router			/api/go-api-template/sample [get]
// @Security		Cognito
func (h *Handler) SampleHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": h.service.SampleMethod(),
	})
}
