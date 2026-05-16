package handler

import (
	"errors"
	"net/http"

	"chemistry-coach/internal/delivery/http/middleware"
	"chemistry-coach/internal/usecase"

	"github.com/gin-gonic/gin"
)

type SessionHandler struct {
	uc *usecase.SessionUseCase
}

func NewSessionHandler(uc *usecase.SessionUseCase) *SessionHandler {
	return &SessionHandler{uc: uc}
}

// Create godoc
// @Summary      Create session
// @Tags         sessions
// @Accept       json
// @Produce      json
// @Param        X-User-Id header string true "User ID"
// @Param        body body usecase.CreateSessionInput true "Session"
// @Success      201 {object} usecase.CreateSessionOutput
// @Router       /sessions [post]
func (h *SessionHandler) Create(c *gin.Context) {
	var in usecase.CreateSessionInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	out, err := h.uc.Create(c.Request.Context(), middleware.GetUserID(c), in)
	if err != nil {
		writeSessionError(c, err)
		return
	}
	c.JSON(http.StatusCreated, out)
}

// SendMessage godoc
// @Summary      Send chat message
// @Tags         sessions
// @Accept       json
// @Produce      json
// @Param        X-User-Id header string true "User ID"
// @Param        sessionId path string true "Session ID"
// @Param        body body usecase.SendMessageInput true "Message"
// @Success      200 {object} usecase.SendMessageOutput
// @Router       /sessions/{sessionId}/messages [post]
func (h *SessionHandler) SendMessage(c *gin.Context) {
	var in usecase.SendMessageInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	out, err := h.uc.SendMessage(c.Request.Context(), middleware.GetUserID(c), c.Param("sessionId"), in)
	if err != nil {
		writeSessionError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// Suggest godoc
// @Summary      Magic wand suggestion
// @Tags         sessions
// @Accept       json
// @Produce      json
// @Param        X-User-Id header string true "User ID"
// @Param        sessionId path string true "Session ID"
// @Param        body body usecase.SuggestInput true "Draft"
// @Success      200 {object} usecase.SuggestOutput
// @Router       /sessions/{sessionId}/suggest [post]
func (h *SessionHandler) Suggest(c *gin.Context) {
	var in usecase.SuggestInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	out, err := h.uc.Suggest(c.Request.Context(), middleware.GetUserID(c), c.Param("sessionId"), in)
	if err != nil {
		writeSessionError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// Finish godoc
// @Summary      Finish session and debrief
// @Tags         sessions
// @Accept       json
// @Produce      json
// @Param        X-User-Id header string true "User ID"
// @Param        sessionId path string true "Session ID"
// @Success      200 {object} usecase.FinishOutput
// @Router       /sessions/{sessionId}/finish [post]
func (h *SessionHandler) Finish(c *gin.Context) {
	out, err := h.uc.Finish(c.Request.Context(), middleware.GetUserID(c), c.Param("sessionId"))
	if err != nil {
		writeSessionError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// Get godoc
// @Summary      Get session debrief
// @Tags         sessions
// @Produce      json
// @Param        X-User-Id header string true "User ID"
// @Param        sessionId path string true "Session ID"
// @Success      200 {object} usecase.FinishOutput
// @Router       /sessions/{sessionId} [get]
func (h *SessionHandler) Get(c *gin.Context) {
	out, err := h.uc.GetDebrief(c.Request.Context(), middleware.GetUserID(c), c.Param("sessionId"))
	if err != nil {
		writeSessionError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// List godoc
// @Summary      List sessions history
// @Tags         sessions
// @Produce      json
// @Param        X-User-Id header string true "User ID"
// @Success      200 {object} usecase.SessionsListOutput
// @Router       /sessions [get]
func (h *SessionHandler) List(c *gin.Context) {
	out, err := h.uc.List(c.Request.Context(), middleware.GetUserID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, out)
}

func writeSessionError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, usecase.ErrSessionNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "session not found"})
	case errors.Is(err, usecase.ErrSessionFinished):
		c.JSON(http.StatusConflict, gin.H{"error": "session finished"})
	case errors.Is(err, usecase.ErrInvalidGoal), errors.Is(err, usecase.ErrInvalidPersona):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
