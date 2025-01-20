package main

import (
	tgClient "awesomeProject3/clients/telegram"
	event_consumer "awesomeProject3/consumer/event-consumer"
	"awesomeProject3/events/telegram"
	"awesomeProject3/storage/files"
	"flag"
	"log"
)

const (
	tgBotHost   = "api.telegram.org"
	storagePath = "storage"
	batchSize   = 10
)

func main() {

	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, mustToken()),
		files.New(storagePath),
	)

	log.Print("началась возня")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("ботик сдох(((", err)
	}
}
func mustToken() string {
	token := flag.String(
		"tg1-bot-token",
		"",
		"token nice to tg bot ",
	)
	flag.Parse()

	if *token == "" {
		log.Fatal("token111")
	}
	return *token
}

// 8074614402:AAHPDXmki0H0aW-QhRkHcJtENjRjU8Q5eTs
