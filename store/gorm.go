package store

import (
	"github.com/sing3demons/todoapi/todo"
	"gorm.io/gorm"
)

type GormStore struct {
	db *gorm.DB
}

func NewGormStore(db *gorm.DB) Storer {
	return &GormStore{db: db}
}

func (s *GormStore) New(todo *todo.Todo) error {
	return s.db.Create(todo).Error
}

func (s *GormStore) List(todos *[]todo.Todo) error {
	return s.db.Find(todos).Error
}

func (s *GormStore) Delete(todo *todo.Todo, id int) error {
	return s.db.Delete(todo, id).Error
}
