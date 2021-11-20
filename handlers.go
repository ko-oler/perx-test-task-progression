package main

import (
	"encoding/json"
	"net/http"
)

//Handlers
func GetListHandler(tl *TaskList) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tl.SortTasks()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tl.tasks)
		w.WriteHeader(http.StatusOK)
	}
}

func AddHandler(tl *TaskList) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var task Task
		w.Header().Set("Content-Type", "application/json")
		err := json.NewDecoder(r.Body).Decode(&task)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		tl.AddTaskToQueue(&task)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(task)
	}
}
