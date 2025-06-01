package todo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockStore struct {
	todos []Todo
	err   error
}

func NewMockStore(todos []Todo, err error) storer {
	return &MockStore{
		todos: todos,
		err:   err,
	}
}

func (m *MockStore) New(todo *Todo) error {
	if m.err != nil {
		return m.err
	}
	m.todos = append(m.todos, *todo)
	return nil
}
func (m *MockStore) List(todos *[]Todo) error {
	if m.err != nil {
		return m.err
	}
	*todos = m.todos
	return nil
}
func (m *MockStore) Delete(todo *Todo, id int) error {
	if m.err != nil {
		return m.err
	}

	return nil // or return an error if not found
}

// mock context to simulate the behavior of a web context

type MockContext struct {
	BindCalled bool
	BindInput  interface{}
	BindErr    error

	JSONCalled bool
	JSONCode   int
	JSONValue  interface{}

	ParamCalled    bool
	ParamKey       string
	ParamReturnVal string

	StatusCalled bool
	StatusCode   int
}

func (m *MockContext) Bind(v interface{}) error {
	m.BindCalled = true
	m.BindInput = v
	if m.BindErr != nil {
		return m.BindErr
	}
	return nil
}

func (m *MockContext) JSON(code int, v interface{}) {
	m.JSONCalled = true
	m.JSONCode = code
	m.StatusCode = code
	m.JSONValue = v
}

func (m *MockContext) Param(key string) string {
	m.ParamCalled = true
	m.ParamKey = key
	return m.ParamReturnVal
}

func (m *MockContext) Status(code int) {
	// Mock implementation, no action needed
	m.StatusCalled = true
	m.StatusCode = code
}

func (m *MockContext) Clear() {
	m.BindCalled = false
	m.BindInput = nil
	m.BindErr = nil

	m.JSONCalled = false
	m.JSONCode = 0
	m.JSONValue = nil

	m.ParamCalled = false
	m.ParamKey = ""
	m.ParamReturnVal = ""

	m.StatusCalled = false
	m.StatusCode = 0
}

func TestNewTodo(t *testing.T) {
	todos := []Todo{
		{ID: 1, Title: "test todo"},
	}
	t.Run("should create a new todo", func(t *testing.T) {
		db := NewMockStore(todos, nil)
		handler := NewTodoHandler(db)

		mockContext := &MockContext{}
		handler.NewTask(mockContext)

		// Bind
		assert.Equal(t, mockContext.BindCalled, true)
		assert.NotNil(t, mockContext.BindInput)
		assert.NoError(t, mockContext.BindErr)

		// JSON
		assert.Equal(t, mockContext.StatusCode, 201)
		assert.Equal(t, mockContext.JSONCalled, true)
		assert.NotNil(t, mockContext.JSONValue)

		mockContext.Clear()
	})

	t.Run("should return error on binding", func(t *testing.T) {
		db := NewMockStore(todos, nil)
		handler := NewTodoHandler(db)

		mockContext := &MockContext{
			BindErr: assert.AnError,
		}
		handler.NewTask(mockContext)

		// Bind
		assert.Equal(t, mockContext.BindCalled, true)
		assert.NotNil(t, mockContext.BindInput)
		assert.Equal(t, mockContext.BindErr, assert.AnError)

		// JSON
		assert.Equal(t, mockContext.StatusCode, 400)
		assert.Equal(t, mockContext.JSONCalled, true)
		assert.NotNil(t, mockContext.JSONValue)

		mockContext.Clear()
	})

	t.Run("should return error on store", func(t *testing.T) {
		db := NewMockStore(todos, assert.AnError)
		handler := NewTodoHandler(db)

		mockContext := &MockContext{}
		handler.NewTask(mockContext)

		// Bind
		assert.Equal(t, mockContext.BindCalled, true)
		assert.NotNil(t, mockContext.BindInput)
		assert.NoError(t, mockContext.BindErr)

		// JSON
		assert.Equal(t, mockContext.StatusCode, 400)
		assert.Equal(t, mockContext.JSONCalled, true)
		assert.NotNil(t, mockContext.JSONValue)
		mockContext.Clear()
	})

}

func TestListTodos(t *testing.T) {
	todos := []Todo{
		{ID: 1, Title: "Test Todo"},
		{ID: 2, Title: "Another Todo"},
	}

	t.Run("should list all todos", func(t *testing.T) {
		db := NewMockStore(todos, nil)
		handler := NewTodoHandler(db)

		mockContext := &MockContext{}
		handler.List(mockContext)

		assert.Equal(t, mockContext.JSONCalled, true)
		assert.Equal(t, mockContext.StatusCode, 200)
		assert.Len(t, mockContext.JSONValue.([]Todo), 2)

		mockContext.Clear()
	})

	t.Run("should return error on store error", func(t *testing.T) {
		db := NewMockStore(nil, assert.AnError)
		handler := NewTodoHandler(db)

		mockContext := &MockContext{}
		handler.List(mockContext)

		assert.Equal(t, mockContext.JSONCalled, true)
		assert.Equal(t, mockContext.StatusCode, 500)
		assert.NotNil(t, mockContext.JSONValue)

		mockContext.Clear()
	})
}

func TestRemoveTodo(t *testing.T) {
	todos := []Todo{
		{ID: 1, Title: "Test Todo"},
	}

	t.Run("should remove a todo", func(t *testing.T) {
		db := NewMockStore(todos, nil)
		handler := NewTodoHandler(db)

		mockContext := &MockContext{}
		mockContext.ParamReturnVal = "1" // Simulate getting ID from URL
		handler.Remove(mockContext)

		assert.Equal(t, mockContext.StatusCode, 200)
		assert.Equal(t, mockContext.JSONCalled, true)

		mockContext.Clear()
	})

	t.Run("should return error on store error", func(t *testing.T) {
		db := NewMockStore(nil, assert.AnError)
		handler := NewTodoHandler(db)

		mockContext := &MockContext{}
		mockContext.ParamReturnVal = "1" // Simulate getting ID from URL
		handler.Remove(mockContext)

		assert.Equal(t, mockContext.StatusCode, 500)
		assert.Equal(t, mockContext.JSONCalled, true)
		assert.NotNil(t, mockContext.JSONValue)

		mockContext.Clear()
	})

	t.Run("should return error on invalid ID", func(t *testing.T) {
		db := NewMockStore(todos, nil)
		handler := NewTodoHandler(db)

		mockContext := &MockContext{}
		mockContext.ParamReturnVal = "invalid" // Simulate invalid ID
		handler.Remove(mockContext)

		assert.Equal(t, mockContext.StatusCode, 400)
		assert.Equal(t, mockContext.JSONCalled, true)
		assert.NotNil(t, mockContext.JSONValue)

		mockContext.Clear()
	})
}
