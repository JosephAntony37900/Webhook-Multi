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

func ProcessDeployEvent(rawData []byte) (int, string) {
	var eventPayload value_objects.DeployEvent

	if err := json.Unmarshal(rawData, &eventPayload); err != nil {
		log.Printf("Error al deserializar el payload: %v", err)
		return 403, "Error al deserializar el payload"
	}

	log.Printf("Evento de despliegue recibido con acción de %s", eventPayload.Action)

	var message string
	switch eventPayload.Status {
	case "on":
		if eventPayload.Success {
			message = formatDeploySuccessMessage(eventPayload)
		} else {
			message = formatDeployFailedMessage(eventPayload)
		}
	case "off":
		message = formatDeployOffMessage(eventPayload)
	default:
		message = "Estado de despliegue no manejado"
	}

	log.Printf("Mensaje a enviar a Discord: %s", message) 

	if err := sendToDiscord2(message); err != nil {
		log.Printf("Error al enviar el mensaje a Discord: %v", err)
		return 500, "Error al enviar el mensaje a Discord"
	}

	return 200, message
}

func sendToDiscord2(message string) error {
	webhookURL := os.Getenv("DISCORD_WEBHOOK_TEST") // Usar el webhook de pruebas
	if webhookURL == "" {
		return fmt.Errorf("DISCORD_WEBHOOK_TEST no está definido en el archivo .env")
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

func formatDeploySuccessMessage(event value_objects.DeployEvent) string {
	return fmt.Sprintf("✅ *Despliegue Exitoso*\n"+
		"📂 *Repositorio:* %s\n"+
		"👤 *Usuario:* %s\n"+
		"🟢 *Estado:* Encendido\n"+
		"📅 *Fecha:* %s\n"+
		"----------------------------------------------------",
		event.Repo.Name,
		event.Sender.Login,
		time.Now().Format("2006-01-02 15:04:05"))
}

func formatDeployFailedMessage(event value_objects.DeployEvent) string {
	return fmt.Sprintf("❌ *Despliegue Fallido*\n"+
		"📂 *Repositorio:* %s\n"+
		"👤 *Usuario:* %s\n"+
		"🔴 *Estado:* Encendido (pero falló)\n"+
		"📅 *Fecha:* %s\n"+
		"----------------------------------------------------",
		event.Repo.Name,
		event.Sender.Login,
		time.Now().Format("2006-01-02 15:04:05"))
}

func formatDeployOffMessage(event value_objects.DeployEvent) string {
	return fmt.Sprintf("⚠ *API Apagada*\n"+
		"📂 *Repositorio:* %s\n"+
		"👤 *Usuario:* %s\n"+
		"⚫ *Estado:* Apagado\n"+
		"📅 *Fecha:* %s\n"+
		"----------------------------------------------------",
		event.Repo.Name,
		event.Sender.Login,
		time.Now().Format("2006-01-02 15:04:05"))
}