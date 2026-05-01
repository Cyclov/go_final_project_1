package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"scheduler/pkg/db"
	"time"
)

func jsonError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	type MsgJson struct {
		Error string `json:"error"`
	}

	msgStr := MsgJson{
		Error: msg,
	}
	json.NewEncoder(w).Encode(msgStr)

}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		jsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	if task.Title == "" {
		jsonError(w, http.StatusBadRequest, "title is required")
		return
	}

	now := time.Now()
	today := now.Format("20060102")

	if task.Date == "" {
		task.Date = today
	}

	taskDate, err := time.Parse("20060102", task.Date)
	if err != nil {
		jsonError(w, http.StatusBadRequest, fmt.Sprintf("invalid date format: %s", task.Date))
		return
	}

	if task.Repeat != "" {
		_, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			jsonError(w, http.StatusBadRequest, fmt.Sprintf("invalid repeat rule: %v", err))
			return
		}
	}

	if taskDate.Before(now.Truncate(24 * time.Hour)) {
		if task.Repeat == "" {
			task.Date = today
		} else {
			next, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				jsonError(w, http.StatusBadRequest, fmt.Sprintf("invalid repeat rule: %v", err))
				return
			}
			task.Date = next
		}
	}

	id, err := db.AddTask(&task)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, `{"id":"%d"}`, id)
}
