// Package main wires up the application and starts the HTTP server
package main

import (
	"fmt"
	"log"
	"os"
	taskservices "study/api/internal/TaskServices"
	userservices "study/api/internal/UserServices"
	"study/api/internal/db"
	"study/api/internal/handlers"
	"study/api/internal/web/tasks"
	"study/api/internal/web/users"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// func (s *Server) handlerTasks(w http.ResponseWriter, r *http.Request) {

// 	switch r.Method {
// 	case http.MethodGet:
// 		var tasks []Task

// 		if err := s.Db.Find(&tasks).Error; err != nil {
// 			writeJSONError(w, http.StatusInternalServerError, "could not fetch tasks")
// 			return
// 		}
// 		w.Header().Set("Content-Type", "application/json")
// 		if len(tasks) == 0 {
// 			json.NewEncoder(w).Encode([]Task{})
// 			return
// 		}

// 		json.NewEncoder(w).Encode(tasks)
// 		return

// 	case http.MethodPost:
// 		var req struct {
// 			Text string `json:"text"`
// 		}
// 		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 			writeJSONError(w, http.StatusBadRequest, "bad request")
// 			return
// 		}
// 		defer r.Body.Close()
// 		if req.Text == "" {
// 			writeJSONError(w, http.StatusBadRequest, "text required")
// 			return
// 		}

// 		task := Task{
// 			Text:   req.Text,
// 			Status: "new",
// 		}
// 		if err := s.Db.Create(&task).Error; err != nil {
// 			writeJSONError(w, http.StatusInternalServerError, "could not create task")
// 			return
// 		}
// 		w.Header().Set("Location", "/tasks/"+strconv.Itoa(task.ID))
// 		w.Header().Set("Content-Type", "application/json")
// 		w.WriteHeader(http.StatusCreated)
// 		json.NewEncoder(w).Encode(task)

//		default:
//			writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
//			return
//		}
//	}
// func (s *Server) handlerTaskByID(w http.ResponseWriter, r *http.Request) {
// 	idStr := r.URL.Path[len("/tasks/"):]
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		writeJSONError(w, http.StatusBadRequest, "invalid id")
// 		return
// 	}
// 	var task Task
// 	if err := s.Db.First(&task, id).Error; err != nil {
// 		writeJSONError(w, http.StatusNotFound, "task not found")
// 		return
// 	}

// 	switch r.Method {
// 	case http.MethodPatch:
// 		var req Task
// 		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 			writeJSONError(w, http.StatusBadRequest, "bad request")
// 			return
// 		}
// 		defer r.Body.Close()
// 		if req.Text != "" {
// 			task.Text = req.Text
// 		}
// 		if req.Status != "" {
// 			task.Status = req.Status
// 		}

// 		if err := s.Db.Save(&task).Error; err != nil {
// 			writeJSONError(w, http.StatusInternalServerError, "could not update task")
// 			return
// 		}
// 		w.Header().Set("Content-Type", "application/json")
// 		json.NewEncoder(w).Encode(task)

// 	case http.MethodDelete:
// 		if err := s.Db.Delete(&task).Error; err != nil {
// 			writeJSONError(w, http.StatusInternalServerError, "could not delete task")
// 			return
// 		}
// 		w.WriteHeader(http.StatusNoContent)
// 	default:
// 		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
// 	}
// }

func main() {
	database, err := db.InitDB()
	if err != nil {
		fmt.Printf("Database error: %v", err)
		os.Exit(1)
	}
	taskRepo := taskservices.NewTaskRepository(database)
	userRepo := userservices.NewUserRepository(database)

	taskSrv := taskservices.NewTaskService(taskRepo, userRepo)
	userSrv := userservices.NewUserService(userRepo, taskRepo)

	taskHandlers := handlers.NewTaskHandlers(taskSrv)
	userHanlders := handlers.NewUserHandlers(userSrv)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	strictHandler := tasks.NewStrictHandler(taskHandlers, nil)
	tasks.RegisterHandlers(e, strictHandler)

	strictHandlerUser := users.NewStrictHandler(userHanlders, nil)
	users.RegisterHandlers(e, strictHandlerUser)

	if err := e.Start(":8080"); err != nil {
		log.Fatalf("failed to start with err: %v", err)
	}
}
