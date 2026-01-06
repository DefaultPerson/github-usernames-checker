package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
	"path/filepath"
	"strings"
)

func writeUsernameInFile(username string) {
	file, err := os.OpenFile(fmt.Sprintf("%s/users_cache", FilesPath), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Error().Msgf("Unable to open or create user cache file: %v", err)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Error().Msgf("Unable to close user cache file: %v", err)
		}
	}(file)

	writer := bufio.NewWriter(file)
	_, err = writer.WriteString(username + "\n")
	if err != nil {
		log.Error().Msgf("Unable to write username to user cache file: %v", err)
		return
	}

	err = writer.Flush()
	if err != nil {
		log.Error().Msgf("Unable to flush writer: %v", err)
	}
}

func ReadProxies() []string {
	data, err := os.ReadFile(fmt.Sprintf("%s/proxy.txt", FilesPath))
	if err != nil {
		log.Fatal().Msgf("Error loading proxy file: %v", err)
	}
	log.Info().Msg("Proxy has been loaded")
	return strings.Fields(string(data))
}

func GetUsersFromFile() ([]User, error) {
	directoryPath := fmt.Sprintf("%s/users_set", FilesPath)
	dirEntries, err := os.ReadDir(directoryPath)
	if err != nil {
		log.Fatal().Msgf("Error reading directory: %v", err)
		return nil, err
	}

	var allUsers []User
	userSet := make(map[string]struct{})

	for _, entry := range dirEntries {
		if entry.IsDir() {
			continue
		}

		filePath := filepath.Join(directoryPath, entry.Name())
		file, err := os.Open(filePath)
		if err != nil {
			log.Fatal().Msgf("Error opening users file: %v", err)
			return nil, err
		}

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			user := User{
				Username: scanner.Text(),
				Filename: entry.Name(),
			}
			allUsers = append(allUsers, user)
			userSet[user.Username] = struct{}{}
		}

		if err := scanner.Err(); err != nil {
			log.Fatal().Msgf("Error scanning users file: %v", err)
			err := file.Close()
			if err != nil {
				log.Fatal().Msgf("Error scanning users file: %v", err)
			}
			return nil, err
		}

		if err := file.Close(); err != nil {
			log.Error().Msgf("Error closing users file: %v", err)
		}
	}

	cacheFilePath := fmt.Sprintf("%s/users_cache", FilesPath)
	cacheFile, err := os.Open(cacheFilePath)
	if err != nil {
		log.Fatal().Msgf("Error opening users cache file: %v", err)
		return nil, err
	}
	defer func() {
		if err := cacheFile.Close(); err != nil {
			log.Error().Msgf("Error closing users cache file: %v", err)
		}
	}()

	cacheUsers := make(map[string]struct{})
	scanner := bufio.NewScanner(cacheFile)
	for scanner.Scan() {
		cacheUsers[scanner.Text()] = struct{}{}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal().Msgf("Error scanning users cache file: %v", err)
		return nil, err
	}

	var filteredUsers []User
	for _, user := range allUsers {
		if _, found := cacheUsers[user.Username]; !found {
			filteredUsers = append(filteredUsers, user)
		}
	}

	return filteredUsers, nil
}

func GetAccountsTokens() {
	data, err := os.ReadFile(fmt.Sprintf("%s/github-accounts.json", FilesPath))
	if err != nil {
		log.Fatal().Msgf("Load accounts data error: %s", err)
	}
	err = json.Unmarshal(data, &githubAccounts)
	if err != nil {
		log.Fatal().Msgf("JSON decode error: %s", err)
	}
	log.Info().Msg("Accounts tokens has been loaded")
}
