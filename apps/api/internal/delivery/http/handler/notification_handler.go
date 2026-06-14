package handler

import (
	"announcement-api/internal/delivery/http/middleware"
	"announcement-api/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type NotificationHandler struct {
	notifUC *usecase.NotificationUsecase
}

func NewNotificationHandler(notifUC *usecase.NotificationUsecase) *NotificationHandler {
	return &NotificationHandler{notifUC: notifUC}
}

func (h *NotificationHandler) List(c *gin.Context) {
	notifications, err := h.notifUC.List(c.Request.Context(), middleware.CurrentUserID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch notifications"})
		return
	}
	c.JSON(http.StatusOK, notifications)
}

func (h *NotificationHandler) UnreadCount(c *gin.Context) {
	userID := middleware.CurrentUserID(c)

	count, err := h.notifUC.UnreadCount(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count notifications"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"unread_count": count})
}

func (h *NotificationHandler) MarkRead(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	userID := middleware.CurrentUserID(c)

	if err := h.notifUC.MarkRead(c.Request.Context(), id, userID); err != nil {
		if err == usecase.ErrNotificationNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to mark as read"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "marked as read"})
}

func (h *NotificationHandler) MarkAllRead(c *gin.Context) {
	userID := middleware.CurrentUserID(c)

	if err := h.notifUC.MarkAllRead(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to mark all as read"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "all marked as read"})
}
