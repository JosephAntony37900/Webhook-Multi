package infrastructure

import (
	"encoding/json"
	"github/pull_request_webhook/application"
	"github/pull_request_webhook/domain/value_objects"
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
)

func HandlePullRequestEvent(ctx *gin.Context) {
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
	case "ping":
		ctx.JSON(http.StatusOK, gin.H{"status": "success"})
		return
	case "pull_request":
		statusCode, message = application.ProcessPullRequestEvent(rawData)
	default:
		ctx.JSON(http.StatusOK, gin.H{"success": "Normal"})
		return
	}

	var payload value_objects.PullRequestEvent
	if err := json.Unmarshal(rawData, &payload); err != nil {
		log.Printf("Error al deserializar el payload del pull request: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error al procesar el payload del pull request"})
		return
	}

	if payload.Action == "closed" {
		log.Printf("Pull request cerrado")
		log.Println("Repo", payload.Repository)
		log.Println("usuario", payload.PullRequest.User.Login)
		log.Println("desde", payload.PullRequest.Head.Ref)
		log.Println("hacia", payload.PullRequest.Base.Ref)
	}

	switch statusCode {
	case 200:
		ctx.JSON(http.StatusOK, gin.H{"success": "Pull Request procesado con éxito", "message": message})
	case 403:
		log.Printf("Error al deserializar el payload del pull request: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error al procesar el payload del pull request"})
	default:
		ctx.JSON(http.StatusOK, gin.H{"success": "Normal"})
	}
}