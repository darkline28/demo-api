package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	taskservices "study/api/internal/TaskServices"
)

type TaskHandlers struct {
	service taskservices.TaskService
}

func NewTaskHandlers(s taskservices.TaskService) *TaskHandlers {
	return &TaskHandlers{
		service: s,
	}
}

func (h *TaskHandlers) HandleTasks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		tasks, err := h.service.List()
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "could not get tasks")
			return
		}
		writeJSON(w, http.StatusOK, tasks)
	case http.MethodPost:
		var req taskservices.TaskRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSONError(w, http.StatusBadRequest, "bad request")
			return
		}
		defer r.Body.Close()
		task, err := h.service.Create(req.Task)
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, "could not create task")
			return
		}
		w.Header().Set("Location", "/tasks/"+strconv.Itoa(task.ID))
		writeJSON(w, http.StatusOK, task)
	default:
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
}
func (h *TaskHandlers) HandleTaskByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/tasks/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid id")
		return
	}
	switch r.Method {
	case http.MethodPatch:
		var req taskservices.Task
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSONError(w, http.StatusBadRequest, "bad request")
			return
		}
		defer r.Body.Close()

		updateTask, err := h.service.Update(id, req.Text, req.Status)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "could not update task")
			return
		}
		writeJSON(w, http.StatusOK, updateTask)

	case http.MethodDelete:
		if err := h.service.Delete(id); err != nil {
			writeJSONError(w, http.StatusInternalServerError, "could not delete task")
			return
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}

}

func writeJSONError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{
		"error":   http.StatusText(code),
		"message": msg,
	})
}
func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}
