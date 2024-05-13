package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/pyKis/go_final_project/configs/models"
	"github.com/pyKis/go_final_project/internal/storage"
)

func responseWithError(w http.ResponseWriter, errorText string, err error) {
	errorResponse := models.ErrorResponse{
		Error: fmt.Errorf("%s: %w", errorText, err).Error()}
	errorData, _ := json.Marshal(errorResponse)
	w.WriteHeader(http.StatusBadRequest)
	_, err = w.Write(errorData)

	if err != nil {
		http.Error(w, fmt.Errorf("error: %w", err).Error(), http.StatusBadRequest)
	}
}

func TaskAddPost(w http.ResponseWriter, r *http.Request) {
	var taskData models.Task
	var buffer bytes.Buffer

	if _, err := buffer.ReadFrom(r.Body); err != nil {
		responseWithError(w, "body getting error", err)
		return
	}

	if err := json.Unmarshal(buffer.Bytes(), &taskData); err != nil {
		responseWithError(w, "JSON encoding error", err)
		return
	}

	if len(taskData.Date) == 0 {
		taskData.Date = time.Now().Format(models.DatePattern)
	} else {
		date, err := time.Parse(models.DatePattern, taskData.Date)
		if err != nil {
			responseWithError(w, "bad data format", err)
			return
		}

		if date.Before(time.Now()) {
			taskData.Date = time.Now().Format(models.DatePattern)
		}
	}

	if len(taskData.Title) == 0 {
		responseWithError(w, "invalid title", errors.New("title is empty"))
		return
	}

	if len(taskData.Repeat) > 0 {
		if _, err := NextDate(time.Now(), taskData.Date, taskData.Repeat); err != nil {
			responseWithError(w, "invalid repeat format", errors.New("no such format"))
			return
		}
	}

	taskId, err := storage.InsertTask(taskData)
	if err != nil {
		responseWithError(w, "failed to create task", err)
		return
	}

	taskIdData, err := json.Marshal(models.TaskIdResponse{Id: uint(taskId)})
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(taskIdData)
	log.Println(fmt.Sprintf("Added task with id=%d", taskId))

	if err != nil {
		responseWithError(w, "writing task id error", err)
	}
}

func TasksReadGet(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")

	var tasks []models.Task
	if len(search) > 0 {
		date, err := time.Parse("02.01.2006", search)
		if err != nil {
			tasks, err = storage.SearchTasks(search)
		} else {
			tasks, err = storage.SearchTasksByDate(date.Format(models.DatePattern))
		}
	} else {
		err := errors.New("")
		tasks, err = storage.ReadTasks()
		if err != nil {
			responseWithError(w, "failed to get tasks", err)
			return
		}
	}

	tasksData, err := json.Marshal(models.Tasks{Tasks: tasks})
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(tasksData)
	log.Println(fmt.Sprintf("Read %d tasks", len(tasks)))

	if err != nil {
		responseWithError(w, "writing tasks error", err)
	}
}

func TaskReadGet(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	task, err := storage.ReadTask(id)
	if err != nil {
		responseWithError(w, "failed to get task", err)
		return
	}

	tasksData, err := json.Marshal(task)
	if err != nil {
		log.Panicln("JSON encoding error", err)
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(tasksData)

	if err != nil {
		responseWithError(w, "writing task error", err)
	}
}

func TaskUpdatePut(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	var buffer bytes.Buffer

	if _, err := buffer.ReadFrom(r.Body); err != nil {
		responseWithError(w, "body getting error", err)
		return
	}

	if err := json.Unmarshal(buffer.Bytes(), &task); err != nil {
		responseWithError(w, "JSON encoding error", err)
		return
	}

	if len(task.Id) == 0 {
		responseWithError(w, "invalid id", errors.New("id is empty"))
		return
	}

	if _, err := strconv.Atoi(task.Id); err != nil {
		responseWithError(w, "invalid id", err)
		return
	}

	if _, err := time.Parse(models.DatePattern, task.Date); err != nil {
		responseWithError(w, "invalid date", err)
		return
	}

	if len(task.Title) == 0 {
		responseWithError(w, "invalid title", errors.New("title is empty"))
		return
	}

	if len(task.Repeat) > 0 {
		if _, err := NextDate(time.Now(), task.Date, task.Repeat); err != nil {
			responseWithError(w, "invalid repeat format", errors.New("no such format"))
			return
		}
	}

	_, err := storage.UpdateTask(task)
	if err != nil {
		responseWithError(w, "invalid title", errors.New("failed to update task"))
		return
	}

	taskIdData, err := json.Marshal(task)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(taskIdData)
	log.Println(fmt.Sprintf("Updated task with id=%s", task.Id))

	if err != nil {
		responseWithError(w, "updating task error", err)
		return
	}
}

func TaskDonePost(w http.ResponseWriter, r *http.Request) {
	task, err := storage.ReadTask(r.URL.Query().Get("id"))
	if err != nil {
		responseWithError(w, "failed to get task", err)
		return
	}

	if len(task.Repeat) == 0 {
		err = storage.DeleteTask(task.Id)
		if err != nil {
			responseWithError(w, "failed to delete task", err)
			return
		}
	} else {
		task.Date, err = NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			responseWithError(w, "failed to get next date", err)
			return
		}

		_, err = storage.UpdateTask(*task)
		if err != nil {
			responseWithError(w, "failed to update task", err)
			return
		}
	}

	tasksData, err := json.Marshal(struct{}{})
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(tasksData)
	log.Println(fmt.Sprintf("Done task with id=%s", task.Id))

	if err != nil {
		responseWithError(w, "writing task error", err)
	}
}

func TaskDelete(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	err := storage.DeleteTask(id)
	if err != nil {
		responseWithError(w, "failed to delete task", err)
		return
	}

	tasksData, err := json.Marshal(struct{}{})
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(tasksData)
	log.Println(fmt.Sprintf("Deleted task with id=%s", id))

	if err != nil {
		responseWithError(w, "writing task error", err)
		return
	}
}