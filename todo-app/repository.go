package main

import (
	"errors"
	"sync"
)

// Todo Database/Repository Mock

type TodoDB interface {
	GetCategoryByName(name string) (Category, error)
	GetUserById(id int) (User, error)
	GetUsers() []User
	GetTodos() []Todo
	CreateTodo(todo Todo) Todo
	DeleteTodo(id int) bool
}

type todoDB struct {
	todos      []Todo
	categories []Category
	users      []User
	nextTodoID int
	mu         sync.Mutex
}

func NewTodoDB() TodoDB {
	return &todoDB{
		todos:      []Todo{},
		nextTodoID: 0,
		mu:         sync.Mutex{},
		categories: categories,
		users: []User{
			{ID: 1, Name: "John Doe", Email: "johndoe@mail.com", Password: "1234"},
			{ID: 2, Name: "Jane Doe", Email: "janedoe@mail.com", Password: "1234"},
		},
	}
}

func (db *todoDB) GetCategoryByName(name string) (Category, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	for _, category := range db.categories {
		if category.Name == name {
			return category, nil
		}
	}
	return Category{}, errors.New("category not found")
}

func (db *todoDB) nextTodoIDInc() int {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.nextTodoID++
	return db.nextTodoID
}

func (db *todoDB) GetUsers() []User {
	db.mu.Lock()
	defer db.mu.Unlock()
	return db.users
}

func (db *todoDB) GetUserById(id int) (User, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	for _, user := range db.users {
		if user.ID == id {
			return user, nil
		}
	}
	return User{}, errors.New("User Not Found")
}

func (db *todoDB) GetTodos() []Todo {
	db.mu.Lock()
	defer db.mu.Unlock()
	return db.todos
}

func (db *todoDB) CreateTodo(todo Todo) Todo {
	todo.ID = db.nextTodoIDInc()
	db.mu.Lock()
	defer db.mu.Unlock()
	db.todos = append(db.todos, todo)
	return todo
}

func (db *todoDB) DeleteTodo(id int) bool {
	db.mu.Lock()
	defer db.mu.Unlock()
	for i, todo := range db.todos {
		if todo.ID == id {
			db.todos = append(db.todos[:i], db.todos[i+1:]...)
			return true
		}
	}
	return false
}
