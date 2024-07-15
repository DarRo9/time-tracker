package controllers

import (
	"net/http"
	"strconv"

	db "time-tracker/internal/database"
	"time-tracker/internal/logger"
	"time-tracker/internal/models"

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
			}).Error("Invalid value of the specified limit")
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
			}).Error("Invalid indent value specified")
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
		}).Error("Failed to get user information")
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
		}).Error("Failed to bind model to data")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to bind model to data"})
		return
	}

	if err := uc.userRepo.CreateUser(c, &user); err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"user":  user,
			"error": err,
		}).Error("An error occurred while trying to create a user")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred while trying to create a user"})
		return
	}

	logger.Logger.WithFields(logrus.Fields{
		"user": user,
	}).Info("the user has been created and added!")

	c.JSON(http.StatusCreated, gin.H{"msg": "the user has been created and added!"})
}

func (uc *UserController) UpdateUser(c *gin.Context) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"error": err,
		}).Error("Failed to bind model to data")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to bind model to data"})
		return
	}

	if err := uc.userRepo.UpdateUser(c, &user); err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"user":  user,
			"error": err,
		}).Error("An error occurred while trying to update user information")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Logger.WithFields(logrus.Fields{
		"user": user,
	}).Info("User information has been successfully updated")
	c.JSON(http.StatusOK, gin.H{"msg": "User information has been successfully updated"})
}

func (uc *UserController) DeleteUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"userID": c.Param("userID"),
			"error":  err,
		}).Error("Invalid user ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	if err := uc.userRepo.DeleteUser(c, userID); err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"userID": c.Param("userID"),
			"error":  err,
		}).Error("Failed to delete user")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	logger.Logger.WithFields(logrus.Fields{
		"userID": c.Param("userID"),
		"error":  err,
	}).Info("The user has been deleted")
	c.JSON(http.StatusOK, gin.H{"msg": "The user has been deleted"})
}
