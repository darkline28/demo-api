package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type Task struct {
	ID     int    `json:"id"`
	Text   string `json:"text"`
	Status string `json:"status"`
}

type Server struct {
	Tasks  map[int]Task
	NextID int
}

type ServerReq struct {
	Task string `json:"task"`
}

func (s *Server) handlerTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		// вернуть все задачи (массив, не map)
		tasks := make([]Task, 0, len(s.Tasks))
		for _, t := range s.Tasks {
			tasks = append(tasks, t)
		}
		json.NewEncoder(w).Encode(tasks)

	case http.MethodPost:
		// создать новую задачу
		var req struct {
			Text string `json:"text"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"bad request"}`, http.StatusBadRequest)
			return
		}
		if req.Text == "" {
			http.Error(w, `{"error":"text required"}`, http.StatusBadRequest)
			return
		}

		task := Task{
			ID:     s.NextID,
			Text:   req.Text,
			Status: "new",
		}
		s.Tasks[s.NextID] = task
		s.NextID++

		w.Header().Set("Location", "/tasks/"+strconv.Itoa(task.ID))
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(task)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
func (s *Server) handlerTaskByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := r.URL.Path[len("/tasks/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
		return
	}

	task, ok := s.Tasks[id]
	if !ok {
		http.Error(w, `{"error":"task not found"}`, http.StatusNotFound)
		return
	}

	switch r.Method {

	case http.MethodGet:
		json.NewEncoder(w).Encode(task)

	case http.MethodPatch:
		var req Task
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"bad request"}`, http.StatusBadRequest)
			return
		}
		if req.Text != "" {
			task.Text = req.Text
		}
		if req.Status != "" {
			task.Status = req.Status
		}
		s.Tasks[id] = task
		json.NewEncoder(w).Encode(task)

	case http.MethodDelete:
		delete(s.Tasks, id)
		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	srv := &Server{
		Tasks:  make(map[int]Task),
		NextID: 1,
	}
	http.HandleFunc("/tasks", srv.handlerTasks)
	http.HandleFunc("/tasks/", srv.handlerTaskByID)
	http.ListenAndServe(":8080", nil)
}
