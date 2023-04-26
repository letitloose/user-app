package user

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setupService(t *testing.T) *UserService {

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

func teardownService(service *UserService) {
	service.repository.database.Close()
}

func TestUserService(t *testing.T) {

	t.Run("ListAllUsers returns a user list", func(t *testing.T) {
		userService := setupService(t)
		defer teardownService(userService)

		userList := userService.ListAllUsers()

		userListType := fmt.Sprintf("%T", userList)

		if userListType != "[]*user.User" {
			t.Fatal("did not return a list of users")
		}
	})

	t.Run("FindUserByName returns a user", func(t *testing.T) {
		userService := setupService(t)
		defer teardownService(userService)
		user, err := userService.FindByUsername("test")
		if err != nil {
			t.Fatalf("error finding user: %s", err)
		}
		if user.FirstName != "lou" {
			t.Fatalf("did not return a correct user, expected: lou, got:%s", user.FirstName)
		}
	})

	t.Run("AddUser adds a user", func(t *testing.T) {
		userService := setupService(t)
		defer teardownService(userService)
		err := userService.RemoveUser("test")
		if err != nil {
			t.Fatalf("error removing user: %s", err)
		}

		user, err := userService.FindByUsername("test")
		if err != nil {
			t.Fatalf("error finding user: %s", err)
		}
		if user.FirstName != "" {
			t.Fatalf("did not return a correct user, expected: lou, got:%s", user.FirstName)
		}
	})

	t.Run("RemoveUser removes a user", func(t *testing.T) {
		userService := setupService(t)
		defer teardownService(userService)

		newUser := User{Username: "test1", Password: "pwd", FirstName: "addison", LastName: "garwood", Email: "louis@mail.com"}
		err := userService.AddUser(&newUser)
		if err != nil {
			t.Fatalf("error adding user: %s", err)
		}

		user, err := userService.FindByUsername("test1")
		if err != nil {
			t.Fatalf("error finding user: %s", err)
		}
		if user.FirstName != "addison" {
			t.Fatalf("did not return a correct user, expected: addison, got:%s", user.FirstName)
		}
	})

	t.Run("RemoveUser removes a user", func(t *testing.T) {
		userService := setupService(t)
		defer teardownService(userService)

		user, err := userService.FindByUsername("test")
		if err != nil {
			t.Fatalf("error finding test user: %s", err)
		}
		user.LastName = "banks"
		err = userService.UpdateUser(user)
		if err != nil {
			t.Fatalf("error adding user: %s", err)
		}

		foundUser, err := userService.FindByUsername("test")
		if err != nil {
			t.Fatalf("error finding user: %s", err)
		}
		if foundUser.LastName != "banks" {
			t.Fatalf("did not return a correct user, expected: banks, got:%s", user.LastName)
		}
	})
}
