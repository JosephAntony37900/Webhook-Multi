package infrastructure

import (
	"encoding/json"
	"github/pull_request_webhook/application"
	"github/pull_request_webhook/domain/value_objects"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleDeployEvent(ctx *gin.Context) {
	eventType := ctx.GetHeader("X-GitHub-Event")
	deliveryD := ctx.GetHeader("X-GitHub-Delivery")

	log.Printf("Nuevo evento: %s con ID: %s", eventType, deliveryD)

	rawData, err := ctx.GetRawData()
	if err != nil {
		log.Printf("Error al leer el cuerpo de la solicitud: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "leer datos"})
		return
	}

	var statusCode int
	var message string

	switch eventType {
	case "deploy":
		statusCode, message = application.ProcessDeployEvent(rawData)
	default:
		ctx.JSON(http.StatusOK, gin.H{"success": "Normal"})
		return
	}

	var payload value_objects.DeployEvent
	if err := json.Unmarshal(rawData, &payload); err != nil {
		log.Printf("Error al deserializar el payload del despliegue: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error al procesar el payload del despliegue"})
		return
	}

	switch statusCode {
	case 200:
		ctx.JSON(http.StatusOK, gin.H{"success": "Despliegue procesado con éxito", "message": message})
	case 403:
		log.Printf("Error al deserializar el payload del despliegue: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error al procesar el payload del despliegue"})
	default:
		ctx.JSON(http.StatusOK, gin.H{"success": "Normal"})
	}
}