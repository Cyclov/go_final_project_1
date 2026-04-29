package main

import (
	"log"
	"os"

	"github.com/Cyclov/go_final_project_1/pkg/db"
)

func main() {
	dbPath := getDBPath()

	err := db.Init(dbPath)

	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

	defer db.Close()

}

func getDBPath() string {
	path, exists := os.LookupEnv("TODO_DBFILE")
	if exists && path == "" { //Проверили наличие переменной
		path = "scheduler.db" // Не нашли, вернули значение по умолчанию
	}

	return path
}
