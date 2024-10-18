package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
)

type Task struct {
	TaskName string `json:"taskName"`
	TaskTime string `json:"taskTime"`
	TaskDate string `json:"taskDate"`
}

const dataFile = "tasks.json"

func main() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}
	fmt.Println("Current directory:", dir)

	http.HandleFunc("/", formHandler)
	http.HandleFunc("/create", submitHandler)
	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Server error:", err)
	}
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/eventForm.html")
	if err != nil {
		http.Error(w, "Could not load form", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, "Could not execute template", http.StatusInternalServerError)
	}
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Could not parse form", http.StatusBadRequest)
		return
	}

	task := Task{
		TaskName: r.FormValue("taskName"),
		TaskTime: r.FormValue("taskTime"),
		TaskDate: r.FormValue("taskDate"),
	}

	if err := saveTask(task); err != nil {
		http.Error(w, "Could not save task", http.StatusInternalServerError)
		return
	}

	// Read all tasks from the JSON file
	tasks, err := readAllTasks()
	if err != nil {
		http.Error(w, "Could not read tasks", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("templates/output.html")
	if err != nil {
		http.Error(w, "Could not load output template", http.StatusInternalServerError)
		return
	}

	data := struct {
		NewTask  Task
		AllTasks []Task
	}{
		NewTask:  task,
		AllTasks: tasks,
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Could not execute template", http.StatusInternalServerError)
	}
}

func readAllTasks() ([]Task, error) {
	var tasks []Task
	data, err := ioutil.ReadFile(dataFile)
	if err != nil {
		if os.IsNotExist(err) {
			return tasks, nil
		}
		return nil, err
	}

	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

func saveTask(task Task) error {
	var tasks []Task

	data, err := ioutil.ReadFile(dataFile)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("error reading file: %v", err)
	}

	if len(data) > 0 {
		if err := json.Unmarshal(data, &tasks); err != nil {
			return fmt.Errorf("error decoding JSON: %v", err)
		}
	}

	tasks = append(tasks, task)

	updatedData, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("error encoding JSON: %v", err)
	}

	if err := ioutil.WriteFile(dataFile, updatedData, 0644); err != nil {
		return fmt.Errorf("error writing file: %v", err)
	}

	fmt.Println("Saved task:", task)
	return nil
}
