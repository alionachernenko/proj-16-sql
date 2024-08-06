package database

import (
	"database/sql"
	"fmt"
	"tasks/internal/entities"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

type PostgresStorage struct {
    db *sql.DB
}

func NewPostgresStorage(connString string) (*PostgresStorage, error) {
	db, err := sql.Open("postgres", connString)

	if err != nil {
		return nil, fmt.Errorf("opening database: %v", err)
	}

	return &PostgresStorage{
		db: db,
	}, nil
}


type TasksResource struct {
	Storage *PostgresStorage
}


func (s *PostgresStorage) GetTasks() ([]entities.Task, error) {
	rows, err := s.db.Query("SELECT id, value, done, user_id FROM tasks")

	if err != nil {
		return nil, fmt.Errorf("querying tasks: %v", err)
	}

	defer rows.Close()

	var tasks []entities.Task

	for rows.Next() {
		var task entities.Task

		err := rows.Scan(&task.Id, &task.Value, &task.Done, &task.UserId)

		if err != nil {
			return nil, fmt.Errorf("scanning rows: %v", err)
		}

		tasks = append(tasks, task)
	}

	log.Info().Msgf("Got tasks: %v", tasks)

	return tasks, nil
}

func (s *PostgresStorage) InsertTask(task entities.Task) error {
	_, err := s.db.Exec("INSERT INTO tasks(user_id, value, done) VALUES($1, $2, $3)", task.UserId, task.Value, task.Done)

	if err != nil {
		return fmt.Errorf("inserting task: %v", err)
	}

	return nil
}

func (s *PostgresStorage) UpdateTask(id string, task entities.Task) error {
	_, err := s.db.Exec("UPDATE tasks SET value = $1, done = $2, user_id = $3 WHERE id = $4", task.Value, task.Done, task.UserId, id)

	if err != nil {
		return fmt.Errorf("updating task: %v", err)
	}

	return nil

}

func (s *PostgresStorage) DeleteTask(id string) error {
	_, err := s.db.Exec("DELETE FROM tasks WHERE id = $1", id)

	if err != nil {
		return fmt.Errorf("deleting task: %v", err)
	}

	return nil
}

type UsersResource struct {
	Storage *PostgresStorage
}


func (s *PostgresStorage) GetUser(username string) (entities.User, error) {
	rows, err := s.db.Query("SELECT id, username, password FROM users WHERE username = $1", username)

	if err != nil {
		return entities.User{}, fmt.Errorf("querying user: %v", err)
	}

	var user entities.User

	if rows.Next() {
		err := rows.Scan(&user.Id, &user.Username, &user.Password)
		if err != nil {
			return entities.User{}, fmt.Errorf("scanning rows: %v", err)
		}
	} else {
		return entities.User{}, fmt.Errorf("user not found")
	}

	return user, nil
}

func (s *PostgresStorage) InsertUser(u entities.User) error {
	_, err := s.db.Exec("INSERT INTO users(username, password) VALUES ($1, $2)", u.Username, u.Password)

	if err != nil {
		return fmt.Errorf("inserting user: %v", err)
	}

	return nil
}
