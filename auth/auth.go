package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"simple-game-go/config"
	"simple-game-go/db"
	"simple-game-go/models"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

var redisClient = db.RedisClient()
var user models.User
var storedPassword models.StoredPassword

func LogoutHandler(c *gin.Context) {

	// Extract the token from the request headers
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token not provided"})
		return
	}

	// Here, you would typically validate the token and get the claims

	// Invalidation example (remove token from Redis)
	err := redisClient.Del(context.Background(), tokenString).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to invalidate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func LoginHandler(c *gin.Context) {

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Retrieve the stored password for the provided username from Redis
	storedPassword, err := getPasswordFromRedis(user.Username, redisClient)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Check if the provided username and password match the stored password
	if user.Password == storedPassword {
		// Generate a JWT token
		token, err := generateJWTToken(user.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": token})
		return
	}

	c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
}

func getPasswordFromRedis(username string, client *redis.Client) (string, error) {
	// Retrieve the password for the given username from Redis
	storedPasswordJSON, err := client.Get(context.Background(), username).Result()
	if err != nil {
		if err == redis.Nil {
			return "", fmt.Errorf("username not found")
		}
		return "", fmt.Errorf("failed to retrieve password from Redis: %v", err)
	}

	// Unmarshal the JSON string into a struct

	if err := json.Unmarshal([]byte(storedPasswordJSON), &storedPassword); err != nil {
		return "", fmt.Errorf("failed to unmarshal stored password: %v", err)
	}

	// Return the password
	return storedPassword.Password, nil
}

func generateJWTToken(username string) (string, error) {
	// Create the token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Token expiration time: 24 hours

	// Sign the token with the secret key
	tokenString, err := token.SignedString(config.JWTSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
