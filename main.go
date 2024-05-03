package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	"github.com/pyKis/go_final_project/database"
	"github.com/pyKis/go_final_project/hendlers"
)

const(
	webDir = "./web"
)

func init() {
    if err := godotenv.Load(".env"); err != nil {
        log.Print("No .env file found")
    }
}

func getPort() string {
	port := os.Getenv("TODO_PORT")	
	if len(port) < 0 {
		errors.New("TODO_PORT not set")
	}
	return ":" + port
}
func main() {
	database.ConnectDB()
	myHandler := chi.NewRouter()

	myHandler.Mount("/", http.FileServer(http.Dir(webDir)))

	myHandler.Get("/api/nextdate", handlers.NextDateGET)
	myHandler.Post("/api/task", handlers.TaskPost)
	myHandler.Get("/api/tasks", handlers.TasksRead)
	myHandler.Get("/api/task", handlers.TaskReadByID)
	myHandler.Put("/api/task", handlers.TaskUpdate)
	myHandler.Post("/api/task/done", handlers.TaskDone)
	myHandler.Delete("/api/task", handlers.TaskDelete)

	s:=&http.Server{
		Addr:	getPort(),
		Handler: myHandler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}