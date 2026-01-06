package utils

import (
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"os"
	"strconv"
	"strings"
	"sync"
)

const (
	RequestTimeout = 5
)

// log lvl
var (
	IterationDuration               = [2]int{150, 310}
	RequestDuration                 = [2]int{20, 19}
	telegramUserID                  int
	telegramBotToken                string
	MaxGoroutines                   int
	RetryCount                      int
	FilesPath                       string
	MainAccountCookie               string
	MainAccountAuthToken            string
	EnableUsernameChange            bool
	EnableTelegramMessageIfFound404 bool
	MinUsersRemainToNewIter         int
)

// todo исправить значения, ввести новые, убрать старые
// todo поправить файлы конфигов
var (
	AlreadyUser    int
	ExceptionsUser int
	githubAccounts []GitHubAccount
	Users          []User
	mutex          sync.Mutex
)

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal().Msg("Config file loading error")
	}

	RetryCount = getEnvAsInt("RETRY_COUNT", true)
	MaxGoroutines = getEnvAsInt("MAX_GOROUTINES", true)
	telegramUserID = getEnvAsInt("TG_USER_ID", true)
	FilesPath = getEnvAsString("FILES_PATH", true)
	MainAccountCookie = getEnvAsString("MAIN_ACCOUNT_COOKIE", true)
	MainAccountAuthToken = getEnvAsString("MAIN_ACCOUNT_AUTH_TOKEN", true)
	telegramBotToken = getEnvAsString("TG_BOT_TOKEN", true)
	EnableTelegramMessageIfFound404 = getEnvAsBool("ENABLE_TELEGRAM_MESSAGE_IF_404", true)
	EnableUsernameChange = getEnvAsBool("ENABLE_USERNAME_CHANGE", true)
}

func getEnvAsString(key string, required bool) string {
	value := os.Getenv(key)
	if value == "" && required {
		log.Fatal().Msgf("Error, var %s not set", key)
	}
	return strings.Trim(value, `"`)
}

func getEnvAsInt(key string, required bool) int {
	valueStr := getEnvAsString(key, required)
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Fatal().Msgf("Error converting %s to int: %v", key, err)
	}
	return value
}

func getEnvAsBool(key string, required bool) bool {
	valueStr := getEnvAsString(key, required)
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		log.Fatal().Msgf("Error converting %s to bool: %v", key, err)
	}
	return value
}
