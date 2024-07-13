package models

import "time"

type Task struct {
	ID          int       `json:"id"`
	UserID      int       `json:"userId"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"startTime"`
	EndTime     time.Time `json:"endTime,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type Request struct {
	UserID      uint   `json:"user_id"`
	Description string `json:"description"`
}
