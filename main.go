package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/ququzone/verifying-paymaster-service/config"
)

func main() {
	conf := config.GetValues()

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

	if err := r.Run(fmt.Sprintf(":%d", conf.Port)); err != nil {
		log.Fatal(err)
	}
}
