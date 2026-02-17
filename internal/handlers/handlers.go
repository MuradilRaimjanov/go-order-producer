package handlers

import (
	"encoding/json"
	"go-order-producer/internal/database"
	"go-order-producer/internal/models"
	"net/http"
	"strconv"
	"strings"
)

type Handlers struct {
	store *database.TaskStore
}

func NewHandlers(store *database.TaskStore) *Handlers {
	return &Handlers{store: store}
}

func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}

func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	respondWithJSON(w, statusCode, map[string]string{"error": message})
}

func (h *Handlers) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.store.GetAll()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Ошибка получения задач")
		return
	}

	respondWithJSON(w, http.StatusOK, tasks)
}

func (h *Handlers) GetTaskById(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/tasks/"), "/")
	idStr := pathParts[0]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Некорректный ID задачи")
		return
	}

	task, err := h.store.GetByID(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, task)
}

func (h *Handlers) CreateTask(w http.ResponseWriter, r *http.Request) {
	var input models.CreateTaskInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if strings.TrimSpace(input.Title) == "" {
		respondWithError(w, http.StatusBadRequest, "Некоррекно оправлен заголовок задачи")
		return
	}

	task, err := h.store.Create(input)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, task)
}

func (h *Handlers) UpdateTask(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/tasks/"), "/")
	idStr := pathParts[0]

	id, err := strconv.Atoi(idStr)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	var input models.UpdateTaskInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, http.StatusBadRequest, "Некорректный данные")
		return
	}

	if input.Title != nil && strings.TrimSpace(*input.Title) == "" {
		respondWithError(w, http.StatusBadRequest, "Заголовок не корректен")
		return
	}

	task, err := h.store.Update(id, input)

	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		} else {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	respondWithJSON(w, http.StatusOK, task)
}

func (h *Handlers) DeleteTask(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/tasks/"), "/")
	idStr := pathParts[0]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Некоррекный id задачи")
		return
	}

	err = h.store.Delete(id)
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			respondWithError(w, http.StatusNotFound, err.Error())
			return
		} else {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}
