package api

import (
	"encoding/json"
	"net/http"
)

func Init(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/nextdate", nextDayHandler)
	mux.HandleFunc("GET /api/task", auth(taskHandler))
	mux.HandleFunc("PUT /api/task", auth(taskHandler))
	mux.HandleFunc("POST /api/task", auth(taskHandler))
	mux.HandleFunc("DELETE /api/task", auth(taskHandler))
	mux.HandleFunc("POST /api/task/done", auth(taskDoneHandler))
	mux.HandleFunc("GET /api/tasks", auth(tasksHandler))
	mux.HandleFunc("POST /api/signin", singinHandler)
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
