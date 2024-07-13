package main

import (
	"github.com/DarRo9/time-tracker/db"
	_ "github.com/DarRo9/time-tracker/docs"
	"github.com/DarRo9/time-tracker/internal/tasks"
	"github.com/DarRo9/time-tracker/internal/users"
	"github.com/DarRo9/time-tracker/router"

	"log"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	file, err := os.Create("app.log")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env", "error", err)
	}

	if os.Getenv("DEBUG") == "true" {
		logger := slog.New(slog.NewJSONHandler(file, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		}))
		slog.SetDefault(logger)
	} else {
		logger := slog.New(slog.NewJSONHandler(file, &slog.HandlerOptions{
			AddSource: true,
		}))
		slog.SetDefault(logger)
	}

	dbConn, err := db.NewDatabase()
	if err != nil {
		log.Fatal("Error creating database connection", "error", err)
	}

	userHandler := users.NewHandler(users.NewService(users.NewRepository(dbConn)))
	taskHandler := tasks.NewHandler(tasks.NewService(tasks.NewRepository(dbConn)))

	router := router.NewRouter(userHandler, taskHandler)

	if err := router.Run(); err != nil {
		log.Fatal("Failed to run server", "error", err)
	}
}
