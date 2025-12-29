package handlers

import (
	"context"
	"errors"
	"fmt"
	userservices "study/api/internal/UserServices"
	"study/api/internal/models"
	"study/api/internal/web/users"
)

// UserHandlers groups HTTP handlers for task-related endpoints
type UserHandlers struct {
	service userservices.UserService
}

// NewUserHandlers returns a UserHandlers instance configured with the given user service
func NewUserHandlers(s userservices.UserService) *UserHandlers {
	return &UserHandlers{
		service: s,
	}
}

// GetUsersIDTasks implements users.StrictServerInterface.
func (u *UserHandlers) GetUsersIDTasks(_ context.Context, request users.GetUsersIDTasksRequestObject) (users.GetUsersIDTasksResponseObject, error) {
	userID := request.Id
	tasks, err := u.service.GetTasksForUser(userID)

	if errors.Is(err, userservices.ErrUserNotFound) {
		return users.GetUsersIDTasks404Response{}, nil
	}

	if err != nil {
		return nil, err
	}

	response := users.GetUsersIDTasks200JSONResponse{}

	for _, t := range tasks {
		userTasks := users.Task{
			Id:     &t.ID,
			Task:   &t.Text,
			Status: &t.Status,
			UserId: &t.UserID,
		}
		response = append(response, userTasks)
	}
	return response, nil
}

// GetUsers implements users.StrictServerInterface.
func (u *UserHandlers) GetUsers(_ context.Context, _ users.GetUsersRequestObject) (users.GetUsersResponseObject, error) {
	allUsers, err := u.service.List()
	if err != nil {
		return nil, err
	}
	response := users.GetUsers200JSONResponse{}

	for _, usr := range allUsers {
		userRep := users.User{
			Id:        &usr.ID,
			Email:     &usr.Email,
			Password:  &usr.Password,
			CreatedAt: &usr.CreatedAt,
			UpdatedAt: &usr.UpdatedAt,
		}
		response = append(response, userRep)
	}
	return response, nil
}

// PostUsers implements users.StrictServerInterface.
func (u *UserHandlers) PostUsers(_ context.Context, request users.PostUsersRequestObject) (users.PostUsersResponseObject, error) {
	userRequest := request.Body
	userToCreate := models.User{
		Email:    userRequest.Email,
		Password: userRequest.Password,
	}
	createUser, err := u.service.Create(userToCreate)
	if err != nil {
		return nil, err
	}

	response := users.PostUsers201JSONResponse{
		Id:        &createUser.ID,
		Email:     &createUser.Email,
		Password:  &createUser.Password,
		CreatedAt: &createUser.CreatedAt,
	}

	return response, nil
}

// PatchUsersID implements users.StrictServerInterface.
func (u *UserHandlers) PatchUsersID(_ context.Context, request users.PatchUsersIDRequestObject) (users.PatchUsersIDResponseObject, error) {
	id := request.Id
	userRequest := request.Body
	if userRequest == nil {
		return nil, fmt.Errorf("empty body")
	}
	if userRequest.Email == nil && userRequest.Password == nil {
		return nil, fmt.Errorf("nothing to update")
	}
	user, err := u.service.GetByID(id)
	if err != nil {
		return nil, err
	}
	email, password, err := checkUserPatch(user, userRequest, u.service)
	if err != nil {
		return nil, err
	}

	updatedUser, err := u.service.Update(user.ID, email, password)
	if err != nil {
		return nil, err
	}

	response := users.PatchUsersID200JSONResponse{
		Id:        &updatedUser.ID,
		Email:     &updatedUser.Email,
		Password:  &updatedUser.Password,
		UpdatedAt: &updatedUser.UpdatedAt,
	}
	return response, nil
}

// DeleteUsersID implements users.StrictServerInterface.
func (u *UserHandlers) DeleteUsersID(_ context.Context, request users.DeleteUsersIDRequestObject) (users.DeleteUsersIDResponseObject, error) {
	id := request.Id
	if err := u.service.Delete(id); err != nil {
		if errors.Is(err, userservices.ErrUserNotFound) {
			return users.DeleteUsersID404Response{}, nil
		}
		return nil, err
	}

	return users.DeleteUsersID204Response{}, nil
}

func checkUserPatch(currentUser models.User, patchRequest *users.PatchUsersIDJSONRequestBody, srv userservices.UserService) (string, string, error) {
	updatedEmail := currentUser.Email

	if patchRequest.Email != nil {
		if *patchRequest.Email == "" {
			return "", "", fmt.Errorf("email is empty")
		}
		if *patchRequest.Email != currentUser.Email {
			_, err := srv.GetByEmail(*patchRequest.Email)
			if err == nil {
				return "", "", fmt.Errorf("email already exists")
			}
			if !errors.Is(err, userservices.ErrUserNotFound) {
				return "", "", fmt.Errorf("error db")
			}
		}
		updatedEmail = *patchRequest.Email
	}

	updatedPassword := currentUser.Password

	if patchRequest.Password != nil {
		if *patchRequest.Password == "" {
			return "", "", fmt.Errorf("password is empty")
		}
		updatedPassword = *patchRequest.Password
	}
	return updatedEmail, updatedPassword, nil
}
