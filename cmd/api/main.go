package main

import (
	"errors"
	"go-order-producer/internal/database"
	"go-order-producer/internal/handlers"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		databaseUrl = "postgres://admin:admin@localhost:5432/producerdb?sslmode=disable"
	}
	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8080"
	}

	log.Printf("Starting server on port %s", serverPort)

	db, err := database.Connect(databaseUrl)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	log.Println("Connected to database %s", databaseUrl)

	taskStore := database.NewTaskStore(db)

	handler := handlers.NewHandlers(taskStore)

	e := echo.New()

	e.Use(middleware.CORS())
	e.Use(middleware.Logger())

	// Routes
	e.GET("/tasks", handler.GetAllTasks)
	e.POST("/tasks", handler.CreateTask)
	e.GET("/tasks/:id", handler.GetTaskById)
	e.PUT("/tasks/:id", handler.UpdateTask)
	e.DELETE("/tasks/:id", handler.DeleteTask)

	serverAddr := ":" + serverPort

	if err := e.Start(serverAddr); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
