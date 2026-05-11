package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"scheduler/pkg/db"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		JsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	if task.Title == "" {
		JsonError(w, http.StatusBadRequest, "title is required")
		return
	}

	now := time.Now()
	today := now.Format(db.DateFormat)

	if task.Date == "" {
		task.Date = today
	}

	taskDate, err := time.Parse(db.DateFormat, task.Date)
	if err != nil {
		JsonError(w, http.StatusBadRequest, fmt.Sprintf("invalid date format: %s", task.Date))
		return
	}

	if task.Repeat != "" {
		_, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			JsonError(w, http.StatusBadRequest, fmt.Sprintf("invalid repeat rule: %v", err))
			return
		}
	}

	if taskDate.Before(now.Truncate(24 * time.Hour)) {
		if task.Repeat == "" {
			task.Date = today
		} else {
			next, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				JsonError(w, http.StatusBadRequest, fmt.Sprintf("invalid repeat rule: %v", err))
				return
			}
			task.Date = next
		}
	}

	id, err := db.AddTask(&task)
	if err != nil {
		JsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)

	msgStr := struct {
		Id int64 `json:"id"`
	}{
		Id: id,
	}
	json.NewEncoder(w).Encode(msgStr)

}
