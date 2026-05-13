package main

import (
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{db: db}
}

func (s *UserStore) migrate() error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id         INTEGER PRIMARY KEY AUTOINCREMENT,
		username   TEXT    NOT NULL UNIQUE,
		password   TEXT    NOT NULL
	);`
	_, err := s.db.Exec(query)
	return err
}

func (s *UserStore) Create(username, password string) (User, error) {
	// hash the password — never store plain text
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}

	query := `
	INSERT INTO users (username, password)
	VALUES (?, ?)
	RETURNING id, username`

	var user User
	err = s.db.QueryRow(query, username, string(hashed)).Scan(
		&user.ID, &user.Username,
	)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (s *UserStore) GetByUsername(username string) (User, bool, error) {
	query := `SELECT id, username, password FROM users WHERE username = ?`
	var user User
	err := s.db.QueryRow(query, username).Scan(
		&user.ID, &user.Username, &user.Password,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, false, nil
		}
		return User{}, false, err
	}
	return user, true, nil
}

func (s *UserStore) CheckPassword(user User, password string) error {
	_, ok, err := s.GetByUsername(user.Username)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return errors.New("invalid password")
		}
		return err
	}
	return nil
}
