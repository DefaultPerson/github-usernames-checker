package utils

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

// todo для проверки используется первый аккаунт
// todo подобрать и вынести тайминги
// todo добавить по-больше логов
//todo Спросить LLM как я могу оптимизировать код.
// todo английский и оптимизация
// todo добавить по-больше логов
//todo английский в логах
//todo Спросить LLM как я могу оптимизировать код.
// подумать где нужен дебаг

func makeGitHubRequest(ctx context.Context, username string, pool *TransportPool) int {
	// todo retry counter
	// todo  api.
	//задаржки
	requestURL := "https://github.com/" + username

	for i := 0; i < RetryCount; i++ {
		client := HttpClientPool.Get().(*http.Client)
		client.Transport = pool.GetRandomTransport()

		req, err := http.NewRequestWithContext(ctx, "GET", requestURL, nil)
		if err != nil {
			HttpClientPool.Put(client)
			continue
		}

		resp, err := client.Do(req)
		if err != nil {
			HttpClientPool.Put(client)
			time.Sleep(time.Duration(rand.Intn(370)+220) * time.Millisecond)
			//todo
			continue
		}

		defer func(Body io.ReadCloser) {
			if err := Body.Close(); err != nil {
			}
		}(resp.Body)

		HttpClientPool.Put(client)

		if resp.StatusCode == 200 || resp.StatusCode == 404 {
			return resp.StatusCode
		}
	}

	return 500
}

func makeChangeUsernameRequest(accountCredentials *GitHubAccount, username string) int {
	//todo пул не используется
	//todo исправить headers
	ghUrl := fmt.Sprintf("https://github.com/users/%s/rename", accountCredentials.Username)
	method := "POST"

	payload := strings.NewReader(fmt.Sprintf("authenticity_token=%s&login=%s&timestamp=%s&timestamp_secret=%s", accountCredentials.AuthToken, username, accountCredentials.Timestamp, accountCredentials.TimestampSecret))

	client := &http.Client{}
	req, _ := http.NewRequest(method, ghUrl, payload)
	req.Header.Add("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/jxl,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Add("accept-language", "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Add("cache-control", "max-age=0")
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	req.Header.Add("cookie", accountCredentials.Cookie)
	req.Header.Add("dnt", "1")
	req.Header.Add("origin", "https://github.com")
	req.Header.Add("referer", "https://github.com/settings/admin")
	req.Header.Add("sec-ch-ua", "\"Chromium\";v=\"123\", \"Not:A-Brand\";v=\"8\"")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-ch-ua-platform", "\"Linux\"")
	req.Header.Add("sec-fetch-dest", "document")
	req.Header.Add("sec-fetch-mode", "navigate")
	req.Header.Add("sec-fetch-site", "same-origin")
	req.Header.Add("sec-fetch-user", "?1")
	req.Header.Add("upgrade-insecure-requests", "1")
	req.Header.Add("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36")
	for {
		res, err := client.Do(req)
		if err != nil {
			log.Warn().Msgf("Request failed: %e", err)
			time.Sleep(1 * time.Second)
			continue
		}
		bodyBytes, _ := io.ReadAll(res.Body)
		body := string(bodyBytes)
		func() {
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {
					log.Warn().Msgf("Error when close body: %v", err)
				}
			}(res.Body)
		}()

		if res.StatusCode == 200 {
			return 200
		} else if res.StatusCode == 422 {
			log.Error().Msgf("Request failed 422: %s", body)
			return 422
		}
		time.Sleep(1 * time.Second)
	}
}

func makeUsernameAvailableRequest(ctx context.Context, pool *TransportPool, username string) bool {
	//todo пул не используется
	//todo исправить headers
	client := HttpClientPool.Get().(*http.Client)

	ghUrl := "https://github.com/account/rename_check?suggest_usernames=true"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("authenticity_token", MainAccountAuthToken)
	_ = writer.WriteField("value", username)
	err := writer.Close()
	if err != nil {
		log.Warn().Msgf("Error when close writer: %v", err)
	}

	req, _ := http.NewRequest(method, ghUrl, payload)
	req.Header.Add("accept", "*/*")
	req.Header.Add("accept-language", "en-US,en;q=0.9,ru;q=0.8")
	req.Header.Add("cookie", MainAccountCookie)
	req.Header.Add("dnt", "1")
	req.Header.Add("origin", "https://github.com")
	req.Header.Add("referer", "https://github.com/settings/admin")
	req.Header.Add("sec-ch-ua", "\"Chromium\";v=\"123\", \"Not:A-Brand\";v=\"8\"")
	req.Header.Add("sec-ch-ua-mobile", "?0")
	req.Header.Add("sec-ch-ua-platform", "\"Linux\"")
	req.Header.Add("sec-fetch-dest", "empty")
	req.Header.Add("sec-fetch-mode", "cors")
	req.Header.Add("sec-fetch-site", "same-origin")
	req.Header.Add("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36")
	req.Header.Set("Content-Type", writer.FormDataContentType())

	for {
		resp, err := client.Do(req)
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}
		bodyBytes, _ := io.ReadAll(resp.Body)
		err = resp.Body.Close()
		if err != nil {
			log.Warn().Msgf("Error when close body: %v", err)

		}
		body := string(bodyBytes)

		if resp.StatusCode == 422 && (strings.Contains(body, "is not available") || strings.Contains(body, "is unavailable") || strings.Contains(body, "only contain alphanumeric characters")) {
			return false
		}
		if strings.Contains(body, "is available") && resp.StatusCode == 200 {
			return true
		}
		message := fmt.Sprintf("Error by check! github account https://github.com/%s with code %d.\nRelaunch script", username, resp.StatusCode)
		log.Warn().Msg(message)
		// todo добавить опцию включения тг сообщений
		go SendTelegramMessage(message)
		time.Sleep(2 * time.Second)
	}
}
