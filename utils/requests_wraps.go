package utils

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
)

//todo Спросить LLM как я могу оптимизировать код.

func changeAccountUsername(user User) {
	accountCredentials := popAccountCredentials()
	if accountCredentials != nil && EnableUsernameChange {
		if result := makeChangeUsernameRequest(accountCredentials, user.Username); result == 200 {
			// todo тут проверка, включена отправка или нет

			log.Info().Msgf("Success rename https://github.com/%s - Email:%s - Profile:%s - File:%s - Code:%d", user.Username, accountCredentials.Email, accountCredentials.Username, user.Filename, result)
			message := fmt.Sprintf("✅ Success rename\nURL: https://github.com/%s\nEmail:%s\nProfile:%s\nFile:%s\nCode:%d", user.Username, accountCredentials.Email, accountCredentials.Username, user.Filename, result)
			go SendTelegramMessage(message)
			return
		}
		log.Warn().Msgf("Error by rename https://github.com/%s - Email:%s - Profile:%s- File:%s", user.Username, accountCredentials.Email, accountCredentials.Username, user.Filename)
		message := fmt.Sprintf("❌ Error by rename github\nURL: https://github.com/%s\nEmail:%s\nProfile:%s\nFile:%s", user.Username, accountCredentials.Email, accountCredentials.Username, user.Filename)
		// todo тут проверка, включена отправка или нет

		go SendTelegramMessage(message)
	} else if !EnableUsernameChange {
		log.Warn().Msgf("Rename disable, do it! https://github.com/%s - File:%s", user.Username, user.Filename)
		message := fmt.Sprintf("❌ Rename disable, do it!\nURL: https://github.com/%s\nFile:%s", user.Username, user.Filename)
		// todo тут проверка, включена отправка или нет
		go SendTelegramMessage(message)
	} else {
		log.Warn().Msgf("No accounts available for https://github.com/%s - File:%s", user.Username, user.Filename)
		message := fmt.Sprintf("❌ Error by rename github account, no accounts available for\nURL: https://github.com/%s\nFile:%s", user.Username, user.Filename)
		go SendTelegramMessage(message)
	}
}

func checkUsernameAvailable(ctx context.Context, pool *TransportPool, user User) {
	// todo удалить юзера из utils.Users в конце
	//if visitedURL[thisSite] {
	//	return
	//}
	//if alreadyChecked := checkUsernameInFile(user.Username); alreadyChecked {
	//
	//	return
	//}

	if EnableTelegramMessageIfFound404 {
		message := fmt.Sprintf("🔁 Found new 404 user\nURL: https://github.com/%s\nFile:%s", user.Username, user.Filename)
		// todo убрать нахуй, оставить только лог!

		go SendTelegramMessage(message)
	}
	log.Info().Msgf("Found new 404 user - URL: https://github.com/%s - File:%s", user.Username, user.Filename)

	if result := makeUsernameAvailableRequest(ctx, pool, user.Username); result {
		log.Info().Msgf("Found availble user - https://github.com/%s", user.Username)
		changeAccountUsername(user)
	}
	writeUsernameInFile(user.Username)
}

func popAccountCredentials() *GitHubAccount {
	mutex.Lock()
	defer mutex.Unlock()

	if len(githubAccounts) == 0 {
		return nil
	}
	account := githubAccounts[0]
	githubAccounts = githubAccounts[1:]
	return &account
}

func CheckGitHubUserExists(ctx context.Context, pool *TransportPool, user User) {
	result := makeGitHubRequest(ctx, user.Username, pool)
	AlreadyUser++
	if result == 200 {
	} else if result == 404 {
		go checkUsernameAvailable(ctx, pool, user)
	} else {
		ExceptionsUser++
		log.Debug().Msgf("User request return with code: %d", result)
	}
	return
}
