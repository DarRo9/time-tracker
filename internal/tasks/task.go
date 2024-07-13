package tasks

import (
	"context"
	"time"
)

type task struct {
	Id        int        `gorm:"column:id"`
	TaskName  string     `json:"taskName"`
	UserId    int        `json:"userId"`
	StartTime time.Time  `json:"startTime"`
	EndTime   *time.Time `json:"endTime"`
}

type createTaskReq struct {
	UserId   int    `json:"userId"`
	TaskName string `json:"taskName"`
}

type getTasksReq struct {
	UserId    int        `json:"userId"`
	StartDate *time.Time `json:"startDate"`
	EndDate   *time.Time `json:"endDate"`
}

type getTasksRes struct {
	Id        int           `json:"id"`
	TaskName  string        `json:"taskName"`
	TaskTime  time.Duration `json:"taskTime"`
	Completed bool          `json:"completed"`
}

type Repository interface {
	createTask(ctx context.Context, taskInfo *task) error
	GetTasksByUser(ctx context.Context, taskReq *getTasksReq) (*[]task, error)
	EndTask(ctx context.Context, taskId int, endTime time.Time) error
}

type Service interface {
	createTask(ctx context.Context, task *createTaskReq) error
	GetTasksByUser(ctx context.Context, taskReq *getTasksReq) (*[]getTasksRes, error)
	EndTask(ctx context.Context, taskId int) error
}
