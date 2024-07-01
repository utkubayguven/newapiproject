package main

import (
	"fmt"
	"log"
	"newapiprojet/adapter"
	"newapiprojet/config"
	"newapiprojet/etcd"
	"newapiprojet/handlers"
	"newapiprojet/middlewares"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	client, err := etcd.NewEtcdClient(clientv3.Config{
		Endpoints: []string{
			"http://etcd1:2379",
			"http://etcd2:2378",
			"http://etcd3:2377",
		},
	})
	if err != nil {
		panic(err)
	}
	defer client.Close()

	etcdAdapter := adapter.NewEtcdAdapter(client)

	conf := config.GetConfig()

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	mw := middlewares.NewNewapiprojetMiddlewares()
	r.Use(mw.LogMiddleware())
	r.Use(middlewares.RequestLimitMiddleware())

	h := handlers.NewHandler(etcdAdapter)

	// User routes
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		fmt.Println("JWT secret key is not configured")
		os.Exit(1)
	}
	fmt.Println("JWT Secret Key: ", secretKey)

	userRoutes := r.Group("/user")
	{
		userRoutes.POST("/register", h.Register)
		userRoutes.POST("/login", h.Login)
	}

	// Account routes
	protected := r.Group("/account")
	protected.Use(middlewares.AuthenticateJWT())
	{
		protected.GET("/balance/:accountID", h.GetAccountBalance)
		protected.POST("/withdrawal", h.Withdrawal)
		protected.POST("/deposit", h.Deposit)
		protected.POST("/pin-change/:id", h.PinChange)
		protected.DELETE("/deleteacc/:id", h.DeleteAccountByID)
	}

	protected2 := r.Group("/user")
	protected2.Use(middlewares.AuthenticateJWT())
	{
		protected2.DELETE("delete/:id", h.DeleteUser)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.Run(fmt.Sprintf(":%d", conf.APIPort)) // listen and serve on
}
