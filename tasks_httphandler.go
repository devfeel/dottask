package task

import (
	"encoding/json"
	"net/http"
)

// CounterOutputHttpHandler Http Handler for output counter info
func (service *TaskService) CounterOutputHttpHandler(w http.ResponseWriter, r *http.Request) {
	str, _ := json.Marshal(service.GetAllTaskCountInfo())
	w.Write([]byte(str))
}

// TaskOutputHttpHandler Http Handler for output task info
func (service *TaskService) TaskOutputHttpHandler(w http.ResponseWriter, r *http.Request) {
	str, _ := json.Marshal(service.GetAllTasks())
	w.Write([]byte(str))
}
