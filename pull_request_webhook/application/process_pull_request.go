package application

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github/pull_request_webhook/domain/value_objects"
	"log"
	"net/http"
	"os"
	"time"
)

func ProcessPullRequestEvent(rawData []byte) (int, string) {
	var eventPayload value_objects.PullRequestEvent

	if err := json.Unmarshal(rawData, &eventPayload); err != nil {
		return 403, "Error al deserializar el payload"
	}

	log.Printf("Evento pull request recibido con acción de %s", eventPayload.Action)

	var message string
	switch eventPayload.Action {
	case "reopened":
		message = formatReopenedMessage(eventPayload)
	case "ready_for_review":
		message = formatReadyForReviewMessage(eventPayload)
	case "closed":
		message = formatClosedMessage(eventPayload)
	case "merged":
		message = formatMergedMessage(eventPayload)
	default:
		message = "Acción no manejada"
	}

	// Enviar el mensaje a Discord
	if err := sendToDiscord(message); err != nil {
		log.Printf("Error al enviar el mensaje a Discord: %v", err)
		return 500, "Error al enviar el mensaje a Discord"
	}

	return 200, message
}

func sendToDiscord(message string) error {
	webhookURL := os.Getenv("DISCORD_WEBHOOK")
	if webhookURL == "" {
		return fmt.Errorf("DISCORD_WEBHOOK no está definido en el archivo .env")
	}

	payload := map[string]string{
		"content": message,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error al serializar el payload: %v", err)
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("error al enviar el mensaje a Discord: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("respuesta no exitosa de Discord: %s", resp.Status)
	}

	return nil
}

func formatReopenedMessage(event value_objects.PullRequestEvent) string {
	return fmt.Sprintf("🔓 **Pull Request Reabierto**\n"+
		"📂 **Repositorio:** %s\n"+
		"👤 **Usuario:** %s\n"+
		"🌿 **Desde:** %s\n"+
		"🌳 **Hacia:** %s\n"+
		"📅 **Fecha:** %s\n"+
		"🔗 **URL:** %s\n" +
		"----------------------------------------------------",
		event.Repository.Name,
		event.PullRequest.User.Login,
		event.PullRequest.Head.Ref,
		event.PullRequest.Base.Ref,
		time.Now().Format("2006-01-02 15:04:05"),
		event.PullRequest.URL)
}

func formatReadyForReviewMessage(event value_objects.PullRequestEvent) string {
	return fmt.Sprintf("👀 **Pull Request Listo para Revisión**\n"+
		"📂 **Repositorio:** %s\n"+
		"👤 **Usuario:** %s\n"+
		"🌿 **Desde:** %s\n"+
		"🌳 **Hacia:** %s\n"+
		"📅 **Fecha:** %s\n"+
		"🔗 **URL:** %s\n"+
		"----------------------------------------------------",
		event.Repository.Name,
		event.PullRequest.User.Login,
		event.PullRequest.Head.Ref,
		event.PullRequest.Base.Ref,
		time.Now().Format("2006-01-02 15:04:05"),
		event.PullRequest.URL)
}

func formatClosedMessage(event value_objects.PullRequestEvent) string {
	return fmt.Sprintf("🚫 **Pull Request Cerrado**\n"+
		"📂 **Repositorio:** %s\n"+
		"👤 **Usuario:** %s\n"+
		"🌿 **Desde:** %s\n"+
		"🌳 **Hacia:** %s\n"+
		"📅 **Fecha:** %s\n"+
		"🔗 **URL:** %s\n"+
		"----------------------------------------------------",
		event.Repository.Name,
		event.PullRequest.User.Login,
		event.PullRequest.Head.Ref,
		event.PullRequest.Base.Ref,
		time.Now().Format("2006-01-02 15:04:05"),
		event.PullRequest.URL)
}

func formatMergedMessage(event value_objects.PullRequestEvent) string {
	return fmt.Sprintf("✅ **Pull Request Mergeado**\n"+
		"📂 **Repositorio:** %s\n"+
		"👤 **Usuario:** %s\n"+
		"🌿 **Desde:** %s\n"+
		"🌳 **Hacia:** %s\n"+
		"📅 **Fecha:** %s\n"+
		"🔗 **URL:** %s\n"+
		"----------------------------------------------------",
		event.Repository.Name,
		event.PullRequest.User.Login,
		event.PullRequest.Head.Ref,
		event.PullRequest.Base.Ref,
		time.Now().Format("2006-01-02 15:04:05"),
		event.PullRequest.URL)
}