package store

import "github.com/sing3demons/todoapi/todo"

type Storer interface {
	New(todo *todo.Todo) error
	List(todos *[]todo.Todo) error
	Delete(todo *todo.Todo, id int) error
}
