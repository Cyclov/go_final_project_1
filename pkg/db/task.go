package db

import (
	"database/sql"
	"fmt"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func GetTask(id string) (*Task, error) {

	resp := &Task{}

	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`

	err := db.QueryRow(query, id).Scan(&resp.ID, &resp.Date, &resp.Title, &resp.Comment, &resp.Repeat)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("No task with id %s", id)
		}

		return nil, err
	}

	return resp, nil
}

func UpdateTask(task *Task, UpdateOnlyDate bool) error {
	var resp sql.Result
	var err error

	if UpdateOnlyDate {
		query := `UPDATE scheduler SET date = ? WHERE id = ?`
		resp, err = db.Exec(query, task.Date, task.ID)
	} else {
		query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
		resp, err = db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	}

	if err != nil {
		return err
	}

	count, err := resp.RowsAffected()

	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("No task with id %s", task.ID)
	}

	return nil

}

func DeleteTask(id string) error {

	query := `Delete FROM scheduler WHERE id = ?`

	resp, err := db.Exec(query, id)

	if err != nil {
		return err
	}

	count, err := resp.RowsAffected()

	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("No task with id %s", id)
	}

	return nil
}
