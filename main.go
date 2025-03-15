package main

import (
	"github/pull_request_webhook/infrastructure"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	// Cargar variables de entorno desde el archivo .env
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error cargando el archivo .env: %v", err)
	}

	engine := gin.Default()

	infrastructure.Routes(engine)

	engine.Run(":4000")
}