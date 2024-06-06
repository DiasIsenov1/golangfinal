package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"product/pkg/models"
)

type Handler struct {
	DB *gorm.DB
}

func Init(url string) Handler {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	db.AutoMigrate(&models.Product{})
	db.AutoMigrate(&models.StockDecreaseLog{})

	return Handler{db}
}
func (h *Handler) Close() {
	sqlDB, err := h.DB.DB()
	if err != nil {
		log.Println("Failed to get underlying database connection:", err)
		return
	}
	if err := sqlDB.Close(); err != nil {
		log.Println("Failed to close database connection:", err)
	}
}
