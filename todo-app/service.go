package main

import (
	"errors"
)

// Todo Service
type TodoService struct {
	db       TodoDB
	aiClient *AiService
}

func NewTodoService(db TodoDB, aiClient *AiService) *TodoService {
	return &TodoService{
		db:       db,
		aiClient: aiClient,
	}
}

func (s *TodoService) GetUsers() ([]User, error) {
	return s.db.GetUsers(), nil
}

func (s *TodoService) CreateTodo(todo Todo) (Todo, error) {
	category, err := s.aiClient.PredictCategory(todo)
	if err != nil {
		return todo, err
	}
	cat, err := s.db.GetCategoryByName(category.Name)
	if err != nil {
		return todo, err
	}
	user, err := s.db.GetUserById(todo.UserID)
	if err != nil {
		return todo, err
	}
	todo.CategoryID = cat.ID
	todo.Category = &cat
	todo.User = &user
	createdTodo := s.db.CreateTodo(todo)
	return createdTodo, nil
}

func (s *TodoService) GetTodos() ([]Todo, error) {
	return s.db.GetTodos(), nil
}

func (s *TodoService) DeleteTodo(id int) error {
	if s.db.DeleteTodo(id) {
		return nil
	}
	return errors.New("todo not found")
}
