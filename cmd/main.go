package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	"github.com/pyKis/go_final_project/internal/handlers"
	"github.com/pyKis/go_final_project/internal/storage"
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
	fmt.Println("Connect DB")
	storage.InstallDb()
	fmt.Println("DB connected")
	myHandler := chi.NewRouter()
	fmt.Println("Register handlers")
	myHandler.Mount("/", http.FileServer(http.Dir(webDir)))



	myHandler.Get("/api/nextdate", handlers.NextDateReadGET)
	myHandler.Post("/api/task", handlers.TaskAddPost)
	myHandler.Get("/api/tasks", handlers.TasksReadGet)
	myHandler.Get("/api/task", handlers.TaskReadGet)
	myHandler.Put("/api/task", handlers.TaskUpdatePut)
	myHandler.Post("/api/task/done", handlers.TaskDonePost)
	myHandler.Delete("/api/task", handlers.TaskDelete)
	myHandler.Post("/api/signin", handlers.SignInPOST)

	fmt.Printf("Starting server on port %s\n", getPort())

	s:=&http.Server{
		Addr:	getPort(),
		Handler: myHandler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}