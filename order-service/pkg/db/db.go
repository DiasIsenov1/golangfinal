package db

import (
	"log"
	"order-service/pkg/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"order-service/pkg/models"
)

type Handler struct {
	DB     *gorm.DB
	Logger *log.Logger
}

func Init(url string) Handler {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	db.AutoMigrate(&models.Order{})

	return Handler{DB: db, Logger: logger.Info}
}
