package bank

import (
	"fmt"
	"net/http"
	"simple-game-go/db"
	"simple-game-go/models"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func RegisterBankAccountHandler(c *gin.Context) {
	// Extract the username from the JWT token
	username, err := getUsernameFromJWT(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or missing JWT token"})
		return
	}

	// Bind the bank account data from the request
	var bankAccount models.BankAccount
	if err := c.ShouldBindJSON(&bankAccount); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println(bankAccount)
	// Set the username
	bankAccount.Username = username

	// Save the bank account to the database using GORM
	if err := saveBankAccount(bankAccount); err != nil {
		// Check for the specific error indicating a duplicate account
		if strings.Contains(err.Error(), "bank account already registered for account number") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Bank account already registered for this account number"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register bank account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bank account registered successfully"})
}

func getUsernameFromJWT(c *gin.Context) (string, error) {
	// Extract the JWT token from the request header
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		return "", fmt.Errorf("JWT token not provided")
	}

	// Parse the JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Replace this with your actual JWT secret
		return []byte("your_jwt_secret_key"), nil
	})
	if err != nil || !token.Valid {
		return "", fmt.Errorf("invalid or expired JWT token")
	}

	// Extract the username from the JWT claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("failed to extract claims from JWT token")
	}

	username, ok := claims["username"].(string)
	if !ok {
		return "", fmt.Errorf("failed to extract username from JWT token")
	}

	fmt.Println(username)

	return username, nil
}

func saveBankAccount(bankAccount models.BankAccount) error {
	// Connect to the PostgreSQL database
	db, err := db.ConnectToPostgres(nil)
	if err != nil {
		return err
	}
	defer db.Close()

	var existingAccount models.BankAccount
	if err := db.Table("accounts").Where("username = ?", bankAccount.Username).First(&existingAccount).Error; err == nil {
		// An account already exists for this account number
		return fmt.Errorf("bank account already registered for account number: %s", bankAccount.AccountNumber)
	}

	// Save the bank account to the database using GORM
	if err := db.Table("accounts").Create(&bankAccount).Error; err != nil {
		return err
	}

	return nil
}

// Handler to top up the wallet balance
func TopUpBalanceHandler(c *gin.Context) {
	// Extract the username from the JWT token
	username, err := getUsernameFromJWT(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or missing JWT token"})
		return
	}

	// Bind the balance data from the request
	var payload struct {
		Balance float64 `json:"balance"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save the wallet balance
	err = saveWalletBalance(username, payload.Balance)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to top up wallet balance"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Wallet balance topped up successfully"})
}

func saveWalletBalance(username string, balance float64) error {
	// Connect to the PostgreSQL database
	db, err := db.ConnectToPostgres(nil)
	if err != nil {
		return err
	}
	defer db.Close()

	// Check if a wallet already exists for the given username
	var existingWallet models.Wallet
	if err := db.Where("username = ?", username).First(&existingWallet).Error; err == nil {
		// Wallet already exists, update the balance
		existingWallet.Balance += balance
		if err := db.Save(&existingWallet).Error; err != nil {
			return fmt.Errorf("failed to update wallet balance: %v", err)
		}
		return nil
	}

	// Create a new Wallet record
	wallet := models.Wallet{
		Username: username,
		Balance:  balance,
	}

	// Save the wallet to the database using GORM
	if err := db.Create(&wallet).Error; err != nil {
		return fmt.Errorf("failed to create wallet: %v", err)
	}

	return nil
}
