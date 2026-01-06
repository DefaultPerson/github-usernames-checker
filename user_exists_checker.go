package main

import (
	"context"
	"github-checker/utils"
	"github.com/rs/zerolog/log"
	"math/rand"
	"sync"
	"time"
)

func queueLoader(qu chan<- utils.User) {
	log.Info().Msg("Queue loader running")
	for {
		startTime := time.Now()
		totalUsers := len(utils.Users)
		log.Info().Msgf("%d users loaded in queue", totalUsers)
		for _, user := range utils.Users {
			qu <- user
		}
		for utils.AlreadyUser+utils.MinUsersRemainToNewIter < totalUsers {
			time.Sleep(300 * time.Millisecond)
			if utils.AlreadyUser > 0 {
				elapsed := time.Since(startTime)
				avgTimePerUser := elapsed / time.Duration(utils.AlreadyUser)
				remainingUsers := totalUsers - utils.AlreadyUser
				estimatedTimeRemaining := time.Duration(remainingUsers) * avgTimePerUser
				excPercent := (float64(utils.ExceptionsUser) / float64(utils.AlreadyUser)) * 100
				log.Info().Msgf("READY/ALL:REMAIN %-5d/%d:%-5d(%-5d/%.2f%% EXC) - PER/REM/FACT: %v/%v/%v",
					utils.AlreadyUser, totalUsers, remainingUsers, utils.ExceptionsUser, excPercent, avgTimePerUser, estimatedTimeRemaining, elapsed)
			}
		}
		log.Info().Msgf("Iteration completed in %v", time.Since(startTime))
		utils.AlreadyUser = 0
		utils.ExceptionsUser = 0
	}
}

func runWorker(queue <-chan utils.User, ctx context.Context, pool *utils.TransportPool, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		for user := range queue {
			utils.CheckGitHubUserExists(ctx, pool, user)
			time.Sleep(time.Duration(rand.Intn(utils.IterationDuration[1])+utils.IterationDuration[0]) * time.Millisecond)
		}
	}
}

func main() {
	utils.GetLogger()
	utils.LoadConfig()
	ctx := context.Background()
	proxies := utils.ReadProxies()
	pool := utils.NewTransportPool(proxies)
	utils.GetAccountsTokens()
	utils.Users, _ = utils.GetUsersFromFile()
	queue := make(chan utils.User, len(utils.Users))
	go queueLoader(queue)
	var wg sync.WaitGroup
	for i := 0; i < utils.MaxGoroutines; i++ {
		wg.Add(1)
		go runWorker(queue, ctx, pool, &wg)
	}
	log.Info().Msg("Workers is running")
	wg.Wait()
}
