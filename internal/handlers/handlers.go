// Package handlers exposes HTTP handlers that adapt the task service to the web API layer
package handlers

import (
	"context"
	"fmt"
	taskservices "study/api/internal/TaskServices"
	"study/api/internal/web/tasks"
)

// TaskHandlers groups HTTP handlers for task-related endpoints
type TaskHandlers struct {
	service taskservices.TaskService
}

// NewTaskHandlers returns a TaskHandlers instance configured with the given task service
func NewTaskHandlers(s taskservices.TaskService) *TaskHandlers {
	return &TaskHandlers{
		service: s,
	}
}

// GetTasks implements tasks.StrictServerInterface.
func (h *TaskHandlers) GetTasks(_ context.Context, _ tasks.GetTasksRequestObject) (tasks.GetTasksResponseObject, error) {
	alltasks, err := h.service.List()
	if err != nil {
		return nil, err
	}
	response := tasks.GetTasks200JSONResponse{}

	for _, tsk := range alltasks {
		task := tasks.Task{
			Id:     &tsk.ID,
			Task:   &tsk.Text,
			Status: &tsk.Status,
		}
		response = append(response, task)
	}
	return response, nil
}

// PostTasks implements tasks.StrictServerInterface.
func (h *TaskHandlers) PostTasks(_ context.Context, request tasks.PostTasksRequestObject) (tasks.PostTasksResponseObject, error) {
	taskRequest := request.Body
	taskToCreate := taskservices.Task{
		Text:   *taskRequest.Task,
		Status: *taskRequest.Status,
	}
	createdTask, err := h.service.Create(taskToCreate)
	if err != nil {
		return nil, err
	}
	response := tasks.PostTasks201JSONResponse{
		Id:     &createdTask.ID,
		Task:   &createdTask.Text,
		Status: &createdTask.Status,
	}
	return response, nil
}

// PatchTasksID handles partial update of a task
func (h *TaskHandlers) PatchTasksID(_ context.Context, request tasks.PatchTasksIdRequestObject) (tasks.PatchTasksIdResponseObject, error) {
	id := request.Id
	taskRequest := request.Body
	if taskRequest == nil {
		return nil, fmt.Errorf("empty body")
	}
	if taskRequest.Task == nil && taskRequest.Status == nil {
		return nil, fmt.Errorf("nothing to update")
	}
	task, err := h.service.GetByID(id)
	if err != nil {
		return nil, err
	}
	if taskRequest.Task != nil {
		task.Text = *taskRequest.Task
	}
	if taskRequest.Status != nil {
		task.Status = *taskRequest.Status
	}

	updateTask, err := h.service.Update(task.ID, task.Text, task.Status)
	if err != nil {
		return nil, err
	}
	response := tasks.PatchTasksId200JSONResponse{
		Id:     &updateTask.ID,
		Task:   &updateTask.Text,
		Status: &updateTask.Status,
	}
	return response, nil

}

// DeleteTasksID handles deletion of a task by ID
func (h *TaskHandlers) DeleteTasksID(_ context.Context, request tasks.DeleteTasksIdRequestObject) (tasks.DeleteTasksIdResponseObject, error) {
	id := request.Id
	if err := h.service.Delete(id); err != nil {
		return nil, err
	}
	return tasks.DeleteTasksId204Response{}, nil

}

// func (h *TaskHandlers) HandleTasks(w http.ResponseWriter, r *http.Request) {
// 	switch r.Method {
// 	case http.MethodGet:

// 	case http.MethodPost:
// 		var req taskservices.TaskRequest
// 		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 			writeJSONError(w, http.StatusBadRequest, "bad request")
// 			return
// 		}
// 		defer r.Body.Close()
// 		task, err := h.service.Create(req.Task)
// 		if err != nil {
// 			writeJSONError(w, http.StatusBadRequest, "could not create task")
// 			return
// 		}
// 		w.Header().Set("Location", "/tasks/"+strconv.Itoa(task.ID))
// 		writeJSON(w, http.StatusOK, task)
// 	default:
// 		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
// 		return
// 	}
// }
// func (h *TaskHandlers) HandleTaskByID(w http.ResponseWriter, r *http.Request) {
// 	idStr := r.URL.Path[len("/tasks/"):]
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		writeJSONError(w, http.StatusBadRequest, "invalid id")
// 		return
// 	}
// 	switch r.Method {
// 	case http.MethodPatch:
// 		var req taskservices.Task
// 		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 			writeJSONError(w, http.StatusBadRequest, "bad request")
// 			return
// 		}
// 		defer r.Body.Close()

// 		updateTask, err := h.service.Update(id, req.Text, req.Status)
// 		if err != nil {
// 			writeJSONError(w, http.StatusInternalServerError, "could not update task")
// 			return
// 		}
// 		writeJSON(w, http.StatusOK, updateTask)

// 	case http.MethodDelete:
// 		if err := h.service.Delete(id); err != nil {
// 			writeJSONError(w, http.StatusInternalServerError, "could not delete task")
// 			return
// 		}
// 		w.WriteHeader(http.StatusNoContent)
// 	default:
// 		writeJSONError(w, http.StatusMethodNotAllowed, "Method not allowed")
// 	}

// }

// func writeJSONError(w http.ResponseWriter, code int, msg string) {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(code)
// 	json.NewEncoder(w).Encode(map[string]string{
// 		"error":   http.StatusText(code),
// 		"message": msg,
// 	})
// }
// func writeJSON(w http.ResponseWriter, code int, v any) {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(code)
// 	_ = json.NewEncoder(w).Encode(v)
// }
