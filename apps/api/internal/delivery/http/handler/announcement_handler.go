package handler

import (
	"announcement-api/internal/delivery/http/middleware"
	"announcement-api/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AnnouncementHandler struct {
	announcementUC *usecase.AnnouncementUsecase
}

func NewAnnouncementHandler(annUC *usecase.AnnouncementUsecase) *AnnouncementHandler {
	return &AnnouncementHandler{announcementUC: annUC}
}

type AnnouncementRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

func (h *AnnouncementHandler) Create(c *gin.Context) {
	var req AnnouncementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	requesterID := middleware.CurrentUserID(c)

	ann, err := h.announcementUC.Create(c.Request.Context(), requesterID, req.Title, req.Content)
	if err != nil {
		handleAnnouncementError(c, err)
		return
	}

	c.JSON(http.StatusCreated, ann)
}

func (h *AnnouncementHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req AnnouncementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	requesterID := middleware.CurrentUserID(c)

	ann, err := h.announcementUC.Update(c.Request.Context(), requesterID, id, req.Title, req.Content)
	if err != nil {
		handleAnnouncementError(c, err)
		return
	}

	c.JSON(http.StatusOK, ann)
}

func (h *AnnouncementHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	requesterID := middleware.CurrentUserID(c)

	if err := h.announcementUC.Delete(c.Request.Context(), requesterID, id); err != nil {
		handleAnnouncementError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *AnnouncementHandler) List(c *gin.Context) {
	anns, err := h.announcementUC.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch announcements"})
		return
	}
	c.JSON(http.StatusOK, anns)
}

func (h *AnnouncementHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	ann, err := h.announcementUC.Get(c.Request.Context(), id)
	if err != nil {
		handleAnnouncementError(c, err)
		return
	}

	c.JSON(http.StatusOK, ann)
}

func handleAnnouncementError(c *gin.Context, err error) {
	switch err {
	case usecase.ErrForbidden:
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
	case usecase.ErrNotFound:
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
