package handler

import (
	"chemistry-coach/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CatalogHandler struct {
	uc *usecase.CatalogUseCase
}

func NewCatalogHandler(uc *usecase.CatalogUseCase) *CatalogHandler {
	return &CatalogHandler{uc: uc}
}

// Goals godoc
// @Summary      List goals
// @Tags         catalog
// @Produce      json
// @Success      200 {object} map[string]interface{}
// @Router       /goals [get]
func (h *CatalogHandler) Goals(c *gin.Context) {
	rec := c.Query("recommendedGoalId")
	c.JSON(http.StatusOK, gin.H{"goals": h.uc.Goals(rec)})
}

// Personas godoc
// @Summary      List personas
// @Tags         catalog
// @Produce      json
// @Success      200 {object} map[string]interface{}
// @Router       /personas [get]
func (h *CatalogHandler) Personas(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"personas": h.uc.Personas()})
}
