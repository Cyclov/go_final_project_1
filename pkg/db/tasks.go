package db

func Tasks(limit int) ([]*Task, error) {

	query := `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?`
	//now := time.Now().Format("20060102")
	resp, err := db.Query(query, limit)

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
		return tasks, err
	}

	return tasks, nil
}
