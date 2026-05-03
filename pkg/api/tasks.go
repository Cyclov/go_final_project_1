package api

import (
	"net/http"
	"scheduler/pkg/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.Tasks(50)
	if err != nil {
		JsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJson(w, TasksResp{
		Tasks: tasks,
	})
}
