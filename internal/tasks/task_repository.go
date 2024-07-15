package tasks

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"time-tracker/internal/apperrors"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) createTask(ctx context.Context, taskInfo *task) error {
	slog.Debug("Creating task", "user", taskInfo)
	tx := r.db.WithContext(ctx)
	if err := tx.Create(&taskInfo).Error; err != nil {
		slog.Error("Error creating user", "error", err)
		var perr *pgconn.PgError
		if errors.As(err, &perr) {
			if perr.Code == "23505" {
				return &apperrors.DuplicateKeyError{Message: fmt.Sprintf("Task with name %v already exists for this user", taskInfo.TaskName)}
			} else if perr.Code == "23503" {
				return &apperrors.NoUserError{Message: fmt.Sprintf("User with id %v doesn't exist", taskInfo.UserId)}
			}
		}
		return err
	}
	slog.Info("Task succsesfully created", "task", taskInfo)
	return nil
}

func (r *repository) GetTasksByUser(ctx context.Context, taskReq *getTasksReq) (*[]task, error) {
	slog.Debug("Getting tasks", "task period to get", taskReq)
	var tasks []task
	tx := r.db.WithContext(ctx).
		Where("user_id = ? AND start_time >= ? AND (end_time <= ? OR end_time IS NULL)", taskReq.UserId, taskReq.StartDate, taskReq.EndDate).
		Find(&tasks)
	if tx.Error != nil {
		slog.Error("Error getting user tasks", "task request", taskReq)
		return nil, tx.Error
	}

	slog.Info("Successfully get tasks", "tasks", tasks)
	return &tasks, nil
}

func (r *repository) EndTask(ctx context.Context, taskId int, endTime time.Time) error {
	slog.Debug("Ending task", "task id", taskId, "end time", endTime)
	var task task
	if err := r.db.WithContext(ctx).First(&task, taskId).Error; err != nil {
		return &apperrors.NoTaskError{Message: fmt.Sprintf("No task with id %v", taskId)}
	}

	if task.EndTime != nil {
		return &apperrors.TaskAlreadyEndedError{Message: "Task already ended"}
	}

	tx := r.db.WithContext(ctx)
	if err := tx.Model(&task).Where("id = ?", taskId).Update("end_time", endTime).Error; err != nil {
		slog.Error("Error ending task", "error", err)
		return err
	}

	slog.Info("Task succsesfully ended", "task id", taskId)
	return nil
}
