package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

type TodoModel struct {
	Id          int
	Description string
	Completed   bool
}

var (
	dburl  = "postgres://postgres:admin12345@localhost:5432/postgres"
	dbpool *pgxpool.Pool
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetReportCaller(true)
}

func main() {
	var err error
	dbpool, err = pgxpool.New(context.Background(), dburl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	fmt.Println("hello from todolist API")

	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "hello from todolist API")
	})

	router.HandleFunc("/todo", CreateTask).Methods("POST")

	log.Fatal(http.ListenAndServe("localhost:8000", router))
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	q := `INSERT INTO todo_list (description, status) VALUES ($1, $2) RETURNING id;`
	description := r.FormValue("description")
	logrus.WithFields(logrus.Fields{"description": description}).Info("Add new TodoTask. Saving to database.")
	todo := &TodoModel{Description: description, Completed: false}
	if err := dbpool.QueryRow(context.Background(), q, todo.Description, todo.Completed).Scan(&todo.Id); err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	if err := GetTaskByID(id); !err {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"updated": false, "error": "Record Not Found"}`)
	} else {
		completed, _ := strconv.ParseBool(r.FormValue("completed"))
		logrus.WithFields(logrus.Fields{"Id": id, "Completed": completed}).Info("Updating TodoItem")
		todo := &TodoModel{}
		todo.Completed = completed
		dbpool.QueryRow(context.Background(), `UPDATE todo_list SET completed = $1 WHERE id = $2`, todo.Completed, todo.Id)
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"updated": true}`)
	}
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	if err := GetTaskByID(id); !err {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"deletedd": false, "error": "Record Not Found"}`)
	} else {
		logrus.WithFields(logrus.Fields{"Id": id}).Info("Deleting TodoItem")
		todo := &TodoModel{}
		dbpool.QueryRow(context.Background(), `DELETE FROM todo_list WHERE id = $1`, todo.Id)
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"deleted": true}`)
	}
}

func GetTaskByID(id int) bool {
	q := `SELECT description, status 
		  FROM todo_list
		  WHERE id = $1`
	todo := &TodoModel{}
	err := dbpool.QueryRow(context.Background(), q, id).Scan(todo.Description, todo.Completed)
	if err != nil {
		logrus.Warn(err)
		return false
	}
	return true
}

func GetCompletedTasks(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Get completed TodoItems")
	completedTodoTasks := GetTodoTasks(true)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(completedTodoTasks)
}

func GetIncompleteTasks(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Get incomplete TodoItems")
	incompleteTodoTasks := GetTodoTasks(false)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(incompleteTodoTasks)
}

func GetTodoTasks(completed bool) interface{} {
	//todos := &TodoModel{}
	//q := `SELECT * FROM todo_list WHERE completed = $1 RETURNING id, description, completed`
	//err := dbpool.QueryRow(context.Background(), q, completed).Scan(todos.Id, todos.)
}
