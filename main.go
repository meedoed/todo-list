package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"os"
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

	router.HandleFunc("/todo", createTask).Methods("POST")

	log.Fatal(http.ListenAndServe("localhost:8000", router))
}

func createTask(w http.ResponseWriter, r *http.Request) {
	q := `INSERT INTO todo_list (description, status) VALUES ($1, $2) RETURNING id;`
	desctiption := r.FormValue("description")
	logrus.WithFields(logrus.Fields{"description": desctiption}).Info("Add new TodoTask. Saving to database.")
	todo := &TodoModel{Description: desctiption, Completed: false}
	if err := dbpool.QueryRow(context.Background(), q, todo.Description, todo.Completed).Scan(&todo.Id); err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

func updateTask(w http.ResponseWriter, r *http.Request) {

}

func deleteTask(w http.ResponseWriter, r *http.Request) {

}

func getCompletedTasks(w http.ResponseWriter, r *http.Request) {
}

func getIncompleteTasks(w http.ResponseWriter, r *http.Request) {
}
