package tasks

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/DarRo9/time-tracker/internal/apperrors"
	"github.com/DarRo9/time-tracker/internal/responses"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service
}

func NewHandler(s Service) *Handler {
	return &Handler{Service: s}
}

func (h *Handler) CreateTask(c *gin.Context) {
	var taskReq createTaskReq
	if err := c.ShouldBindJSON(&taskReq); err != nil {
		slog.Error("Error binding JSON", "error", err)
		c.JSON(http.StatusBadRequest, responses.ErrResponse{Error: err.Error()})
		return
	}
	err := h.Service.createTask(c.Request.Context(), &taskReq)
	if err != nil {
		switch err.(type) {
		case *apperrors.NoUserError:
			c.JSON(http.StatusBadRequest, responses.ErrResponse{Error: err.Error()})
		case *apperrors.DuplicateKeyError:
			c.JSON(http.StatusConflict, responses.ErrResponse{Error: err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, responses.ErrResponse{Error: err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, responses.SuccessResponse{Message: "Task created"})
}

func (h *Handler) GetTasksByUser(c *gin.Context) {
	userId, err := strconv.Atoi(c.Query("userId"))
	if err != nil {
		slog.Error("Invalid user Id", "error", err)
		c.JSON(http.StatusBadRequest, responses.ErrResponse{Error: "Invalid user Id"})
		return
	}

	var startDate, endDate *time.Time

	startDateStr := c.Query("startTime")
	if startDateStr != "" {
		startDateParsed, err := time.Parse(time.RFC3339, startDateStr)
		if err != nil {
			slog.Error("Invalid start date format", "error", err)
			c.JSON(http.StatusBadRequest, responses.ErrResponse{Error: "Invalid start date format"})
			return
		}
		startDate = &startDateParsed
	}

	endDateStr := c.Query("endTime")
	if endDateStr != "" {
		endDateParsed, err := time.Parse(time.RFC3339, endDateStr)
		if err != nil {
			slog.Error("Invalid end date format", "error", err)
			c.JSON(http.StatusBadRequest, responses.ErrResponse{Error: "Invalid end date format"})
			return
		}
		endDate = &endDateParsed
	}

	taskReq := &getTasksReq{
		UserId:    userId,
		StartDate: startDate,
		EndDate:   endDate,
	}

	tasks, err := h.Service.GetTasksByUser(c.Request.Context(), taskReq)
	if err != nil {
		switch err.(type) {
		case *apperrors.NoTaskError:
			c.JSON(http.StatusNoContent, responses.SuccessResponse{})
		default:
			c.JSON(http.StatusInternalServerError, responses.ErrResponse{Error: "Failed to get user tasks"})
		}
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func (h *Handler) EndTask(c *gin.Context) {
	taskId, err := strconv.Atoi(c.Query("taskId"))
	if err != nil {
		slog.Error("Invalid task Id", "error", err)
		c.JSON(http.StatusBadRequest, responses.ErrResponse{Error: "Invalid task Id"})
		return
	}

	if err := h.Service.EndTask(c.Request.Context(), taskId); err != nil {
		switch err.(type) {
		case *apperrors.TaskAlreadyEndedError:
			c.JSON(http.StatusConflict, responses.ErrResponse{Error: err.Error()})
		case *apperrors.NoTaskError:
			c.JSON(http.StatusNoContent, responses.SuccessResponse{})
		default:
			c.JSON(http.StatusInternalServerError, responses.ErrResponse{Error: "Failed to update task end time"})
		}
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse{Message: "Task ended"})
}
