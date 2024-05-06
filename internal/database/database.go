package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"github.com/pyKis/go_final_project/configs/models"
	"os"
	_"modernc.org/sqlite"
)

var db *sql.DB

func getDbFilePath() string {
	dbFilePath := "scheduler.db"

	envDbFilePath := os.Getenv("TODO_DBFILE")
	if len(envDbFilePath) > 0 {
		dbFilePath = envDbFilePath
	}

	return dbFilePath
}

func createDbFile(dbFilePath string) (*sql.DB, error) {
	_, err := os.Create(dbFilePath)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite", dbFilePath)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func createTable(db *sql.DB) {
	_, err := db.Exec(
		"CREATE TABLE IF NOT EXISTS `scheduler` (`id` INTEGER PRIMARY KEY AUTOINCREMENT, `date` VARCHAR(8) NULL, `title` VARCHAR(64) NOT NULL, `comment` VARCHAR(255) NULL, `repeat` VARCHAR(128) NULL)")
	if err != nil {
		log.Fatal(err)
	}
}

func InstallDb() {
	dbFilePath := getDbFilePath()
	_, err := os.Stat(dbFilePath)

	if err != nil {
		db, err = createDbFile(dbFilePath)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		db, err = sql.Open("sqlite", dbFilePath)
	}

	if err != nil {
		log.Fatal(err)
	}
	createTable(db)
}

func InsertTask(task models.Task) (int, error) {
	result, err := db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func ReadTasks() ([]models.Task, error) {
	var tasks []models.Task

	rows, err := db.Query("SELECT * FROM scheduler ORDER BY date")
	if err != nil {
		return []models.Task{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return []models.Task{}, err
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return []models.Task{}, err
	}

	if tasks == nil {
		tasks = []models.Task{}
	}

	return tasks, nil
}

func SearchTasks(search string) ([]models.Task, error) {
	var tasks []models.Task

	search = fmt.Sprintf("%%%s%%", search)
	rows, err := db.Query("SELECT * FROM scheduler WHERE title LIKE :search OR comment LIKE :search ORDER BY date",
		sql.Named("search", search))
	if err != nil {
		return []models.Task{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return []models.Task{}, err
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return []models.Task{}, err
	}

	if tasks == nil {
		tasks = []models.Task{}
	}

	return tasks, nil
}

func SearchTasksByDate(date string) ([]models.Task, error) {
	var tasks []models.Task

	rows, err := db.Query("SELECT * FROM scheduler WHERE date = :date",
		sql.Named("date", date))
	if err != nil {
		return []models.Task{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return []models.Task{}, err
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return []models.Task{}, err
	}

	if tasks == nil {
		tasks = []models.Task{}
	}

	return tasks, nil
}

func ReadTask(id string) (models.Task, error) {
	var task models.Task

	row := db.QueryRow("SELECT * FROM scheduler WHERE id = :id",
		sql.Named("id", id))
	if err := row.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
		return models.Task{}, err
	}

	return task, nil
}

func UpdateTask(task models.Task) (models.Task, error) {
	result, err := db.Exec("UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", task.Id))
	if err != nil {
		return models.Task{}, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return models.Task{}, err
	}

	if rowsAffected == 0 {
		return models.Task{}, errors.New("failed to update")
	}

	return task, nil
}

func DeleteTaskDb(id string) error {
	result, err := db.Exec("DELETE FROM scheduler WHERE id = :id",
		sql.Named("id", id))
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("failed to delete")
	}

	return err
}