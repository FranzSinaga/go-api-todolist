package main

import (
	fmtLog "log"
	"net/http"
	"os"

	"go-todos-api/config"
	"go-todos-api/middleware"
	"go-todos-api/src/api"

	"github.com/go-kit/log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	var err error

	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	if os.Getenv("ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			fmtLog.Println("Warning: No .env file found")
		}
	}

	db := config.DBConnection(logger)
	api.SetDB(db)
	router := mux.NewRouter()

	router.Use(middleware.LoggingMiddleware(logger))
	router.Use(middleware.ErrorHandlingMiddleware)

	// Define API routes
	router.HandleFunc("/api/todos", api.GetTodos).Methods("GET")
	router.HandleFunc("/api/todo/{id}", api.GetTodo).Methods("GET")
	router.HandleFunc("/api/todo/create", api.CreateTodo).Methods("POST")
	router.HandleFunc("/api/todo/update/{id}", api.UpdateTodo).Methods("PUT")
	router.HandleFunc("/api/todo/hard-delete/{id}", api.HardDeleteTodo).Methods("DELETE")

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmtLog.Println("Server connected to port", port)
	err = http.ListenAndServe(":"+port, router)
	if err != nil {
		logger.Log("level", "error", "message", err.Error())
		panic(err)
	}
}
