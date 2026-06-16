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

// @Summary     Create new announcement
// @Description Create new announcement and broadcasted via SSE. Only admin.
// @Tags        Announcements
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       body  body      dto.AnnouncementRequest  true  "Announcement data"
// @Success     201   {object}  entity.Announcement
// @Failure     400   {object}  dto.ErrorResponse  "Validation failed"
// @Failure     401   {object}  dto.ErrorResponse  "Token not valid"
// @Failure     403   {object}  dto.ErrorResponse  "Not admin"
// @Failure     500   {object}  dto.ErrorResponse  "Internal server error"
// @Router      /announcements [post]
func (h *AnnouncementHandler) Create(c *gin.Context) {
	var req dto.AnnouncementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
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

// @Summary     Update announcement
// @Description Update announcement's title and content. Update notification broadcasted via SSE. Only admin.
// @Tags        Announcements
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       id    path      string                   true  "Announcement UUID"
// @Param       body  body      dto.AnnouncementRequest  true  "Updated announcement data"
// @Success     200   {object}  entity.Announcement
// @Failure     400   {object}  dto.ErrorResponse  "UUID or body not valid"
// @Failure     401   {object}  dto.ErrorResponse  "Token not valid"
// @Failure     403   {object}  dto.ErrorResponse  "Not admin"
// @Failure     404   {object}  dto.ErrorResponse  "Announcement not found"
// @Failure     500   {object}  dto.ErrorResponse  "Internal server error"
// @Router      /announcements/{id} [put]
func (h *AnnouncementHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}

	var req AnnouncementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
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

// @Summary     Delete announcement
// @Description Delete announcement and all related notification. Event delete broadcasted via SSE. Only admin.
// @Tags        Announcements
// @Produce     json
// @Security    BearerAuth
// @Param       id  path      string  true  "Announcement UUID"
// @Success     200  {object}  dto.MessageResponse
// @Failure     400  {object}  dto.ErrorResponse  "UUID not valid"
// @Failure     401  {object}  dto.ErrorResponse  "Token not valid"
// @Failure     403  {object}  dto.ErrorResponse  "Not admin"
// @Failure     500  {object}  dto.ErrorResponse  "Internal server error"
// @Router      /announcements/{id} [delete]
func (h *AnnouncementHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
		return
	}

	requesterID := middleware.CurrentUserID(c)

	if err := h.announcementUC.Delete(c.Request.Context(), requesterID, id); err != nil {
		handleAnnouncementError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.MessageResponse{Message: "announcement deleted"})
}

// @Summary     List of announcements
// @Description Return all announcements ordered by recent.
// @Tags        Announcements
// @Produce     json
// @Security    BearerAuth
// @Success     200  {array}   entity.Announcement
// @Failure     401  {object}  dto.ErrorResponse  "Token not valid or no token detected"
// @Failure     500  {object}  dto.ErrorResponse  "Internal server error"
// @Router      /announcements [get]
func (h *AnnouncementHandler) List(c *gin.Context) {
	anns, err := h.announcementUC.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "failed to fetch announcements"})
		return
	}
	c.JSON(http.StatusOK, anns)
}

// @Summary     Detail announcement
// @Description Return one announcement based on ID
// @Tags        Announcements
// @Produce     json
// @Security    BearerAuth
// @Param       id   path      string  true  "Announcement UUID"  example("550e8400-e29b-41d4-a716-446655440000")
// @Success     200  {object}  entity.Announcement
// @Failure     400  {object}  dto.ErrorResponse  "UUID not valid"
// @Failure     401  {object}  dto.ErrorResponse  "Token not valid"
// @Failure     404  {object}  dto.ErrorResponse  "Announcement not found"
// @Router      /announcements/{id} [get]
func (h *AnnouncementHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid id"})
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
		c.JSON(http.StatusForbidden, dto.ErrorResponse{Error: err.Error()})
	case usecase.ErrNotFound:
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal server error"})
	}
}

var _ entity.Announcement
