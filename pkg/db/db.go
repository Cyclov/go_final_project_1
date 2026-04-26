package db

import (
	"os"
)

func getDBPath() string {
	path, exists := os.LookupEnv("TODO_DBFILE")
	if exists && path == "" { //Проверили наличие переменной
		path = "scheduler.db" // Не нашли, вернули значение по умолчанию
	}

	return path
}
