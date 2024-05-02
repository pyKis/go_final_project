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
	"github.com/pyKis/go_final_project/database"
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
	fmt.Println("DB connect")
	database.ConnectDB()
	fmt.Println("DB connected")
	myHandler := chi.NewRouter()
	myHandler.Mount("/", http.FileServer(http.Dir(webDir)))

	s:=&http.Server{
		Addr:	getPort(),
		Handler: myHandler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}