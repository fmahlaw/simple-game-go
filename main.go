package main

import (
	"simple-game-go/auth"
	"simple-game-go/bank"
	"simple-game-go/middleware"
	"simple-game-go/player"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type StoredPassword struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	r := gin.Default()

	r.POST("/dashboard", middleware.JwtMiddleware(), player.DashboardHandler)

	r.POST("/register", player.RegisterHandler)

	r.POST("/login", auth.LoginHandler)
	r.POST("/logout", auth.LogoutHandler)
	r.POST("/add-bank", middleware.JwtMiddleware(), bank.RegisterBankAccountHandler)
	r.POST("/topup", middleware.JwtMiddleware(), bank.TopUpBalanceHandler)
	r.POST("/players", middleware.JwtMiddleware(), player.GetAllPlayersHandler)
	r.GET("/players/:username", middleware.JwtMiddleware(), player.GetPlayerWithWalletHandler)

	r.Run(":8080")

}
