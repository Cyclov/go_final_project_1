package api

import "net/http"

func Init(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/nextdate", nextDayHandler)
	mux.HandleFunc("GET /api/task", taskHandler)
	mux.HandleFunc("POST /api/task", taskHandler)
}
