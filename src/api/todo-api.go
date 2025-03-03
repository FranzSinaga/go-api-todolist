package api

import (
	"database/sql"
	"encoding/json"
	api_helper "go-todos-api/src/api/helper"
	"go-todos-api/src/model"
	"io"
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
	rows, err := db.Query("SELECT id, title, description, status, created_date, updated_date FROM todos WHERE status NOT LIKE 'Deleted'")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	hasData := false
	for rows.Next() {
		var todo model.Todos
		if err := rows.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Status, &todo.CreatedDate, &todo.UpdatedDate); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		todos = append(todos, todo)
		hasData = true
	}

	if !hasData {
		todos = []model.Todos{}
	}

	response := api_helper.SetResponse(todos, false, http.StatusOK, "Success get all data todos")
	responseJson, _ := json.Marshal(response)
	api_helper.SendResponse(w, responseJson, http.StatusOK)
}

func GetTodo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var todo model.Todos
	err := db.QueryRow("SELECT id, title, description, status FROM todos WHERE id = ?", id).Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			response := api_helper.SetResponse(nil, true, http.StatusOK, "No user found with the given ID")
			responseJson, _ := json.Marshal(response)
			api_helper.SendResponse(w, responseJson, http.StatusOK)

			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := api_helper.SetResponse(todo, false, http.StatusOK, "Success get data")
	responseJson, _ := json.Marshal(response)
	api_helper.SendResponse(w, responseJson, http.StatusOK)
}

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	var todo model.CreateTodoRequest

	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil { // cek apakah ada kesalahan
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := todo.Validate(); err != nil {
		response := api_helper.SetResponse(nil, true, http.StatusBadRequest, "Validation failed: "+err.Error())
		responseJson, _ := json.Marshal(response)
		api_helper.SendResponse(w, responseJson, http.StatusBadRequest)
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

	response := api_helper.SetResponse(id, false, http.StatusOK, "Successfully created Todo")
	responseJson, _ := json.Marshal(response)
	api_helper.SendResponse(w, responseJson, http.StatusOK)
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

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := todo.Validate(); err != nil {
		response := api_helper.SetResponse(nil, true, http.StatusBadRequest, "Validation failed: "+err.Error())
		responseJson, _ := json.Marshal(response)
		api_helper.SendResponse(w, responseJson, http.StatusBadRequest)
		return
	}

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

	response := api_helper.SetResponse(id, false, http.StatusOK, "Successfully updated Todo")
	responseJson, _ := json.Marshal(response)
	api_helper.SendResponse(w, responseJson, http.StatusOK)
}

func HardDeleteTodo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	_, err := db.Exec("DELETE FROM todos WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := api_helper.SetResponse(id, false, http.StatusOK, "Successfully deleted todo")
	responseJson, _ := json.Marshal(response)
	api_helper.SendResponse(w, responseJson, http.StatusOK)
}
