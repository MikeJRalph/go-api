package main

import (
	"html/template"
	"net/http"
	"strconv"
	"sync"
)

// Task represents a single task in the todo list
type Task struct {
	ID   int
	Text string
}

// App holds the application state
type App struct {
	tasks  []Task
	nextID int
	mu     sync.Mutex
	tmpl   *template.Template
}

// NewApp initializes a new App instance
func NewApp() *App {
	return &App{
		tasks:  []Task{},
		nextID: 1,
		tmpl:   template.Must(template.ParseFiles("templates/index.html")),
	}
}

// Adds new tasks to the list
func (a *App) addTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusBadRequest)
		return
	}

	taskText := r.FormValue("task")
	if taskText == "" {
		http.Error(w, "Task cannot be empty", http.StatusBadRequest)
		return
	}

	a.mu.Lock()
	a.tasks = append(a.tasks, Task{ID: a.nextID, Text: taskText})
	a.nextID++
	a.mu.Unlock()

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Deletes tasks from the list
func (a *App) deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not found", http.StatusBadRequest)
		return
	}

	idStr := r.FormValue("id")
	if idStr == "" {
		http.Error(w, "Missing Task ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Task ID", http.StatusBadRequest)
		return
	}

	a.mu.Lock()
	for i, task := range a.tasks {
		if task.ID == id {
			a.tasks = append(a.tasks[:i], a.tasks[i+1:]...)
			break
		}
	}
	a.mu.Unlock()

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Reders the index page with the current list of tasks
func (a *App) indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not found", http.StatusBadRequest)
		return
	}

	a.mu.Lock()
	tasks := a.tasks
	a.mu.Unlock()

	err := a.tmpl.ExecuteTemplate(w, "index.html", struct {
		Tasks []Task
	}{Tasks: tasks})
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func main() {
	app := NewApp()

	// Set up routes
	http.HandleFunc("/", app.indexHandler)
	http.HandleFunc("/add", app.addTaskHandler)
	http.HandleFunc("/delete", app.deleteTaskHandler)

	// Serve static files (e.g., CSS, JS)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Start the server
	println("Server is running on :8080")
	http.ListenAndServe(":8080", nil)
}
