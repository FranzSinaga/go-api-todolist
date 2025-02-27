package main

import (
	"database/sql"
	fmtLog "log"
	"net/http"
	"os"

	"go-todos-api/middleware"
	"go-todos-api/src/api"

	"github.com/go-kit/log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB

func main() {
	var err error

	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	db, err = sql.Open("mysql", "dbeaver:dbeaver@tcp(127.0.0.1:3306)/todolist?charset=utf8")
	if err != nil {
		logger.Log("ERROR MESSAGE", err)
	}

	_, err = db.Query("SELECT * FROM todos")
	if err != nil {
		panic(err.Error())
	}
	fmtLog.Printf("DB Connection Successfully")
	defer db.Close()

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
	port := ":8000"
	fmtLog.Printf("Server connected to port :8000")
	http.ListenAndServe(port, router)

}
