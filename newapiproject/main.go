package main

import (
	"fmt"
	"log"
	"newapiprojet/config"
	"newapiprojet/database"
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

	// etcd endpoint'lerini belirleyin
	endpoints := []string{
		"http://etcd1:2379",
		"http://etcd2:2378",
		"http://etcd3:2377",
	}
	// her handlerın ıcıne clıentı cagır adaptor pattern
	//raft uygulaması
	// etcd client'ı başlatın
	etcdClient, err := database.InitEtcd(endpoints)
	if err != nil {
		log.Fatalf("Error initializing etcd client: %v", err)
	}

	if etcdClient == nil {
		log.Fatalf("etcdClient is nil after initialization")
	}

	// Config dosyasını yükleyin
	conf := config.GetConfig()

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	mw := middlewares.NewNewapiprojetMiddlewares()
	r.Use(mw.LogMiddleware())
	r.Use(middlewares.RequestLimitMiddleware())

	h := handlers.NewHandler(etcdClient.Client)

	// JWT gizli anahtarının ayarlandığını kontrol et
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		fmt.Println("JWT secret key is not configured")
		os.Exit(1)
	}
	fmt.Println("JWT Secret Key: ", secretKey) // Debugging için ekledik

	userRoutes := r.Group("/user")
	{
		userRoutes.POST("/register", h.Register)
		// userRoutes.POST("/login", h.Login)

	}

	// Account routes
	protected := r.Group("/account")
	protected.Use(middlewares.AuthenticateJWT())
	{
		// protected.GET("/balance/:accountNumber", h.GetAccountBalance)
		// protected.POST("/withdrawal", h.Withdrawal)
		// protected.POST("/deposit", h.Deposit)
		// protected.POST("/pin-change/:id", h.PinChange)
		// protected.DELETE("/deleteacc/:accountNumber", h.DeleteAccount)
	}

	// protected2 := r.Group("/user")
	// protected2.Use(middlewares.AuthenticateJWT())
	// {
	// 	protected2.DELETE("/:id", h.DeleteUser)
	// }

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.Run(fmt.Sprintf(":%d", conf.APIPort)) // API portunu config dosyasından alın
}
