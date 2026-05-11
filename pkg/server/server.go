package server

import (
	"net/http"
	"time"

	"scheduler/pkg/api"
)

const webDir = "web"

func Init(port string) *http.Server {

	router := http.NewServeMux()
	strPort := ":" + port

	router.Handle("GET /", http.FileServer(http.Dir(webDir)))
	api.Init(router)

	return &http.Server{
		Addr:         strPort,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

}
