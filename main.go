package main

import (
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
)

const dsn = "test.db"

func main() {
	logger := logrus.New()
	logrus.SetOutput(os.Stdout)

	logger.Info("Opening database")
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

}
