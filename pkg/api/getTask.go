package api

import (
	"fmt"
	"net/http"
	"scheduler/pkg/db"
)

func GetTaskHandler(w http.ResponseWriter, r *http.Request) {

	id := r.FormValue("id")
	if id == "" {
		err := fmt.Errorf("Parametr \"id\" is requaire")
		JsonError(w, http.StatusInternalServerError, err.Error())
		return

	}

	task, err := db.GetTask(id)

	if err != nil {
		JsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	WriteJson(w, task)

}
