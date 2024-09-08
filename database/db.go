package database

import (
	"fmt"
	"hacktiv8-go-assignment-2/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"log"
)

var (
	host     = "localhost"
	user     = "postgres"
	password = "postgres"
	dbPort   = "5432"
	dbname   = "postgres"
	db       *gorm.DB
	err      error
)

func StartDB() {
	config := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbname, dbPort)
	db, err = gorm.Open(postgres.Open(config), &gorm.Config{})
	if err != nil {
		log.Fatal("error connecting to database :", err)
	}

	err = db.Debug().AutoMigrate(&models.Order{}, &models.Item{})
	if err != nil {
		log.Fatal("ERROR:", err)
	}
}
func GetDB() *gorm.DB {
	return db
}
