package transport

import (
	"encoding/json"
	"fmt"
	"net/http"
	"tasks/internal/database"
	"tasks/internal/entities"

	"github.com/rs/zerolog/log"
)

type Resourse struct {
	s *database.PostgresStorage
}

func NewResource(s *database.PostgresStorage) *Resourse {
	return &Resourse{s: s}
}

func (tr *Resourse) GetTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := tr.s.GetTasks()
	fmt.Print(tasks)

	if err != nil {
		log.Error().Err(err).Msg("Error getting tasks")
	}

	err = json.NewEncoder(w).Encode(tasks)

	if err != nil {
		fmt.Println("Failed to encode: ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (tr *Resourse) CreateTask(w http.ResponseWriter, r *http.Request) {
	var reqBody entities.Task

	err := json.NewDecoder(r.Body).Decode(&reqBody)

	if err != nil {
		log.Error().Err(err).Msg("Failed to decode")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = tr.s.InsertTask(reqBody)

	if err != nil {
		log.Error().Err(err).Msg("Failed to create task")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (tr *Resourse) UpdateTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var reqBody entities.Task

	err := json.NewDecoder(r.Body).Decode(&reqBody)

	if err != nil {
		fmt.Println("Failed to encode: ", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = tr.s.UpdateTask(id, reqBody)

	if err != nil {
		log.Error().Err(err).Msg("Failed to update task")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (tr *Resourse) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	err := tr.s.DeleteTask(id)

	if err != nil {
		log.Error().Err(err).Msg("Error deleting task")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (ur *Resourse) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user entities.User

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		log.Error().Err(err).Msg("Failed to decode")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = ur.s.InsertUser(user)

	if err != nil {
		log.Error().Err(err).Msg("Error creating user")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
