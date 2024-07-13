package database

import (
	"context"
	"time"

	"github.com/DarRo9/time-tracker/internal/logger"
	"github.com/DarRo9/time-tracker/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type TaskRepository struct {
	db *pgxpool.Pool
}

func NewTaskRepository(db *pgxpool.Pool) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) GetUserTasksByPeriod(ctx context.Context, userID int, start, end time.Time) ([]models.Task, error) {
	logger.Logger.WithFields(logrus.Fields{
		"userID": userID,
		"start":  start,
		"end":    end,
	}).Debug("Запрос тасок пользователя за период")

	var tasks []models.Task
	query := `
	SELECT id, user_id, description, start_time, end_time, created_at, updated_at
	FROM tasks
	WHERE user_id = $1 AND start_time >= $2 AND end_time <= $3
	ORDER BY EXTRACT(EPOCH FROM (end_time - start_time)) DESC
	`
	rows, err := r.db.Query(ctx, query, userID, start, end)
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"error": err,
		}).Error("Произошла ошибка при выполнении запроса получения тасок")
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.UserID, &task.Description, &task.StartTime, &task.EndTime, &task.CreatedAt, &task.UpdatedAt); err != nil {
			logger.Logger.WithFields(logrus.Fields{
				"error": err,
			}).Error("Произошла ошибка при сканировании строки тасок")
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if rows.Err() != nil {
		logger.Logger.WithFields(logrus.Fields{
			"error": err,
		}).Error("Произошла ошибка при итерации по строкам данных")
		return nil, rows.Err()
	}

	logger.Logger.WithFields(logrus.Fields{
		"userID": userID,
		"start":  start,
		"end":    end,
		"count":  len(tasks),
	}).Info("Успешно получены таски юзера за период")

	return tasks, nil
}

func (r *TaskRepository) StartTask(ctx context.Context, userID int, description string) error {
	logger.Logger.WithFields(logrus.Fields{
		"userID":      userID,
		"description": description,
	}).Debug("Начало новой таски")

	query := `
			INSERT INTO tasks (user_id, description, start_time, created_at, updated_at)
			VALUES ($1, $2, NOW(), NOW(), NOW())
		`
	_, err := r.db.Exec(ctx, query, userID, description)
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"userID":      userID,
			"description": description,
			"error":       err,
		}).Error("Произошла ошибка при попытка начать новую таску")
		return err
	}

	logger.Logger.WithFields(logrus.Fields{
		"userID":      userID,
		"description": description,
	}).Info("Таска успешно начата")

	return nil
}

func (r *TaskRepository) EndTask(ctx context.Context, taskID int) error {
	logger.Logger.WithFields(logrus.Fields{
		"taskID": taskID,
	}).Debug("Окончание таски")

	query := `
			UPDATE tasks
			SET end_time = NOW(), updated_at = NOW()
			WHERE id = $1 AND end_time IS NULL
		`
	_, err := r.db.Exec(ctx, query, taskID)
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"taskID": taskID,
			"error":  err,
		}).Error("Произошла ошибка при окончании таски")
		return err
	}

	logger.Logger.WithFields(logrus.Fields{
		"taskID": taskID,
	}).Info("Таска успешно завершена")

	return nil
}
