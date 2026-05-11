package api

import (
	"net/http"

	"scheduler/pkg/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {

	search := r.FormValue("search")
	tasks, err := db.Tasks(50, search)
	if err != nil {
		JsonError(w, http.StatusInternalServerError, err.Error())
		return
	}
	WriteJson(w, TasksResp{
		Tasks: tasks,
	})
}
