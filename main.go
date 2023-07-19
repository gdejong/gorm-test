package main

import (
	"errors"
	"github.com/go-faker/faker/v4"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"strconv"
)

const (
	dsn = "test.db"

	commandCreate = "create"
	commandDelete = "delete"
	commandRead   = "read"
)

func main() {
	// Create a logger.
	logger := logrus.New()
	logger.SetOutput(os.Stdout)

	// Make sure we have a command line argument.
	if len(os.Args) < 2 {
		logger.Info("Usage: go run . <command>")
		os.Exit(1)
	}

	command := os.Args[1]

	logger.WithField("dsn", dsn).Info("Opening database")
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		// TODO can we use the logrus logger here?
		// Logger:
	})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	// Migrate the schema
	logger.Info("Auto migrating database")
	err = db.AutoMigrate(&User{})
	if err != nil {
		panic("failed to auto migrate: " + err.Error())
	}

	switch command {
	case commandCreate:
		handleCreate(db, logger)
	case commandDelete:
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			logger.WithError(err).Fatal("Failed to parse id")
		}
		handleDelete(db, logger, id)
	case commandRead:
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			logger.WithError(err).Fatal("Failed to parse id")
		}
		handleRead(db, logger, id)
	}
}

func handleRead(db *gorm.DB, logger *logrus.Logger, id int) {
	user := &User{}
	result := db.First(user, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		logger.WithField("id", id).Warn("User not found")

		return
	}

	logger.WithFields(logrus.Fields{
		"id":        user.ID,
		"name":      user.Name,
		"createdAt": user.CreatedAt,
	}).Info("User")
}

func handleDelete(db *gorm.DB, logger *logrus.Logger, id int) {
	res := db.Delete(&User{}, id)

	if res.Error != nil {
		logger.WithError(res.Error).Fatal("Failed to delete user")
	}

	if res.RowsAffected == 0 {
		logger.WithField("id", id).Warn("User not found")

		return
	}

	logger.WithField("id", id).Info("Deleted user")
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
