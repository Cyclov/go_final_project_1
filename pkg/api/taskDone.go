package api

import (
	"fmt"
	"net/http"
	"scheduler/pkg/db"
	"time"
)

func taskDoneHandler(w http.ResponseWriter, r *http.Request) {
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

	if task.Repeat == "" {

		err = db.DeleteTask(id)
		if err != nil {
			JsonError(w, http.StatusInternalServerError, err.Error())
			return
		}

	} else {

		task.Date, err = NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			JsonError(w, http.StatusInternalServerError, err.Error())
			return
		}

		err = db.UpdateTask(task, true)

		if err != nil {
			JsonError(w, http.StatusInternalServerError, err.Error())
			return
		}

	}

	WriteJson(w, struct{}{})

}
