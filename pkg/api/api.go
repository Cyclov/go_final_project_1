package api

import (
	"encoding/json"
	"net/http"
)

func Init(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/nextdate", nextDayHandler)
	mux.HandleFunc("GET /api/tasks", tasksHandler)
	mux.HandleFunc("GET /api/task", taskHandler)
	mux.HandleFunc("PUT /api/task", taskHandler)
	mux.HandleFunc("POST /api/task", taskHandler)
	mux.HandleFunc("DELETE /api/task", taskHandler)
	mux.HandleFunc("POST /api/task/done", taskDoneHandler)
}

func JsonError(w http.ResponseWriter, status int, msg string) {
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

func WriteJson(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(v)
}
