package main

import (
	"encoding/json"
	"net/http"
)

//Handlers
func (tl *TaskList) GetListHandler(w http.ResponseWriter, r *http.Request) {
	tl.SortTasks()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tl.tasks)
	w.WriteHeader(http.StatusOK)
}

func (tl *TaskList) AddHandler(w http.ResponseWriter, r *http.Request) {
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
