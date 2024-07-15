package users

import (
	"log/slog"
	"net/http"
	"strconv"

	"time-tracker/internal/apperrors"
	"time-tracker/internal/responses"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service
}

func NewHandler(s Service) *Handler {
	return &Handler{Service: s}
}

func (h *Handler) CreateUser(c *gin.Context) {
	var passportInfo passport
	if err := c.ShouldBindJSON(&passportInfo); err != nil {
		slog.Error("Error binding JSON", "error", err)
		c.JSON(http.StatusBadRequest, responses.ErrResponse{Error: err.Error()})
		return
	}
	userInfo := user{Passport: passportInfo.PassportNumber}
	err := h.Service.createUser(c.Request.Context(), &userInfo)
	if err != nil {
		switch err.(type) {
		case *apperrors.BadRequestError:
			c.JSON(http.StatusBadRequest, responses.ErrResponse{Error: err.Error()})
		case *apperrors.ExternalAPIError:
			c.JSON(http.StatusInternalServerError, responses.ErrResponse{Error: err.Error()})
		case *apperrors.DuplicateKeyError:
			c.JSON(http.StatusConflict, responses.ErrResponse{Error: err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, responses.ErrResponse{Error: "Internal server error"})
		}
		return
	}
	c.JSON(http.StatusCreated, responses.SuccessResponse{Message: "User created"})
}

func (h *Handler) GetUsers(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page <= 0 {
		slog.Error("Invalig page", "error", err)
		c.JSON(http.StatusBadRequest, responses.ErrResponse{Error: "Invalid page number"})
		return
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	if err != nil || pageSize <= 0 {
		slog.Error("Invalid pageSize", "error", err)
		c.JSON(http.StatusBadRequest, responses.ErrResponse{Error: "Invalid page size"})
		return
	}

	filters := make(map[string]string)
	id := c.Query("id")
	if id != "" {
		filters["id"] = id
	}
	passportNumber := c.Query("passportNumber")
	if passportNumber != "" {
		filters["passportNumber"] = passportNumber
	}
	surname := c.Query("surname")
	if surname != "" {
		filters["surname"] = surname
	}
	name := c.Query("name")
	if name != "" {
		filters["name"] = name
	}
	patronymic := c.Query("patronymic")
	if patronymic != "" {
		filters["patronymic"] = patronymic
	}
	address := c.Query("address")
	if address != "" {
		filters["address"] = address
	}
	searchInfo := &searchReq{page: page, pageSize: pageSize, filters: filters}

	users, err := h.Service.getUsers(c, searchInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrResponse{Error: "Internal server error"})
		return
	}
	if len(*users) == 0 {
		c.JSON(http.StatusNoContent, responses.SuccessResponse{})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h *Handler) RemoveUser(c *gin.Context) {
	var userId userId
	if err := c.ShouldBindJSON(&userId); err != nil {
		slog.Error("Error binding JSON", "error", err)
		c.JSON(http.StatusBadRequest, responses.ErrResponse{Error: err.Error()})
		return
	}
	err := h.Service.removeUser(c.Request.Context(), &userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, responses.SuccessResponse{Message: "User deleted"})
}

func (h *Handler) UpdateUser(c *gin.Context) {
	var userInfo updateUserReq
	err := c.ShouldBindJSON(&userInfo)
	if err != nil {
		slog.Error("Error binding JSON", "error", err)
		c.JSON(http.StatusBadRequest, responses.ErrResponse{Error: err.Error()})
		return
	}

	idStr := c.Query("id")
	var id userId
	id.Id, err = strconv.Atoi(idStr)
	if err != nil || id.Id <= 0 {
		slog.Error("Invalid id", "error", err)
		c.JSON(http.StatusBadRequest, responses.ErrResponse{Error: "Invalid id"})
		return
	}

	if userInfo.Name == nil && userInfo.Passport == nil && userInfo.Surname == nil && userInfo.Patronymic == nil && userInfo.Address == nil {
		slog.Error("No info to update", "user info", userInfo)
		c.JSON(http.StatusBadRequest, responses.ErrResponse{Error: "No info to update"})
		return
	}

	err = h.Service.updateUser(c.Request.Context(), &id, &userInfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse{Message: "User updated"})
}
