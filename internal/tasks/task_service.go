package tasks

import (
	"context"
	"log/slog"
	"sort"
	"time"

	"time-tracker/internal/apperrors"
)

type service struct {
	Repository
	timeout time.Duration
}

func NewService(r Repository) Service {
	return &service{
		Repository: r,
		timeout:    time.Duration(2) * time.Second}
}

func (s *service) createTask(c context.Context, taskReq *createTaskReq) error {
	slog.Debug("Creating new task in service", "task", taskReq)
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer func() {
		cancel()
	}()

	taskInfo := &task{
		UserId:    taskReq.UserId,
		TaskName:  taskReq.TaskName,
		StartTime: time.Now().UTC(),
	}

	if err := s.Repository.createTask(ctx, taskInfo); err != nil {
		return err
	}

	return nil
}

func (s *service) GetTasksByUser(ctx context.Context, taskReq *getTasksReq) (*[]getTasksRes, error) {
	slog.Debug("Getting tasks in service", "task", taskReq)
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	if taskReq.StartDate == nil {
		startDate := time.Date(0001, 1, 1, 0, 0, 0, 0, time.UTC)
		taskReq.StartDate = &startDate
	}
	if taskReq.EndDate == nil {
		endDate := time.Now().UTC()
		taskReq.EndDate = &endDate
	}

	tasks, err := s.Repository.GetTasksByUser(ctx, taskReq)
	if err != nil {
		return nil, err
	}

	if len(*tasks) == 0 {
		return nil, &apperrors.NoTaskError{Message: "No tasks with search"}
	}

	taskResp := make([]getTasksRes, 0, len(*tasks))
	for _, task := range *tasks {
		var taskTime time.Duration
		var completed bool

		if task.EndTime != nil {
			taskTime = time.Duration(task.EndTime.Sub(task.StartTime).Minutes())
			completed = true
		} else {
			taskTime = time.Duration(time.Since(task.StartTime).Minutes())
		}

		taskResp = append(taskResp, getTasksRes{
			Id:        task.Id,
			TaskName:  task.TaskName,
			TaskTime:  taskTime,
			Completed: completed,
		})

		sort.Slice(taskResp, func(i, j int) bool {
			return taskResp[i].TaskTime > taskResp[j].TaskTime
		})
	}

	return &taskResp, nil
}

func (s *service) EndTask(ctx context.Context, taskId int) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	endTime := time.Now().UTC()
	slog.Debug("Updating task end time in service", "taskId", taskId, "endTime", endTime)

	if err := s.Repository.EndTask(ctx, taskId, endTime); err != nil {
		return err
	}

	return nil
}
