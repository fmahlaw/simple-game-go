package db

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// ConnectToPostgres connects to the PostgreSQL database and returns the *gorm.DB instance
func ConnectToPostgres(c *gin.Context) (*gorm.DB, error) {
	connStr := "user=postgres dbname=gaming_app password=mahlaw17 sslmode=disable"
	db, err := gorm.Open("postgres", connStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to the database"})
		return nil, err
	}
	return db, nil
}
