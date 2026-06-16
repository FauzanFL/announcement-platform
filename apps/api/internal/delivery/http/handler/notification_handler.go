package handler

import (
	"announcement-api/internal/delivery/http/dto"
	"announcement-api/internal/delivery/http/middleware"
	"announcement-api/internal/domain/entity"
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

// @Summary     List of user notifications
// @Description Return all notifications
// @Tags        Notifications
// @Produce     json
// @Security    BearerAuth
// @Success     200  {array}   entity.Notification
// @Failure     401  {object}  dto.ErrorResponse  "Token not valid"
// @Failure     500  {object}  dto.ErrorResponse  "Internal server error"
// @Router      /notifications [get]
func (h *NotificationHandler) List(c *gin.Context) {
	notifications, err := h.notifUC.List(c.Request.Context(), middleware.CurrentUserID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to fetch notifications"})
		return
	}
	c.JSON(http.StatusOK, notifications)
}

// @Summary     Unread all notifications count
// @Description Return total of unread notifications for login user
// @Tags        Notifications
// @Produce     json
// @Security    BearerAuth
// @Success     200  {object}  dto.UnreadCountResponse
// @Failure     401  {object}  dto.ErrorResponse  "Token not valid"
// @Failure     500  {object}  dto.ErrorResponse  "Internal server error"
// @Router      /notifications/unread-count [get]
func (h *NotificationHandler) UnreadCount(c *gin.Context) {
	userID := middleware.CurrentUserID(c)

	count, err := h.notifUC.UnreadCount(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to count notifications"})
		return
	}

	c.JSON(http.StatusOK, dto.UnreadCountResponse{UnreadCount: count})
}

// @Summary     Mark as read notification
// @Description Mark one notification as read based on ID. User can only mark their own notifications
// @Tags        Notifications
// @Produce     json
// @Security    BearerAuth
// @Param       id  path      string  true  "Notification UUID"  example("550e8400-e29b-41d4-a716-446655440000")
// @Success     200  {object}  dto.MessageResponse
// @Failure     400  {object}  dto.ErrorResponse  "UUID not valid"
// @Failure     401  {object}  dto.ErrorResponse  "Token not valid"
// @Failure     404  {object}  dto.ErrorResponse  "Notification not found"
// @Failure     500  {object}  dto.ErrorResponse  "Internal server error"
// @Router      /notifications/{id}/read [put]
func (h *NotificationHandler) MarkRead(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}

	userID := middleware.CurrentUserID(c)

	if err := h.notifUC.MarkRead(c.Request.Context(), id, userID); err != nil {
		if err == usecase.ErrNotificationNotFound {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to mark as read"})
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{Message: "marked as read"})
}

// @Summary     Mark all notifications as read
// @Description Mark all unread notifications for login user as read
// @Tags        Notifications
// @Produce     json
// @Security    BearerAuth
// @Success     200  {object}  dto.MessageResponse
// @Failure     401  {object}  dto.ErrorResponse  "Token not valid"
// @Failure     500  {object}  dto.ErrorResponse  "Internal server error"
// @Router      /notifications/read-all [put]

func (h *NotificationHandler) MarkAllRead(c *gin.Context) {
	userID := middleware.CurrentUserID(c)

	if err := h.notifUC.MarkAllRead(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to mark all as read"})
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{Message: "all marked as read"})
}

var _ entity.Notification
