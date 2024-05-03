package database

import (
	"errors"
	"fmt"

	"log"
	"os"

	"github.com/joho/godotenv"
	
	"github.com/pyKis/go_final_project/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	_"modernc.org/sqlite"
)


func init() {
    if err := godotenv.Load(".env"); err != nil {
        log.Print("No .env file found")
    }
}

var Db *gorm.DB
func getDBPath() string {
	 DBFile := os.Getenv("TODO_DBFILE")
	 if len(DBFile) < 0 {
		errors.New("TODO_PORT not set")
	}

	return DBFile
}



func ConnectDB() {
	// создаём подключение к базе данных.
	// В &gorm.Config настраивается логер,
	// который будет сохранять информацию
	// обо всех активностях с базой данных.
	fmt.Println(getDBPath())
	db, err := gorm.Open(sqlite.Open(getDBPath()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database.\n", err)
	}

	log.Println("connected")
	db.Logger = logger.Default.LogMode(logger.Info)

	log.Println("running migration")
	db.AutoMigrate(&models.Task{})

	Db = db
}

func ReadTasks() ([]models.Task, error) {
	var tasks []models.Task
	result := Db.Order("date").Limit(30).Find(&tasks)
	return tasks, result.Error
}

func ReadTaskByID(id uint) (models.Task, error) {
	var tasks models.Task
	result := Db.Limit(1).Find(&tasks, models.Task{ID: id})
	if result.RowsAffected == 0 {
		return tasks, errors.New("No record found")
	}
	return tasks, result.Error
}