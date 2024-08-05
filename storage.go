package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(connString string) (*Storage, error) {
	db, err := sql.Open("postgres", connString)

	if err != nil {
		return nil, fmt.Errorf("opening database: %v", err)
	}

	return &Storage{
		db: db,
	}, nil
}

func (s *Storage) GetTasks() ([]Task, error) {
	rows, err := s.db.Query("SELECT id, value, done, user_id FROM tasks")

	if err != nil {
		return nil, fmt.Errorf("querying tasks: %v", err)
	}

	defer rows.Close()

	var tasks []Task

	for rows.Next() {
		var task Task

		err := rows.Scan(&task.Id, &task.Value, &task.Done, &task.UserId)

		if err != nil {
			return nil, fmt.Errorf("scanning rows: %v", err)
		}

		tasks = append(tasks, task)
	}

	log.Info().Msgf("Got tasks: %v", tasks)

	return tasks, nil
}

func (s *Storage) InsertTask(task Task) error {
	_, err := s.db.Exec("INSERT INTO tasks(user_id, value, done) VALUES($1, $2, $3)", task.UserId, task.Value, task.Done)

	if err != nil {
		return fmt.Errorf("inserting task: %v", err)
	}

	return nil
}

func (s *Storage) UpdateTask(id string, task Task) error {
	_, err := s.db.Exec("UPDATE tasks SET value = $1, done = $2, user_id = $3 WHERE id = $4", task.Value, task.Done, task.UserId, id)

	if err != nil {
		return fmt.Errorf("updating task: %v", err)
	}

	return nil

}

func (s *Storage) DeleteTask(id string) error {
	_, err := s.db.Exec("DELETE FROM tasks WHERE id = $1", id)

	if err != nil {
		return fmt.Errorf("deleting task: %v", err)
	}

	return nil
}

func (s *Storage) GetUser(username string) (User, error) {
	rows, err := s.db.Query("SELECT id, username, password FROM users WHERE username = $1", username)

	if err != nil {
		return User{}, fmt.Errorf("querying user: %v", err)
	}

	var user User

	if rows.Next() {
		err := rows.Scan(&user.Id, &user.Username, &user.Password)
		if err != nil {
			return User{}, fmt.Errorf("scanning rows: %v", err)
		}
	} else {
		return User{}, fmt.Errorf("user not found")
	}

	return user, nil
}

func (s *Storage) InsertUser(u User) error {
	_, err := s.db.Exec("INSERT INTO users(username, password) VALUES ($1, $2)", u.Username, u.Password)

	if err != nil {
		return fmt.Errorf("inserting user: %v", err)
	}

	return nil
}
