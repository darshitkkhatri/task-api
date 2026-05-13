package main

import (
	"database/sql"
	"time"
)

type SQLiteStore struct {
	db *sql.DB
}

// new — accepts *sql.DB directly
func NewSQLiteStore(db *sql.DB) (*SQLiteStore, error) {
	store := &SQLiteStore{db: db}
	if err := store.migrate(); err != nil {
		return nil, err
	}
	return store, nil
}
func (s *SQLiteStore) migrate() error {
	query := `
	CREATE TABLE IF NOT EXISTS tasks (
		id         INTEGER PRIMARY KEY AUTOINCREMENT,
		title      TEXT    NOT NULL,
		done       BOOLEAN NOT NULL DEFAULT FALSE,
		created_at DATETIME NOT NULL
	);`
	_, err := s.db.Exec(query)
	return err
}

var _ Storage = (*SQLiteStore)(nil)

func (s *SQLiteStore) Add(title string) (Task, error) {
	query := `
	INSERT INTO tasks (title, done, created_at)
	VALUES (?, ?, ?)
	RETURNING id, title, done, created_at`

	var task Task
	err := s.db.QueryRow(query, title, false, time.Now()).Scan(
		&task.ID, &task.Title, &task.Done, &task.CreatedAt,
	)
	if err != nil {
		return Task{}, err
	}
	return task, nil
}

func (s *SQLiteStore) GetAll() ([]Task, error) {
	query := `SELECT id, title, done, created_at FROM tasks ORDER BY id`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Done, &task.CreatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (s *SQLiteStore) GetByID(id int) (Task, bool, error) {
	query := `SELECT id, title, done, created_at FROM tasks WHERE id = ?`

	var task Task
	err := s.db.QueryRow(query, id).Scan(
		&task.ID, &task.Title, &task.Done, &task.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return Task{}, false, nil
	}
	if err != nil {
		return Task{}, false, err
	}
	return task, true, nil
}

func (s *SQLiteStore) Update(id int, title string, done bool) (Task, bool, error) {
	query := `
	UPDATE tasks SET title = ?, done = ?
	WHERE id = ?
	RETURNING id, title, done, created_at`

	var task Task
	err := s.db.QueryRow(query, title, done, id).Scan(
		&task.ID, &task.Title, &task.Done, &task.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return Task{}, false, nil
	}
	if err != nil {
		return Task{}, false, err
	}
	return task, true, nil
}

func (s *SQLiteStore) Delete(id int) (bool, error) {
	query := `DELETE FROM tasks WHERE id = ?`

	result, err := s.db.Exec(query, id)
	if err != nil {
		return false, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return rows > 0, nil
}
