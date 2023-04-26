package main

import (
	"log"
	"os"
	"testing"

	"github.com/go-yaml/yaml"
	"github.com/letitloose/user-app/cmd/config"
)

func setup() {
	var config = config.Config{}
	config.Db.Database = "database"
	config.Db.Host = "localhost"
	config.Db.Port = 3306
	config.Db.Password = "pass"
	config.Db.Username = "user"

	configFile, err := os.Create("config.yml")
	if err != nil {
		log.Fatalf("could not mock config %s", err)
	}
	defer configFile.Close()

	configBytes, err := yaml.Marshal(config)
	if err != nil {
		log.Fatalf("failed to marshal dummy config to bytes:%s", err)
	}

	err = os.WriteFile("config.yml", configBytes, 0644)
	if err != nil {
		log.Fatalf("failed to write config to disk:%s", err)

	}
}

func tearDown() {
	os.Remove("config.yml")
}

func TestMain(t *testing.T) {
	t.Run("assembleConnectString creates the string from the config", func(t *testing.T) {
		setup()
		defer tearDown()

		config := config.GetConfig()
		err := config.ReadConfig("config.yml")
		if err != nil {
			t.Fatalf("error reading config file:%s", err)
		}
		connString := assembleConnectString(config)
		if connString != "user:pass@tcp(localhost:3306)/database?parseTime=true" {
			t.Fatalf("connect string did not come out correctly: %s", connString)
		}
	})
}
