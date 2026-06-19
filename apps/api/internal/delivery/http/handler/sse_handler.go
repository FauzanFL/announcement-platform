package handler

import (
	"announcement-api/internal/delivery/http/dto"
	"announcement-api/internal/delivery/http/middleware"
	"announcement-api/internal/domain/entity"
	"announcement-api/internal/usecase"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SSEHandler struct {
	annUC   *usecase.AnnouncementUsecase
	notifUC *usecase.NotificationUsecase
}

func NewSSEHandler(annUC *usecase.AnnouncementUsecase, notifUC *usecase.NotificationUsecase) *SSEHandler {
	return &SSEHandler{annUC: annUC, notifUC: notifUC}
}

// @Summary     SSE stream real-time
// @Description Open Server-Sent Events (SSE) connection to sent real-time event. token sent via query parameter. There are two type of events: `announcement` (new/update/delete) and `unread_count` (number of unread notification). Heartbeat sent 15 second as SSE comment.
// @Tags        SSE
// @Produce     text/event-stream
// @Security    BearerAuth
// @Param       token  query     string  true  "JWT token"  example("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...")
// @Success     200    {string}  string  "Stream of SSE events"
// @Failure     401    {object}  dto.ErrorResponse  "Token not valid"
// @Failure     500    {object}  dto.ErrorResponse  "Streaming not supported"
// @Router      /stream [get]
func (h *SSEHandler) Stream(c *gin.Context) {
	userID := middleware.CurrentUserID(c)

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "streaming unsupported"})
		return
	}

	ctx := c.Request.Context()

	msgCh, cleanup, err := h.annUC.SubscribeToEvents(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to subscribe"})
		return
	}
	defer cleanup()

	h.sendUnreadCount(c, ctx, userID, flusher)

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case payload, open := <-msgCh:
			if !open {
				return
			}

			var event entity.AnnouncementEvent
			if err := json.Unmarshal([]byte(payload), &event); err != nil {
				continue
			}

			data, _ := json.Marshal(event)
			fmt.Fprintf(c.Writer, "event: announcement\ndata: %s\n\n", data)
			flusher.Flush()

			h.sendUnreadCount(c, ctx, userID, flusher)

		case <-ticker.C:
			fmt.Fprintf(c.Writer, ": heartbeat\n\n")
			flusher.Flush()
		}
	}
}

func (h *SSEHandler) sendUnreadCount(c *gin.Context, ctx context.Context, userID uuid.UUID, flusher http.Flusher) {
	count, err := h.notifUC.UnreadCount(ctx, userID)
	if err != nil {
		return
	}
	payload, _ := json.Marshal(map[string]interface{}{"unread_count": count})
	fmt.Fprintf(c.Writer, "event: unread_count\ndata: %s\n\n", payload)
	flusher.Flush()
}
