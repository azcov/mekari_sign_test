package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Todo struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	Text        string    `json:"text"`
	Description string    `json:"description"`
	CategoryID  int       `json:"category_id"`
	Completed   bool      `json:"completed"`
	Category    *Category `json:"category,omitempty"`
	User        *User     `json:"user,omitempty"`
}

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var (
	categories = []Category{
		{ID: 1, Name: "Others"},
		{ID: 2, Name: "Work"},
		{ID: 3, Name: "Personal"},
		{ID: 4, Name: "Shopping"},
		{ID: 5, Name: "Health"},
		{ID: 6, Name: "Finance"},
		{ID: 7, Name: "Education"},
		{ID: 8, Name: "Travel"},
		{ID: 9, Name: "Fitness"},
		{ID: 10, Name: "Hobby"},
	}
)

// Router
func main() {
	h := NewHandler(NewTodoService(NewTodoDB(), NewAiServiceClient()))
	r := mux.NewRouter()
	r.Use(corsMiddleware)
	r.Use(loggerMiddleware)
	r.HandleFunc("/users", h.getUsers).Methods("GET")
	r.HandleFunc("/todos", h.getTodos).Methods("GET")
	r.HandleFunc("/todos", h.createTodo).Methods("POST")
	r.HandleFunc("/todos/{id}", h.deleteTodo).Methods("DELETE")
	r.HandleFunc("/", homeHandler)

	port := ":8080"
	log.Printf("Server started on %s\n", port)
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}

// Handler
type Handler struct {
	todoService *TodoService
}

func NewHandler(todoService *TodoService) *Handler {
	return &Handler{
		todoService: todoService,
	}
}

func (h *Handler) getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	users, err := h.todoService.GetUsers()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to get users"))
		return
	}

	if err := json.NewEncoder(w).Encode(users); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to encode response"))
		return
	}
}

func (h *Handler) getTodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	todos, err := h.todoService.GetTodos()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to get todos"))
		return
	}

	if err := json.NewEncoder(w).Encode(todos); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to encode response"))
		return
	}
}

func (h *Handler) createTodo(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request body"))
		return
	}
	todo, err := h.todoService.CreateTodo(todo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to create todo"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(todo); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to encode response"))
		return
	}
}

func (h *Handler) deleteTodo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	paramId, exists := params["id"]
	if !exists {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing id parameter"))
		return
	}

	id, err := strconv.Atoi(paramId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid id parameter"))
		return
	}

	if err := h.todoService.DeleteTodo(id); err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Todo not found"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Todo deleted"})
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	// Serve a specific HTML file for the homepage
	http.ServeFile(w, r, "index.html")
}

// Middleware

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s\n", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
