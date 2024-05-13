package main

import (
	"fmt"
	"newapiprojet/database"
	"newapiprojet/docs"
	"newapiprojet/handlers"
	"newapiprojet/middlewares"
	"os"

	"github.com/gin-gonic/gin"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	db, err := database.InitDb()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	mw := middlewares.NewNewapiprojetMiddlewares()
	r.Use(mw.LogMiddleware())

	docs.SwaggerInfo.BasePath = ""

	h := handlers.NewHandler(db)

	r.POST("/register", h.Register)
	r.POST("/login", h.Login)
	r.GET("/balance/:accountNumber", h.GetAccountBalance)
	r.POST("pin-change/:id", h.PinChange)
	r.DELETE("deleteacc/:accountNumber", h.DeleteAccount)
	r.DELETE("deleteuser/:id", h.DeleteUser)
	r.POST("/withdrawal", h.Withdrawal)
	r.POST("/deposit", h.Deposit)
	r.GET("/account/:id", h.GetAccountByID)

	protected := r.Group("/api")
	protected.Use(middlewares.AuthenticateJWT())
	// protected.GET("/balance", h.BalanceInquiry)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.Run()
}
