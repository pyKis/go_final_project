package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

func init() {
    if err := godotenv.Load(".env"); err != nil {
        log.Print("No .env file found")
    }
}

func getPort(path string) string {
	port := 0
	envPort := os.Getenv("TODO_PORT")
	if len(envPort) > 0 {
		if eport, err := strconv.ParseInt(envPort, 10, 32); err == nil {
			port = int(eport)
		}
	}
	path = strings.ReplaceAll(strings.TrimPrefix(path, `../web/`), `\`, `/`)
	return fmt.Sprintf("http://localhost:%d/%s", port, path)
}
func main() {
	fmt.Println(getPort("../web/index.html"))
}