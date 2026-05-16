package handler

import (
	"chemistry-coach/internal/usecase"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	uc *usecase.AuthUseCase
}

func NewAuthHandler(uc *usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{uc: uc}
}

// Start godoc
// @Summary      Start onboarding / auth
// @Description  Creates anonymous user or returns existing
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body body usecase.AuthStartInput true "Onboarding data"
// @Success      200 {object} usecase.AuthStartOutput
// @Failure      400 {object} map[string]string
// @Router       /auth/start [post]
func (h *AuthHandler) Start(c *gin.Context) {
	var in usecase.AuthStartInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	if existing := c.GetHeader("X-User-Id"); existing != "" {
		in.UserID = existing
	}
	out, err := h.uc.Start(c.Request.Context(), in)
	if err != nil {
		if errors.Is(err, usecase.ErrUnderage) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "underage"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, out)
}
