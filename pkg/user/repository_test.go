package user

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setup(t *testing.T) *userRepository {

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to connect to DB: %s", err)
	}

	return NewUserRepository(db)
}

func tearDown(userRepo *userRepository) {
	userRepo.database.Close()
}

func TestRepository(t *testing.T) {

	t.Run("AddUser inserts a user into the users table", func(t *testing.T) {
		userRepo := setup(t)
		defer tearDown(userRepo)
		userRepo.createUserTable()
		newUser := User{Username: "test", Password: "pwd", FirstName: "lou", LastName: "garwood", Email: "louis@mail.com"}

		err := userRepo.addUser(&newUser)
		if err != nil {
			t.Fatalf("error inserting user: %s", err)
		}
	})

	t.Run("FindUser returns the correct user.", func(t *testing.T) {
		userRepo := setup(t)
		defer tearDown(userRepo)

		userRepo.createUserTable()
		newUser := User{Username: "test", Password: "pwd", FirstName: "lou", LastName: "garwood", Email: "louis@mail.com"}
		err := userRepo.addUser(&newUser)

		user, err := userRepo.findUser("test")
		if err != nil {
			t.Fatalf("failed to get user: %s", err)
		}

		if user == nil {
			t.Fatal("user not found")
		}
		if user.FirstName != "lou" {
			t.Fatal("correct user not found")
		}
	})

	t.Run("RemoveUser removes the user from the database", func(t *testing.T) {
		userRepo := setup(t)
		defer tearDown(userRepo)

		userRepo.createUserTable()
		user := &User{Username: "test", Password: "test", FirstName: "brian", LastName: "boblan", Email: "lou@email.borg"}
		userRepo.addUser(user)

		err := userRepo.removeUser("test")
		if err != nil {
			t.Fatalf("failed to remove user: %s", err)
		}

		foundUser, err := userRepo.findUser("test")
		if err != nil {
			t.Fatalf("failed to find user: %s", err)
		}

		if *foundUser != (User{}) {
			t.Fatalf("failed to remove user: %s", foundUser)
		}
	})

	t.Run("updateUser updates the user from the database", func(t *testing.T) {
		userRepo := setup(t)
		defer tearDown(userRepo)

		userRepo.createUserTable()
		user := &User{Username: "test", Password: "test", FirstName: "brian", LastName: "boblan", Email: "lou@email.borg"}
		userRepo.addUser(user)

		user.LastName = "updateski"
		err := userRepo.updateUser(user)
		if err != nil {
			t.Fatalf("failed to remove user: %s", err)
		}

		foundUser, err := userRepo.findUser("test")
		if err != nil {
			t.Fatalf("failed to find user: %s", err)
		}

		if foundUser.LastName != "updateski" {
			t.Fatalf("failed to update user: %s", foundUser)
		}
	})
}
