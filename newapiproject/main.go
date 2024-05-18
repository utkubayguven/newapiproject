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
		userRoutes.POST("/login", h.Login)
	}

	// Account routes

	protected := r.Group("/account")
	protected.Use(middlewares.AuthenticateJWT())
	{
		protected.GET("/balance/:accountNumber", h.GetAccountBalance)
		protected.POST("/withdrawal", h.Withdrawal) // with json body parameter
		protected.POST("/deposit", h.Deposit)       // with json body parameter
		protected.POST("/pin-change/:id", h.PinChange)
		protected.DELETE("/deleteacc/:accountNumber", h.DeleteAccount)
	}

	protected2 := r.Group("/user")
	protected2.Use(middlewares.AuthenticateJWT())
	{
		protected2.DELETE("/:id", h.DeleteUser)
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.Run()
}
