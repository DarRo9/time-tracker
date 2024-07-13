package controllers

import (
	"net/http"
	"strconv"

	db "github.com/DarRo9/time-tracker/internal/database"
	"github.com/DarRo9/time-tracker/internal/logger"
	"github.com/DarRo9/time-tracker/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type UserController struct {
	userRepo *db.UserRepository
}

func NewUserController(userRepo *db.UserRepository) *UserController {
	return &UserController{userRepo: userRepo}
}

func (uc *UserController) GetUsers(c *gin.Context) {
	limit := 1
	offset := 0

	if limitStr := c.Query("limit"); limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{
				"limit": limitStr,
				"error": err,
			}).Error("Неверное значение указанного лимита")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit value"})
			return
		}
		limit = parsedLimit
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		parsedOffset, err := strconv.Atoi(offsetStr)
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{
				"offset": offsetStr,
				"error":  err,
			}).Error("Указано неверное значение отступа ")
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset value"})
			return
		}
		offset = parsedOffset
	}

	filter := make(map[string]interface{})
	allowedFilters := []string{"passport_number", "surname", "name", "patronymic", "address"}

	for _, filterReq := range allowedFilters {
		if filterVal := c.Query(filterReq); filterVal != "" {
			filter[filterReq] = filterVal
		}
	}

	users, err := uc.userRepo.GetUsers(c, filter, limit, offset)
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"filter": filter,
			"limit":  limit,
			"offset": offset,
			"error":  err,
		}).Error("Не удалось получить информацию о пользователях")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (uc *UserController) AddUser(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"error": err,
		}).Error("Не удалось забиндить модель с данными")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный запрос"})
		return
	}

	if err := uc.userRepo.CreateUser(c, &user); err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"user":  user,
			"error": err,
		}).Error("Произошла ошибка при попытке создать пользователя")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "упс, не получилось создать пользователя"})
		return
	}

	logger.Logger.WithFields(logrus.Fields{
		"user": user,
	}).Info("Ура, пользователь был создан и добавлен!")

	c.JSON(http.StatusCreated, gin.H{"msg": "Пользователь добавлен!"})
}

func (uc *UserController) UpdateUser(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"error": err,
		}).Error("Не удалось забиндить модель с данными")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный запрос"})
		return
	}

	if err := uc.userRepo.UpdateUser(c, &user); err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"user":  user,
			"error": err,
		}).Error("Произошла ошибка при попытке обновить информацию об пользователе")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Logger.WithFields(logrus.Fields{
		"user": user,
	}).Info("Информация об пользователе была успешно обновлена")
	c.JSON(http.StatusOK, gin.H{"msg": "Информация пользователя была успешно изменена!"})
}

func (uc *UserController) DeleteUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"userID": c.Param("userID"),
			"error":  err,
		}).Error("Неверный юзер-айди")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверная айдишка пользователя"})
		return
	}
	if err := uc.userRepo.DeleteUser(c, userID); err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"userID": c.Param("userID"),
			"error":  err,
		}).Error("Не получилось удалить юзера")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	logger.Logger.WithFields(logrus.Fields{
		"userID": c.Param("userID"),
		"error":  err,
	}).Info("Юзер был удален")
	c.JSON(http.StatusOK, gin.H{"msg": "Пользователь был успешно удален"})
}
