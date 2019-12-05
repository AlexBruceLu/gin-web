package main

import (
	"gin-web/example/jaeger/sing/app/config"
	"gin-web/example/jaeger/sing/app/router"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(config.AppMode)

	engine := gin.New()
	router.SetupRouter(engine)

	log.Println("server Listen port", config.AppPort)

	if err := engine.Run(config.AppPort); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
