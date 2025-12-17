
# Problem 2: Find and Fix Issues (10-15 minutes)
## Description
A junior developer wrote this simple todo API in Go, but it has some problems. Find at least 3
issues and explain how to fix them.
### Buggy Go Code:

```go
package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Todo struct {
	ID        int    `json:"id"`
	Text      string `json:"text"`
	Completed bool   `json:"completed"`
}

var todos []Todo
var nextID = 1

func getTodos(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(todos)
}

func createTodo(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	json.NewDecoder(r.Body).Decode(&todo)
	todo.ID = nextID
	nextID++
	todo.Completed = false
	todos = append(todos, todo)
	json.NewEncoder(w).Encode(todo)
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])
	for i, todo := range todos {
		if todo.ID == id {
			todos = append(todos[:i], todos[i+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Todo deleted"})
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/todos", getTodos).Methods("GET")
	r.HandleFunc("/todos", createTodo).Methods("POST")
	r.HandleFunc("/todos/{id}", deleteTodo).Methods("DELETE")
	http.ListenAndServe(":8080", r)
}
```

### What to do:
1. List the Suggestion/feedback you found
2. Explain why each is a problem
### Example format:
Issue 1: Missing Content-Type header
Problem: The API doesn't set proper JSON content type
#### Some hints of what to look for:
● Are proper HTTP headers set?
● Is there any error handling for JSON decoding?
● What happens if the ID conversion fails?
● Are there CORS headers for frontend communication?
● Is there proper validation for incoming data?

## Answer
- Issue 1: Missing Content-Type header 
Problem: The API doesn't set proper JSON content type

- Issue 2: Missing error handling for JSON decoding
Problem: The API doesn't handle JSON decoding errors

- Issue 3: Missing error handling for ID conversion
Problem: The API doesn't handle ID conversion errors

- Issue 4: Missing CORS headers
Problem: The API doesn't handle CORS headers

- Issue 5: Missing validation for incoming data
Problem: The API doesn't handle validation for incoming data

- Issue 6: Missing mutex for concurrent access
Problem: The API doesn't handle concurrent access
