package api

import (
	"fmt"
	"net/http"

	"scheduler/pkg/db"
)

func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		err := fmt.Errorf("Parametr \"id\" is requaire")
		JsonError(w, http.StatusInternalServerError, err.Error())
		return

	}

	err := db.DeleteTask(id)

	if err != nil {
		JsonError(w, http.StatusInternalServerError, err.Error())
		return
	}

	WriteJson(w, struct{}{})
}
