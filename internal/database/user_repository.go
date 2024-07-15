package database

import (
	"fmt"

	"context"
	"time"

	"time-tracker/internal/logger"
	"time-tracker/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	logger.Logger.WithFields(logrus.Fields{
		"passport_number": user.PassportNumber,
		"surname":         user.Surname,
		"name":            user.Name,
		"patronymic":      user.Patronymic,
		"address":         user.Address,
	}).Debug("Создание юзера")

	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt
	query := `INSERT INTO users (passport_number, surname, name, patronymic, address, created_at, updated_at) 
              VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	err := r.db.QueryRow(ctx, query, user.PassportNumber, user.Surname, user.Name, user.Patronymic, user.Address, user.CreatedAt, user.UpdatedAt).Scan(&user.ID)
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"error": err,
		}).Error("An error occurred while creating a user")
	}

	logger.Logger.WithFields(logrus.Fields{
		"userID": user.ID,
	}).Info("The user has been created")

	return nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	logger.Logger.WithFields(logrus.Fields{
		"userID": id,
	}).Debug("Getting a user by ID")

	user := &models.User{}
	query := `SELECT id, passport_number, surname, name, patronymic, address, created_at, updated_at FROM users WHERE id=$1`
	err := r.db.QueryRow(ctx, query, id).Scan(&user.ID, &user.PassportNumber, &user.Surname, &user.Name, &user.Patronymic, &user.Address, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"userID": id,
			"error":  err,
		}).Error("An error occurred while retrieving the user's ID")
		return nil, err
	}

	logger.Logger.WithFields(logrus.Fields{
		"userID": id,
	}).Info("User data successfully received")

	return user, nil

}

func (r *UserRepository) UpdateUser(ctx context.Context, user *models.User) error {
	logger.Logger.WithFields(logrus.Fields{
		"userID":          user.ID,
		"passport_number": user.PassportNumber,
		"username":        user.Surname,
		"name":            user.Name,
		"patronymic":      user.Patronymic,
		"address":         user.Address,
	}).Debug("Updating user data")

	user.UpdatedAt = time.Now()
	query := `UPDATE users SET passport_number=$1, surname=$2, name=$3, patronymic=$4, address=$5, updated_at=$6 WHERE id=$7`

	_, err := r.db.Exec(ctx, query, user.PassportNumber, user.Surname, user.Name, user.Patronymic, user.Address, user.UpdatedAt, user.ID)
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"userID": user.ID,
			"error":  err,
		}).Error("An error occurred while updating user data")
		return err
	}
	logger.Logger.WithFields(logrus.Fields{
		"userID": user.ID,
	}).Info("User data has been successfully updated")
	return nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, id int) error {
	logger.Logger.WithFields(logrus.Fields{
		"userID": id,
	}).Debug("Deleting a user")

	query := `DELETE FROM users WHERE id=$1`
	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"userID": id,
			"error":  err,
		}).Error("An error occurred when deleting a user")
		return err
	}
	logger.Logger.WithFields(logrus.Fields{
		"userID": id,
	}).Info("The user was successfully deleted")

	return nil
}

func (r *UserRepository) GetUsers(ctx context.Context, filter map[string]interface{}, limit, offset int) ([]models.User, error) {
	logger.Logger.WithFields(logrus.Fields{
		"filter": filter,
		"limit":  limit,
		"offset": offset,
	}).Debug("Obtaining users with the ability to filter and paginate")

	var argID int = 1
	query := "SELECT id, passport_number, surname, name, patronymic, address, created_at, updated_at FROM users WHERE true"
	args := []interface{}{}

	for key, val := range filter {
		query += fmt.Sprintf(" AND %s = $%d", key, argID)
		args = append(args, val)
		argID++
	}

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argID, argID+1)
	args = append(args, limit, offset)
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"error": err,
		}).Error("An error occurred while retrieving users")
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.PassportNumber, &user.Surname, &user.Name, &user.Patronymic, &user.Address, &user.CreatedAt, &user.UpdatedAt); err != nil {
			logger.Logger.WithFields(logrus.Fields{
				"error": err,
			}).Error("An error occurred while scanning user strings")
			return nil, err
		}

		users = append(users, user)
	}
	if rows.Err() != nil {
		logger.Logger.WithFields(logrus.Fields{
			"error": rows.Err(),
		}).Error("An error occurred when trying to iterate over rows with users")
		return nil, rows.Err()
	}

	logger.Logger.WithFields(logrus.Fields{
		"filter": filter,
		"limit":  limit,
		"offset": offset,
		"count":  len(users),
	}).Info("User information successfully received")

	return users, nil
}
