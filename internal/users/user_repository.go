package users

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/DarRo9/time-tracker/internal/apperrors"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) createUser(ctx context.Context, user *user) error {
	slog.Debug("Creating user", "user", user)
	tx := r.db.WithContext(ctx)
	if err := tx.Create(&user).Error; err != nil {
		slog.Error("Error creating user", "error", err)
		var perr *pgconn.PgError
		if errors.As(err, &perr) {
			if perr.Code == "23505" {
				return &apperrors.DuplicateKeyError{Message: "User already exists"}
			}
		}
		return err
	}
	slog.Info("User succsesfully created", "user", user)
	return nil
}

func (r *repository) getUsers(ctx context.Context, s *searchReq) (*[]user, error) {
	slog.Debug("Getting users", "search info", s)
	var users []user
	tx := r.db.WithContext(ctx).Limit(s.pageSize).Offset(s.offset)

	for key, value := range s.filters {
		switch key {
		case "id":
			tx = tx.Where("id = ?", value)
		case "passportNumber":
			tx = tx.Where("passport_number LIKE ?", "%"+value+"%")
		case "surname":
			tx = tx.Where("surname LIKE ?", "%"+value+"%")
		case "name":
			tx = tx.Where("name LIKE ?", "%"+value+"%")
		case "patronymic":
			tx = tx.Where("patronymic LIKE ?", "%"+value+"%")
		case "address":
			tx = tx.Where("address LIKE ?", "%"+value+"%")
		}
	}

	if err := tx.Find(&users).Error; err != nil {
		slog.Error("Error getting users", "error", err)
		return nil, err
	}
	slog.Info("Successfully get users", "users", users)
	return &users, nil
}

func (r *repository) removeUser(ctx context.Context, id *userId) error {
	slog.Debug("Deleting user", "user", id)
	tx := r.db.WithContext(ctx)
	res := tx.Delete(&user{}, id.Id)
	if res.Error != nil {
		slog.Error("Error deleting user", "error", res.Error)
		return res.Error
	} else if res.RowsAffected == 0 {
		return &apperrors.NoRowsAffectedError{Message: fmt.Sprintf("No users with id%v", id.Id)}
	}
	slog.Info("Successfully deleted user")
	return nil
}

func (r *repository) updateUser(ctx context.Context, id *userId, updates *map[string]interface{}) error {
	slog.Debug("Updating user", "user", updates)
	tx := r.db.WithContext(ctx).Model(&user{}).Where("id = ?", id.Id)
	_ = tx
	if len(*updates) > 0 {
		if err := tx.Updates(updates).Error; err != nil {
			slog.Error("Error updating user", "error", err)
			return err
		}
	}

	slog.Info("Successfully updated user")
	return nil
}
