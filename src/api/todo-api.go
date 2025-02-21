package api

import (
	"database/sql"
	"encoding/json"
	"go-todos-api/src/model"
	"io"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB

func SetDB(database *sql.DB) {
	db = database
}

func GetTodos(w http.ResponseWriter, r *http.Request) {
	var todos []model.Todos
	rows, err := db.Query("SELECT id, title, description, status FROM todos WHERE status NOT LIKE 'Deleted'")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	hasData := false
	for rows.Next() {
		var todo model.Todos
		if err := rows.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Status); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		todos = append(todos, todo)
		hasData = true
	}

	if !hasData {
		todos = []model.Todos{}
	}

	response := map[string]interface{}{
		"data":    todos,
		"error":   false,
		"status":  http.StatusOK,
		"message": "Success get all data todos",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetTodo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var todo model.Todos
	err := db.QueryRow("SELECT id, title, description, status FROM todos WHERE id = ?", id).Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			response := map[string]interface{}{
				"data":    nil,
				"error":   true,
				"status":  http.StatusNotFound,
				"message": "No user found with the given ID",
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)

			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"data":    todo,
		"error":   false,
		"status":  http.StatusOK,
		"message": "",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	var todo model.CreateTodoRequest

	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil { // cek apakah ada kesalahan
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := db.Exec("INSERT INTO todos (title, description, status) VALUES (?, ?, ?)", todo.Title, todo.Description, "To Do")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"data":    id,
		"error":   false,
		"status":  http.StatusOK,
		"message": "Successfully created Todo",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	// Baca isi r.Body dan simpan ke dalam variabel
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var todo model.CreateTodoRequest
	err = json.Unmarshal(body, &todo)

	if err != nil { // cek apakah ada kesalahan
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validasi input menggunakan validator
	if err := todo.Validate(); err != nil {
		response := map[string]interface{}{
			"data":    nil,
			"error":   true,
			"status":  http.StatusBadRequest,
			"message": "Validation failed: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}
	log.Print(todo.Title)
	log.Print(todo.Description)

	result, err := db.Exec("UPDATE todos SET title = ?, description = ?, status = ? WHERE id = ?", todo.Title, todo.Description, todo.Status, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = result.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"data":    id,
		"error":   false,
		"status":  http.StatusOK,
		"message": "Successfully updated Todo",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func HardDeleteTodo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	_, err := db.Exec("DELETE FROM todos WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"data":    id,
		"error":   false,
		"status":  http.StatusOK,
		"message": "Successfully deleted todo",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
