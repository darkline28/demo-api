package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Server struct {
	Task string
}

type ServerReq struct {
	Task string `json:"task"`
}

func (s *Server) handlerPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ServerReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	s.Task = req.Task
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handlerGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if s.Task != "" {
		fmt.Fprintf(w, "Hello, %s", s.Task)
	} else {
		fmt.Fprintln(w, "Hello World!")
	}
}

func main() {
	srv := &Server{}
	http.HandleFunc("/", srv.handlerGet)
	http.HandleFunc("/task", srv.handlerPost)
	http.ListenAndServe(":8080", nil)
}
