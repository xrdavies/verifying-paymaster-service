package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/ququzone/verifying-paymaster-service/config"
	"github.com/ququzone/verifying-paymaster-service/jsonrpc"
	"github.com/ququzone/verifying-paymaster-service/signer"
)

func main() {
	conf := config.GetValues()

	signerApi, err := signer.NewSigner()
	if err != nil {
		log.Fatalf("instance signer error: %v", err)
	}

	gin.SetMode(conf.GinMode)
	r := gin.New()
	if err := r.SetTrustedProxies(nil); err != nil {
		log.Fatal(err)
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
	r.POST("/", handlers...)

	if err := r.Run(fmt.Sprintf(":%d", conf.Port)); err != nil {
		log.Fatal(err)
	}
}
