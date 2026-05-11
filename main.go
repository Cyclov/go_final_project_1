package main

import (
	"log"
	"os"

	"scheduler/pkg/db"
	"scheduler/pkg/server"
)

func main() {
	dbPath := getDBPath()
	serverPort := getServerPort()
	err := db.Init(dbPath)

	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

	defer db.Close()

	srv := server.Init(serverPort)

	log.Println("Server running on http://localhost" + srv.Addr)

	log.Fatal(srv.ListenAndServe())

}

func getDBPath() string {
	path, exists := os.LookupEnv("TODO_DBFILE")
	if !exists || path == "" { //Проверили наличие переменной
		path = "scheduler.db" // Не нашли, вернули значение по умолчанию
	}

	return path
}

func getServerPort() string {
	port, exists := os.LookupEnv("TODO_PORT")
	if !exists || port == "" { //Проверили наличие переменной
		port = "7540" // Не нашли, вернули значение по умолчанию
	}

	return port
}
