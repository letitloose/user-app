package user

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func setupHandlers(t *testing.T) *UserService {

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to connect to DB: %s", err)
	}

	userRepo := NewUserRepository(db)
	err = userRepo.createUserTable()
	if err != nil {
		t.Fatalf("failed to create user table: %s", err)
	}
	newUser := User{Username: "test", Password: "pwd", FirstName: "lou", LastName: "garwood", Email: "louis@mail.com"}
	err = userRepo.addUser(&newUser)
	if err != nil {
		t.Fatalf("failed to add user: %s", err)
	}

	return NewUserService(userRepo)
}

func teardownHandlers(service *UserService) {
	service.repository.database.Close()
}

func TestHandlers(t *testing.T) {

	t.Run("test listUsers", func(t *testing.T) {
		userService := setupHandlers(t)
		defer teardownHandlers(userService)
		request, err := http.NewRequest("GET", "/users", nil)
		if err != nil {
			t.Fatal(err)
		}

		recorder := httptest.NewRecorder()
		handler := http.HandlerFunc(userService.ServeHTTP)

		request.Header.Set("Content-Type", "application/json")
		handler.ServeHTTP(recorder, request)

		if status := recorder.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		expected := `[{"user-name":"test","password":"pwd","first-name":"lou","last-name":"garwood","email":"louis@mail.com"}]`
		if recorder.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v",
				recorder.Body.String(), expected)
		}
	})

	t.Run("test findUserByName", func(t *testing.T) {
		userService := setupHandlers(t)
		defer teardownHandlers(userService)
		request, err := http.NewRequest("GET", "/users/test", nil)
		request.Header.Set("Content-Type", "application/json")
		if err != nil {
			t.Fatal(err)
		}

		recorder := httptest.NewRecorder()
		handler := http.HandlerFunc(userService.ServeHTTP)

		handler.ServeHTTP(recorder, request)

		if status := recorder.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		expected := `{"user-name":"test","password":"pwd","first-name":"lou","last-name":"garwood","email":"louis@mail.com"}`
		if recorder.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v",
				recorder.Body.String(), expected)
		}
	})

	t.Run("test removeUser", func(t *testing.T) {
		userService := setupHandlers(t)
		defer teardownHandlers(userService)
		request, err := http.NewRequest("DELETE", "/users/test", nil)
		request.Header.Set("Content-Type", "application/json")
		if err != nil {
			t.Fatal(err)
		}

		recorder := httptest.NewRecorder()
		handler := http.HandlerFunc(userService.ServeHTTP)

		handler.ServeHTTP(recorder, request)

		if status := recorder.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		expected := "user successfully deleted"
		if recorder.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v",
				recorder.Body.String(), expected)
		}

		request, err = http.NewRequest("GET", "/users/test", nil)
		request.Header.Set("Content-Type", "application/json")
		if err != nil {
			t.Fatal(err)
		}

		recorder = httptest.NewRecorder()
		handler = http.HandlerFunc(userService.ServeHTTP)

		handler.ServeHTTP(recorder, request)

		expected = `{"user-name":"","password":"","first-name":"","last-name":"","email":""}`
		if recorder.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v",
				recorder.Body.String(), expected)
		}
	})

	t.Run("test createUser", func(t *testing.T) {
		userService := setupHandlers(t)
		defer teardownHandlers(userService)
		user := &User{Username: "test1", Password: "pass", FirstName: "lou", LastName: "gar"}
		userJson, err := json.Marshal(user)
		request, err := http.NewRequest("POST", "/users", bytes.NewBuffer(userJson))
		request.Header.Set("Content-Type", "application/json")
		if err != nil {
			t.Fatal(err)
		}

		recorder := httptest.NewRecorder()
		handler := http.HandlerFunc(userService.ServeHTTP)

		handler.ServeHTTP(recorder, request)

		if status := recorder.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		expected := "user successfully added"
		if recorder.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v",
				recorder.Body.String(), expected)
		}

		request, err = http.NewRequest("GET", "/users/test1", nil)
		request.Header.Set("Content-Type", "application/json")
		if err != nil {
			t.Fatal(err)
		}

		recorder = httptest.NewRecorder()
		handler = http.HandlerFunc(userService.ServeHTTP)

		handler.ServeHTTP(recorder, request)

		expected = `{"user-name":"test1","password":"pass","first-name":"lou","last-name":"gar","email":""}`
		if recorder.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v",
				recorder.Body.String(), expected)
		}
	})

	t.Run("test updateUser", func(t *testing.T) {
		userService := setupHandlers(t)
		defer teardownHandlers(userService)
		user := &User{Username: "test", Password: "pass", FirstName: "lou", LastName: "gar"}
		userJson, err := json.Marshal(user)
		request, err := http.NewRequest("PUT", "/users/test", bytes.NewBuffer(userJson))
		request.Header.Set("Content-Type", "application/json")
		if err != nil {
			t.Fatal(err)
		}

		recorder := httptest.NewRecorder()
		handler := http.HandlerFunc(userService.ServeHTTP)

		handler.ServeHTTP(recorder, request)

		if status := recorder.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		expected := "user successfully updated"
		if recorder.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v",
				recorder.Body.String(), expected)
		}

		request, err = http.NewRequest("GET", "/users/test", nil)
		request.Header.Set("Content-Type", "application/json")
		if err != nil {
			t.Fatal(err)
		}

		recorder = httptest.NewRecorder()
		handler = http.HandlerFunc(userService.ServeHTTP)

		handler.ServeHTTP(recorder, request)

		expected = `{"user-name":"test","password":"pass","first-name":"lou","last-name":"gar","email":""}`
		if recorder.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v",
				recorder.Body.String(), expected)
		}
	})

	t.Run("test renderResponse returns json if content-type is json", func(t *testing.T) {

		userList := []User{{Username: "lou"}}
		response := httptest.NewRecorder()
		response.Header().Set("Content-Type", "application/json")
		renderResponse(response, userList, "")

		expected := `[{"user-name":"lou","password":"","first-name":"","last-name":"","email":""}]`
		if response.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v",
				response.Body.String(), expected)
		}
	})

	t.Run("test renderResponse returns json if content-type is not json", func(t *testing.T) {

		os.WriteFile("user.tmpl", []byte("<html>"), 0755)
		defer os.Remove("user.tmpl")
		userList := []User{{Username: "lou"}}
		response := httptest.NewRecorder()
		response.Header().Set("Content-Type", "application/text")
		renderResponse(response, userList, "user.tmpl")

		expected := "<html>"
		if response.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v",
				response.Body.String(), expected)
		}
	})
}
