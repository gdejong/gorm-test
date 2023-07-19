package main

import (
	"github.com/go-faker/faker/v4"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
)

const dsn = "test.db"
const commandCreate = "create"

func main() {
	// Create a logger.
	logger := logrus.New()
	logrus.SetOutput(os.Stdout)

	// Make sure we have a command line argument.
	if len(os.Args) != 2 {
		logger.Info("Usage: go run . <command>")
		os.Exit(1)
	}

	command := os.Args[1]

	logger.WithField("dsn", dsn).Info("Opening database")
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	// Migrate the schema
	logger.Info("Auto migrating database")
	err = db.AutoMigrate(&User{})
	if err != nil {
		panic("failed to auto migrate: " + err.Error())
	}

	if command == commandCreate {
		handleCreate(db, logger)
	}

}

func handleCreate(db *gorm.DB, logger *logrus.Logger) {
	name := faker.Name()

	user := &User{
		Name: name,
	}

	res := db.Create(user)

	if res.Error != nil {
		logger.WithError(res.Error).Fatal("Failed to create user")
	}

	logger.WithFields(logrus.Fields{
		"id":   user.ID,
		"name": user.Name,
	}).Info("Created user")
}
