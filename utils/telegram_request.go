package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
	"time"
)

func SendTelegramMessage(message string) {
	currentTime := time.Now()
	message += fmt.Sprintf("\n%02d:%02d:%02d.%02d", currentTime.Hour(), currentTime.Minute(), currentTime.Second(), currentTime.Nanosecond()/1e7)
	client := &http.Client{
		Timeout: RequestTimeout * time.Second,
	}
	tgUrl := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", telegramBotToken)
	data := map[string]string{
		"chat_id": strconv.Itoa(telegramUserID),
		"text":    message,
	}
	dataBytes, _ := json.Marshal(data)
	for {
		req, _ := http.NewRequest("POST", tgUrl, bytes.NewBuffer(dataBytes))
		req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			log.Warn().Msgf("An error occurred while trying to send a message to the telegram bot with error %e", err)
			time.Sleep(1 * time.Second)
			continue
		}
		if resp.StatusCode == 200 {
			log.Info().Msgf("Send a message to the Telegram bot using code %d", resp.StatusCode)
			return
		}
		log.Warn().Msgf("An error occurred while trying to send a message to the telegram bot with code %d", resp.StatusCode)
		time.Sleep(1 * time.Second)
	}
}
