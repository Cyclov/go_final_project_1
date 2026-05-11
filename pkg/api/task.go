package api

import (
	"net/http"
)

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodPost:
		addTaskHandler(w, r)

	case http.MethodGet:
		GetTaskHandler(w, r)

	case http.MethodPut:
		UpdateTaskHandler(w, r)

	case http.MethodDelete:
		DeleteTaskHandler(w, r)

	}
}
