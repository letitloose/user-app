package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (userService *UserService) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	username, _ := extractUsername(request)
	switch request.Method {
	case http.MethodGet:
		if username == "" {
			userService.listUsers(writer, request)
			return
		}
		userService.findByUsername(writer, request)
	case http.MethodPost:
		userService.createUser(writer, request)
	case http.MethodPut:
		userService.updateUser(writer, request)
	case http.MethodDelete:
		userService.deleteUser(writer, request)
	default:
		msg := http.StatusText(http.StatusNotFound)
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, msg)
	}
}

func (userService *UserService) AddHandlersToMux(mux *http.ServeMux) {
	mux.HandleFunc("/api/users", userService.ServeHTTP)
	mux.HandleFunc("/api/users/", userService.ServeHTTP)
}

func (userService *UserService) listUsers(writer http.ResponseWriter, request *http.Request) {
	users := userService.ListAllUsers()
	bytes, _ := json.Marshal(users)

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(bytes)
}

func extractUsername(r *http.Request) (string, error) {
	urlPath := strings.Split(r.URL.Path[1:], "/")
	if !pathHasUsername(urlPath) {
		return "", fmt.Errorf("no username provided")
	}
	return urlPath[len(urlPath)-1], nil
}

func pathHasUsername(pathElements []string) bool {
	if pathElements[0] == "api" {
		return len(pathElements) == 3
	}

	return len(pathElements) == 2
}

func (userService *UserService) findByUsername(writer http.ResponseWriter, request *http.Request) {
	username, err := extractUsername(request)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, err.Error())
		return
	}

	users, err := userService.FindByUsername(username)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, err.Error())
		return
	}

	bytes, _ := json.Marshal(users)

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(bytes)
}

func (userService *UserService) deleteUser(writer http.ResponseWriter, request *http.Request) {
	username, err := extractUsername(request)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, err.Error())
		return
	}

	err = userService.RemoveUser(username)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, err.Error())
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("user successfully deleted"))
}

func (userService *UserService) updateUser(writer http.ResponseWriter, request *http.Request) {
	username, err := extractUsername(request)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, err.Error())
		return
	}

	var user = &User{}
	decoder := json.NewDecoder(request.Body)
	err = decoder.Decode(user)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, err.Error())
		return
	}

	if username != user.Username {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, "wrong user specified")
		return
	}

	err = userService.UpdateUser(user)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, err.Error())
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("user successfully updated"))
}

func (userService *UserService) createUser(writer http.ResponseWriter, request *http.Request) {
	var user = &User{}
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(user)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, err.Error())
		return
	}

	err = userService.AddUser(user)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(writer, err.Error())
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("user successfully added"))
}
