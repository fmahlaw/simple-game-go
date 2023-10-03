package player

import (
	"context"
	"encoding/json"
	"net/http"
	"simple-game-go/db"
	"simple-game-go/models"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

var redisClient = db.RedisClient()
var user models.User

func GetAllPlayersHandler(c *gin.Context) {
	// Connect to the PostgreSQL database
	db, err := db.ConnectToPostgres(nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to the database"})
		return
	}
	defer db.Close()

	var players []models.Player
	if err := db.Select("id,username, email").Find(&players).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch players"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"players": players})

}

func GetPlayerWithWalletHandler(c *gin.Context) {
	connStr := "user=postgres dbname=gaming_app password=mahlaw17 sslmode=disable"
	db, _ := gorm.Open("postgres", connStr)
	username := c.Param("username")

	var player models.PlayerResponse
	query := `
		SELECT p.id, p.username, p.email, p.password
		    , w.id as wallet_id, w.balance as wallet_balance
		     , a.id as account_id, a.account_name, a.account_number, a.bank_name
		FROM players p
		LEFT JOIN wallets w ON p.username = w.username
		LEFT JOIN accounts a ON p.username = a.username
		WHERE p.username = ?`

	if err := db.Raw(query, username).Scan(&player).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Player not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"player": player})
}

func RegisterHandler(c *gin.Context) {
	// Initialize the DB connection
	db, err := db.ConnectToPostgres(nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to the database"})
		return
	}
	defer db.Close()

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the username already exists in the database
	var count int
	db.Model(&models.Player{}).Where("username = ?", user.Username).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
		return
	}

	newPlayer := models.Player{
		Username: user.Username,
		Password: user.Password,
		Email:    user.Email,
	}
	db.Model(&models.Player{}).Create(&newPlayer)

	// Store the user in Redis
	err = storeUserInRedis(context.Background(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store user in Redis"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func storeUserInRedis(ctx context.Context, user models.User) error {
	// Serialize the user to JSON
	userJSON, err := json.Marshal(user)
	if err != nil {
		return err
	}

	// Store the user data in Redis with a key based on the username
	err = redisClient.Set(ctx, user.Username, userJSON, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func DashboardHandler(c *gin.Context) {
	// You can assume that if the handler is reached, the token is valid
	c.JSON(http.StatusOK, gin.H{"message": "Token is valid"})
}
