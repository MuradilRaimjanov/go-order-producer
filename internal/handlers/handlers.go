package handlers

import (
	"go-order-producer/internal/database"
	"go-order-producer/internal/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type Handlers struct {
	store *database.TaskStore
}

func NewHandlers(store *database.TaskStore) *Handlers {
	return &Handlers{store: store}
}

func (h *Handlers) GetAllTasks(c echo.Context) error {
	tasks, err := h.store.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Ошибка получения задач",
		})
	}

	return c.JSON(http.StatusOK, tasks)
}

func (h *Handlers) GetTaskById(c echo.Context) error {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Некорректный ID задачи",
		})
	}

	task, err := h.store.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, task)
}

func (h *Handlers) CreateTask(c echo.Context) error {
	var input models.CreateTaskInput

	// Автоматический bind JSON → struct
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	if strings.TrimSpace(input.Title) == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Некорректно отправлен заголовок задачи",
		})
	}

	task, err := h.store.Create(input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, task)
}

func (h *Handlers) UpdateTask(c echo.Context) error {
	// Получаем ID из пути
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Некорректный ID задачи",
		})
	}

	// Читаем JSON из body
	var input models.UpdateTaskInput
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Некорректные данные",
		})
	}

	// Проверка заголовка задачи
	if input.Title != nil && strings.TrimSpace(*input.Title) == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Заголовок не корректен",
		})
	}

	// Обновление задачи
	task, err := h.store.Update(id, input)
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, task)
}

func (h *Handlers) DeleteTask(c echo.Context) error {
	// Получаем ID из пути
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Некорректный ID задачи",
		})
	}

	// Удаляем задачу
	err = h.store.Delete(id)
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{"result": "success"})
}
