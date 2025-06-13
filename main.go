package main

import (
	"html/template"
	"net/http"
	"sync"
)

type Task struct {
	ID   int
	Text string
}

type App struct {
	tasks  []Task
	nextID int
	mu     sync.Mutex
	tmpl   *template.Template
}

func NewApp() *App {
	return &App{
		tasks:  []Task{},
		nextID: 1,
		tmpl:   template.Must(template.ParseFiles("templates/index.html")),
	}
}

func (a *App) addTaskHandler(w http.ResponseWriter, r *http.Request) {

}

func (a *App) deleteTaskHandler(w http.ResponseWriter, r *http.Request) {

}

func (a *App) indexHandler(w http.ResponseWriter, r *http.Request) {

}

func main() {
	app := NewApp()

	http.HandleFunc("/", app.indexHandler)
	http.HandleFunc("/add", app.addTaskHandler)
	http.HandleFunc("/delete", app.deleteTaskHandler)

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	println("Server is running on :8080")
	http.ListenAndServe(":8080", nil)
}
