package controllers

import (
	"net/http"
	"strconv"
	"time"

	db "time-tracker/internal/database"
	"time-tracker/internal/logger"
	"time-tracker/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type TaskController struct {
	taskRepo *db.TaskRepository
}

func NewTaskController(taskRepo *db.TaskRepository) *TaskController {
	return &TaskController{taskRepo: taskRepo}
}

// CreateTags godoc
// @Summary     Get user tasks by period
// @Description Get tasks for a user within a specified time period
// @Tags        tasks
// @Accept      json
// @Produce     json
// @Param       userID path     int    true "User ID"
// @Param       start  query    string true "Start time in RFC3339 format"
// @Param       end    query    string true "End time in RFC3339 format"
// @Success     200    {array}  models.Task
// @Failure     400    {object} gin.H
// @Failure     500    {object} gin.H
// @Router      /users/{userID}/tasks [get]
func (tc *TaskController) GetUserTasksByPeriod(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"userID": c.Param("userID"),
			"error":  err,
		}).Error("User ID is incorrect")
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is incorrect"})
		return
	}

	start, err := time.Parse(time.RFC3339, c.Query("start"))
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"start": c.Query("start"),
			"error": err,
		}).Error("Invalid start time")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start time"})
		return
	}

	end, err := time.Parse(time.RFC3339, c.Query("end"))
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"end":   c.Query("end"),
			"error": err,
		}).Error("Invalid end time")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end time"})
		return
	}

	tasks, err := tc.taskRepo.GetUserTasksByPeriod(c, userID, start, end)
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"userID": c.Param("userID"),
			"start":  start,
			"end":    end,
			"error":  err,
		}).Error("Failed to get user tasks for the period")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	logger.Logger.WithFields(logrus.Fields{
		"userID": c.Param("userID"),
		"start":  start,
		"end":    end,
	}).Info("Successfully receiving user tasks for the period")

	c.JSON(http.StatusOK, tasks)
}

// @Summary     Start a new task
// @Description Start tracking time for a new task
// @Tags        tasks
// @Accept      json
// @Produce     json
// @Param       task body     models.Request true "Task to start"
// @Success     201  {object} models.Task
// @Failure     400  {object} gin.H
// @Failure     500  {object} gin.H
// @Router      /tasks/start [post]
func (tc *TaskController) StartTask(c *gin.Context) {
	var req models.Request

	if err := c.BindJSON(&req); err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"error": err,
		}).Error("Failed to bind model to data")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to bind model to data"})
		return
	}

	if err := tc.taskRepo.StartTask(c, int(req.UserID), req.Description); err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"userID":      int(req.UserID),
			"description": req.Description,
			"error":       err,
		}).Error("Failed to start task")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Logger.WithFields(logrus.Fields{
		"userID":      int(req.UserID),
		"description": req.Description,
	}).Info("The task has been started")
	c.JSON(http.StatusCreated, gin.H{"msg": "The task has been started"})
}

// @Summary     End a task
// @Description End tracking time for a task
// @Tags        tasks
// @Accept      json
// @Produce     json
// @Param       taskID path     int true "Task ID"
// @Success     200    {object} models.Task
// @Failure     400    {object} gin.H
// @Failure     500    {object} gin.H
// @Router      /tasks/end/{taskID} [post]
func (tc *TaskController) EndTask(c *gin.Context) {
	taskID, err := strconv.Atoi(c.Param("taskID"))
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"taskID": taskID,
			"error":  err,
		}).Error("Invalid task ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	if err := tc.taskRepo.EndTask(c, taskID); err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"taskID": taskID,
			"error":  err,
		}).Error("Failed to finish task")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	logger.Logger.WithFields(logrus.Fields{
		"taskID": taskID,
	}).Info("The task was over")

	c.JSON(http.StatusOK, gin.H{"msg": "The task was over"})
}
