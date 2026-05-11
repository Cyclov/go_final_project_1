package db

import (
	"database/sql"
	"time"
)

func Tasks(limit int, search string) ([]*Task, error) {

	var resp *sql.Rows
	var err error
	var query string

	if search != "" {
		if parsed, err := time.Parse("02.01.2006", search); err == nil {
			searchDate := parsed.Format(DateFormat)
			query = `SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? ORDER BY date LIMIT ?`
			resp, err = db.Query(query, searchDate, limit)
		} else {
			query = `SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?`
			resp, err = db.Query(query, "%"+search+"%", "%"+search+"%", limit)
		}

	} else {
		query = `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?`
		resp, err = db.Query(query, limit)
	}

	var tasks []*Task
	tasks = make([]*Task, 0)

	if err != nil {
		return tasks, err
	}
	defer resp.Close()

	for resp.Next() {
		t := &Task{}
		if err := resp.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return tasks, err
		}
		tasks = append(tasks, t)
	}

	if err := resp.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}
