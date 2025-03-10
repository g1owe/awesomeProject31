package main

import (
	tgClient "awesomeProject3/clients/telegram"
	"awesomeProject3/consumer/event-consumer"
	"awesomeProject3/events/telegram"
	"awesomeProject3/storage/SQLite"
	"context"
	"database/sql"
	"flag"
	"log"
)

const (
	tgBotHost         = "api.telegram.org"
	sqliteStoragePath = "C:\\Users\\glowe\\database.db"
	batchSize         = 100
)

func main() {
	db, err := sql.Open("sqlite3", "storage.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	//s := files.New(storagePath)
	s, err := sqlite.New(sqliteStoragePath)
	if err != nil {
		log.Fatal("низя создать хранилище", err)
	}
	if err := s.Init(context.TODO()); err != nil {
		log.Fatal("низя инициализировать хранилище", err)
	}

	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, mustToken()),
		s,
	)

	log.Print("началась возня")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("ботик сдох(((", err)
	}

	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table' AND name='kotiki'")
	if err != nil {
		log.Println(err)
		return
	}
	defer rows.Close()

	if !rows.Next() {
		log.Println("Таблица storage не существует")
		return
	}
}
func mustToken() string {
	token := flag.String(
		"tg-bot-token",
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
