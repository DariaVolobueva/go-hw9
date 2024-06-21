package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
)

type Student struct {
    ID     int
    Name   string
    Grades map[string]float64
}

type Class struct {
    Name    string
    Teacher string
    Students map[int]Student
}

var (
    classes = map[string]Class{
        "10-A": {
            Name:    "10-A",
            Teacher: "teacher123",
            Students: map[int]Student{
                1: {ID: 1, Name: "John Doe", Grades: map[string]float64{"Math": 90, "History": 85}},
                2: {ID: 2, Name: "Jane Smith", Grades: map[string]float64{"Math": 92, "History": 88}},
            },
        },
        "10-B": {
            Name:    "10-B",
            Teacher: "teacher456",
            Students: map[int]Student{
                1: {ID: 1, Name: "Alice Johnson", Grades: map[string]float64{"Math": 85, "History": 80}},
                2: {ID: 2, Name: "Bob Brown", Grades: map[string]float64{"Math": 87, "History": 82}},
            },
        },
    }
    mu sync.Mutex
)

func isAuthenticated(r *http.Request, className string) bool {
    teacherID := r.URL.Query().Get("teacher")
    class, ok := classes[className]
    return ok && teacherID == class.Teacher
}

func classHandler(w http.ResponseWriter, r *http.Request) {
    className := r.URL.Query().Get("class")
    mu.Lock()
    defer mu.Unlock()
    
    if !isAuthenticated(r, className) {
        http.Error(w, "Forbidden", http.StatusForbidden)
        return
    }
    
    class, ok := classes[className]
    if !ok {
        http.Error(w, "Class not found", http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode(class)
}

func studentHandler(w http.ResponseWriter, r *http.Request) {
    mu.Lock()
    defer mu.Unlock()

    className := r.URL.Query().Get("class")
    
    class, ok := classes[className]
    if !ok {
        http.Error(w, "Class not found", http.StatusNotFound)
        return
    }

    idStr := r.URL.Path[len("/student/"):]
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid student ID", http.StatusBadRequest)
        return
    }
    
    student, ok := class.Students[id]
    if !ok {
        http.Error(w, "Student not found", http.StatusNotFound)
        return
    }

	if !isAuthenticated(r, className) {
        http.Error(w, "Forbidden", http.StatusForbidden)
        return
    }
    
    json.NewEncoder(w).Encode(student)
}

func main() {
    http.HandleFunc("/class", classHandler)
    http.HandleFunc("/student/", studentHandler)
    
    fmt.Println("Server is listening on port 8080...")
    http.ListenAndServe(":8080", nil)
}
