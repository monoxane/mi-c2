package main

import (
	"log"
	"mi-c2/internal/api"
	"mi-c2/internal/controller"
	"mi-c2/internal/discord"
	"mi-c2/internal/env"
	"mi-c2/internal/logging"
	"os"
	"os/signal"
)

func main() {
	if !env.PopulateEnvironment() {
		os.Exit(1)
	}

	logging.Init()
	controller.Init()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")

	discord.Connect()
	go api.Listen(stop)

	<-stop

	discord.Cleanup()
}
