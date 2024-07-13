package users

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/DarRo9/time-tracker/internal/apperrors"
)

const infoURL = "http://localhost:8081/info"

type service struct {
	Repository
	timeout time.Duration
}

func NewService(r Repository) Service {
	return &service{
		Repository: r,
		timeout:    time.Duration(2) * time.Second}
}

func (s *service) createUser(c context.Context, user *user) error {
	slog.Debug("Creating new user in service", "user", user)
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer func() {
		cancel()
	}()

	passport := strings.Fields(user.Passport)
	if len(passport) != 2 || len(passport[0]) != 4 || len(passport[1]) != 6 {
		slog.Error("Invalid passport format", "passport", user.Passport)
		return &apperrors.BadRequestError{Message: "Invalid passport format"}
	}
	if _, err := strconv.Atoi(passport[0]); err != nil {
		slog.Error("Invalid passport series", "passport", user.Passport)
		return &apperrors.BadRequestError{Message: "Invalid  passport series"}
	}
	if _, err := strconv.Atoi(passport[1]); err != nil {
		slog.Error("Invalid passport number", "passport", user.Passport)
		return &apperrors.BadRequestError{Message: "Invalid passport number"}
	}

	url := fmt.Sprintf("%s?passportSerie=%s&passportNumber=%s", infoURL, passport[0], passport[1])

	resp, err := http.Get(url)
	if err != nil {
		slog.Error("Error while making request to info API", "error", err)
		return &apperrors.ExternalAPIError{Message: "Error while making request to info API"}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			slog.Error("Error reading external api responce body", "error", err)
			return &apperrors.ExternalAPIError{Message: "Error while making request to info API"}
		}

		slog.Error("Error while making request to info API", "body", string(body))
		return &apperrors.ExternalAPIError{Message: "Error while making request to info API"}
	}

	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		slog.Error("Error unmarshaling info from external API", "error", err)
		return &apperrors.ExternalAPIError{Message: "Error unmarshaling info from external API"}
	}

	if len(user.Address) == 0 || len(user.Name) == 0 || len(user.Surname) == 0 {
		slog.Error("Got not complete data from external API", "user", user)
		return &apperrors.ExternalAPIError{Message: "Got not complete data from external API"}
	}

	if err := s.Repository.createUser(ctx, user); err != nil {
		return err
	}

	return nil
}

func (s *service) getUsers(ctx context.Context, search *searchReq) (*[]user, error) {
	slog.Debug("Calculating offset", "search info", search)
	search.offset = (search.page - 1) * search.pageSize
	return s.Repository.getUsers(ctx, search)
}

func (s *service) removeUser(ctx context.Context, id *userId) error {
	return s.Repository.removeUser(ctx, id)
}

func (s *service) updateUser(ctx context.Context, id *userId, us *updateUserReq) error {
	updates := make(map[string]interface{})
	if us.Passport != nil {
		updates["passport"] = *us.Passport
	}
	if us.Surname != nil {
		updates["surname"] = *us.Surname
	}
	if us.Name != nil {
		updates["name"] = *us.Name
	}
	if us.Patronymic != nil {
		updates["patronymic"] = *us.Patronymic
	}
	if us.Address != nil {
		updates["address"] = *us.Address
	}
	slog.Debug("New update request", "request", updates)

	return s.Repository.updateUser(ctx, id, &updates)
}
