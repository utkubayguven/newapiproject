package main

import (
	"fmt"
	"log"
	"newapiprojet/database"
	"newapiprojet/docs"
	"newapiprojet/handlers"
	"newapiprojet/middlewares"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// .env dosyasını yükleyin
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

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

	// JWT gizli anahtarının ayarlandığını kontrol et
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		fmt.Println("JWT gizli anahtarı yapılandırılmamış")
		os.Exit(1)
	}
	fmt.Println("JWT Secret Key: ", secretKey) // Debugging için ekledik

	userRoutes := r.Group("/user")
	{
		userRoutes.POST("/register", h.Register)
		userRoutes.DELETE("/:id", h.DeleteUser)
		userRoutes.POST("/login", h.Login)
	}

	// Account routes
	accountRoutes := r.Group("/account")
	{
		accountRoutes.GET("/:id", h.GetAccountByID)
		accountRoutes.POST("/withdrawal", h.Withdrawal)
		accountRoutes.POST("/deposit", h.Deposit)
		accountRoutes.POST("/pin-change/:id", h.PinChange)
		accountRoutes.DELETE("/deleteacc/:accountNumber", h.DeleteAccount)
	}

	protected := r.Group("/account/token")
	protected.Use(middlewares.AuthenticateJWT())
	{
		protected.GET("/balance/:accountNumber", h.GetAccountBalance)
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.Run()
}
