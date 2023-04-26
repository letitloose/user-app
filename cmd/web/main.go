package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/letitloose/user-app/cmd/config"
	"github.com/letitloose/user-app/pkg/server"
	"github.com/letitloose/user-app/pkg/user"
)

func main() {

	err := run("app-config.yml")
	if err != nil {
		log.Fatalf("error starting application: %s", err)
	}
}

func run(configFile string) error {
	config := config.GetConfig()
	err := config.ReadConfig(configFile)
	if err != nil {
		return errors.New(fmt.Sprintf("error reading config file: %s", err))
	}

	db, err := setupDatabase(config)
	if err != nil {
		return errors.New(fmt.Sprintf("error setting up database: %s", err))
	}

	userRepo := user.NewUserRepository(db)
	userService := user.NewUserService(userRepo)

	server := server.NewServer(config, userService)
	err = server.Run()
	if err != nil {
		log.Fatalf("error starting server:%s\n", err)
	}
	return nil
}

func setupDatabase(config *config.Config) (*sql.DB, error) {

	connString := assembleConnectString(config)
	db, err := sql.Open("mysql", connString)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	log.Println("db connection successful!")

	return db, nil
}

func assembleConnectString(config *config.Config) string {

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		config.Db.Username,
		config.Db.Password,
		config.Db.Host,
		config.Db.Port,
		config.Db.Database)

	return connectionString
}
