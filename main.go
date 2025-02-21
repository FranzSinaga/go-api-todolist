package main

import (
	"database/sql"
	"log"
	"net/http"

	"go-todos-api/src/api"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB

func main() {
	var err error

	db, err = sql.Open("mysql", "dbeaver:dbeaver@tcp(127.0.0.1:3306)/todolist")

	if err != nil {
		log.Fatal(err)
	}
	log.Print("DB Connection Successfully")
	defer db.Close()

	api.SetDB(db)
	router := mux.NewRouter()

	// Define API routes
	router.HandleFunc("/todos", api.GetTodos).Methods("GET")
	router.HandleFunc("/todo/{id}", api.GetTodo).Methods("GET")
	router.HandleFunc("/todo/create", api.CreateTodo).Methods("POST")
	router.HandleFunc("/todo/update/{id}", api.UpdateTodo).Methods("PUT")

	router.HandleFunc("/todo/hard-delete/{id}", api.HardDeleteTodo).Methods("DELETE")

	log.Print("Server connected to port http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}
