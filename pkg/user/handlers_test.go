package user

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
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
		request, err := http.NewRequest("GET", "/api/users", nil)
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

		expected := `[{"user-name":"test","password":"pwd","first-name":"lou","last-name":"garwood","email":"louis@mail.com"}]`
		if recorder.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v",
				recorder.Body.String(), expected)
		}
	})

	t.Run("test findUserByName", func(t *testing.T) {
		userService := setupHandlers(t)
		defer teardownHandlers(userService)
		request, err := http.NewRequest("GET", "/api/users/test", nil)
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
		request, err := http.NewRequest("DELETE", "/api/users/test", nil)
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

		request, err = http.NewRequest("GET", "/api/users/test", nil)
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
		request, err := http.NewRequest("POST", "/api/users", bytes.NewBuffer(userJson))
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

		request, err = http.NewRequest("GET", "/api/users/test1", nil)
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
		request, err := http.NewRequest("PUT", "/api/users/test", bytes.NewBuffer(userJson))
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

		request, err = http.NewRequest("GET", "/api/users/test", nil)
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
}
