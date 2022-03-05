package todo

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"
)

type DeletedAt sql.NullTime

type Todo struct {
	Title     string `json:"text" binding:"required"`
	ID        uint   `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt DeletedAt `gorm:"index"`
}

func (Todo) TableName() string {
	return "todos"
}

type storer interface {
	New(*Todo) error
	List(*[]Todo) error
	Delete(*Todo, int) error
}

type TodoHandler struct {
	store storer
}

func NewTodoHandler(store storer) *TodoHandler {
	return &TodoHandler{store: store}
}

type Context interface {
	Bind(interface{}) error
	JSON(int, interface{})
	Param(string) string
}

func (t *TodoHandler) NewTask(c Context) {
	var todo Todo
	if err := c.Bind(&todo); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	if err := t.store.New(&todo); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, map[string]interface{}{
		"ID": todo.ID,
	})
}

func (t *TodoHandler) List(c Context) {
	var todos []Todo
	if err := t.store.List(&todos); err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, todos)
}

func (t *TodoHandler) Remove(c Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	if err := t.store.Delete(&Todo{}, id); err != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, map[string]interface{}{
		"status": "success",
	})

}
