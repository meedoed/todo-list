package main

import (
	"context"
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
	dburl = "postgres://postgres:admin12345@localhost:5432/postgres"
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetReportCaller(true)
}

func main() {
	dbpool, err := pgxpool.New(context.Background(), dburl)
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

	log.Fatal(http.ListenAndServe("localhost:8000", router))
}

func createTask() {

}

func updateTask() {

}

func deleteTask() {

}

func getCompletedTasks() {
}

func getIncompleteTasks() {
}
