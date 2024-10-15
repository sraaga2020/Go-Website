package main

import (
	
    "html/template"
    "net/http"
)

func main() {
    http.HandleFunc("/", formHandler) // Handler for editing form
    http.HandleFunc("/create", submitHandler) // Handler for form submission
	// fmt.Fprintf(w, "Hello, you've requested: %s\n", r.URL.Path)
    http.ListenAndServe(":8080", nil)
}

func formHandler(w http.ResponseWriter, r *http.Request) {

    // Load the task creation form from the template file
    tmpl, err := template.ParseFiles("templates/eventForm.html")
    if err != nil {
        http.Error(w, "Could not load form", http.StatusInternalServerError)
        return
    }

    // Execute the template
    err = tmpl.Execute(w, nil)
    if err != nil {
        http.Error(w, "Could not execute template", http.StatusInternalServerError)
    }
}


func submitHandler(w http.ResponseWriter, r *http.Request) {
    
	if r.Method == http.MethodPost {
        r.ParseForm()
        taskName := r.FormValue("taskName")
        taskTime := r.FormValue("taskTime")
		taskDate := r.FormValue("taskDate")
        // Load the output template
        tmpl, err := template.ParseFiles("templates/output.html")
        

        // Create a data structure to pass to the template
        data := struct {
            TaskName string
            TaskTime string
			TaskDate string
        }{
            TaskName: taskName,
            TaskTime: taskTime,
			TaskDate: taskDate,
        }

        // Execute the template with the data
        err = tmpl.Execute(w, data)

        if err != nil {
            http.Error(w, "Could not execute template", http.StatusInternalServerError)
            return
        }
		
		
		
    } else {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
    }


}