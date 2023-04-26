package user

import (
	"database/sql"
	"errors"
)

type User struct {
	Username  string `json:"user-name"`
	Password  string `json:"password"`
	FirstName string `json:"first-name"`
	LastName  string `json:"last-name"`
	Email     string `json:"email"`
}

type userRepository struct {
	database *sql.DB
}

func NewUserRepository(database *sql.DB) *userRepository {
	return &userRepository{database: database}
}

func (repository *userRepository) listAll() []*User {
	rows, _ := repository.database.Query(`SELECT username, password, firstname, lastname, email FROM users;`)
	defer rows.Close()

	users := []*User{}

	for rows.Next() {
		var (
			username  string
			password  string
			firstname string
			lastname  string
			email     string
		)

		rows.Scan(&username, &password, &firstname, &lastname, &email)

		users = append(users, &User{
			Username:  username,
			Password:  password,
			FirstName: firstname,
			LastName:  lastname,
			Email:     email,
		})
	}

	return users
}

func (repository *userRepository) createUserTable() error {

	_, err := repository.database.Exec(`create table if not exists users (username varchar(255) unique, 
		password varchar(255), 
		firstname varchar(255), 
		lastname varchar(255), 
		email varchar(255));`)
	if err != nil {
		return err
	}

	return nil
}

func (repository *userRepository) userTableExists() bool {

	rows, err := repository.database.Query(`desc users`)
	if rows != nil {
		defer rows.Close()
	}

	if err == nil {
		return true
	}

	return false
}

func (repository *userRepository) addUser(user *User) error {

	insertStatement := "insert into users (username, password, firstname, lastname, email) values (?, ?, ?, ?, ?)"

	result, err := repository.database.Exec(insertStatement, user.Username, user.Password, user.FirstName, user.LastName, user.Email)
	if err != nil {
		return err
	}

	numRows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if numRows != 1 {
		return errors.New("user not inserted")
	}
	return nil
}

func (repository *userRepository) updateUser(user *User) error {

	updateStatment := "update users set password=?, firstname=?, lastname=?, email=? where username=?;"

	result, err := repository.database.Exec(updateStatment, user.Password, user.FirstName, user.LastName, user.Email, user.Username)
	if err != nil {
		return err
	}

	numRows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if numRows != 1 {
		return errors.New("user not updated")
	}
	return nil
}

func (repository *userRepository) findUser(usernameParam string) (*User, error) {
	query := "select username,  password, firstname, lastname, email from users where username = ?"

	rows := repository.database.QueryRow(query, usernameParam)

	var (
		username  string
		password  string
		firstname string
		lastname  string
		email     string
	)

	rows.Scan(&username, &password, &firstname, &lastname, &email)

	user := &User{
		Username:  username,
		Password:  password,
		FirstName: firstname,
		LastName:  lastname,
		Email:     email,
	}

	return user, nil
}

func (repository *userRepository) removeUser(username string) error {

	deleteQuery := "delete from users where username = ?;"
	result, err := repository.database.Exec(deleteQuery, username)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected != 1 {
		return errors.New("wrong number of rows affected")
	}
	return nil
}
