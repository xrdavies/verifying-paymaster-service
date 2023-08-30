package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/ququzone/verifying-paymaster-service/api"
	"github.com/ququzone/verifying-paymaster-service/config"
	"github.com/ququzone/verifying-paymaster-service/container"
	"github.com/ququzone/verifying-paymaster-service/db"
	"github.com/ququzone/verifying-paymaster-service/jsonrpc"
	"github.com/ququzone/verifying-paymaster-service/logger"
	"github.com/ququzone/verifying-paymaster-service/models"
)

func main() {
	err := logger.InitLogger()
	if err != nil {
		log.Fatalf("initial logger error: %v", err)
	}
	err = config.InitValues()
	if err != nil {
		log.Fatalf("init config error: %v", err)
	}

	repository := db.NewRepository()
	err = repository.AutoMigrate(&models.User{}, &models.ApiKeys{}, &models.Account{})
	if err != nil {
		logger.S().Fatalf("database migrate error: %v", err)
	}

	signerApi, err := api.NewSigner(container.NewContainer(repository))
	if err != nil {
		logger.S().Fatalf("instance signer error: %v", err)
	}

	conf := config.Config()
	gin.SetMode(conf.GinMode)
	r := gin.New()
	if err := r.SetTrustedProxies(nil); err != nil {
		logger.S().Fatalf("gin set trusted proxies error: %v", err)
	}
	r.Use(
		cors.Default(),
		gin.Recovery(),
	)
	r.GET("/ping", func(g *gin.Context) {
		g.String(http.StatusOK, "ok")
	})
	handlers := []gin.HandlerFunc{
		jsonrpc.Process(signerApi),
	}
	r.POST("/rpc/:key", handlers...)

	if err := r.Run(fmt.Sprintf(":%d", conf.Port)); err != nil {
		logger.S().Fatalf("gin run error: %v", err)
	}
}
