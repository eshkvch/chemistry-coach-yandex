package handler

import (
	"chemistry-coach/internal/usecase"
	"net/http"

	"chemistry-coach/internal/delivery/http/middleware"
	"github.com/gin-gonic/gin"
)

type ProfileHandler struct {
	uc *usecase.ProfileUseCase
}

func NewProfileHandler(uc *usecase.ProfileUseCase) *ProfileHandler {
	return &ProfileHandler{uc: uc}
}

// Get godoc
// @Summary      Get user profile
// @Tags         profile
// @Produce      json
// @Param        X-User-Id header string true "User ID"
// @Success      200 {object} usecase.ProfileOutput
// @Failure      401 {object} map[string]string
// @Router       /profile [get]
func (h *ProfileHandler) Get(c *gin.Context) {
	out, err := h.uc.Get(c.Request.Context(), middleware.GetUserID(c))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, out)
}
