package telegram

import (
	"awesomeProject3/lib/e"
	"awesomeProject3/storage"
	"context"
	"database/sql"
	"errors"
	"log"
	"net/url"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
	KittyCmd = "/SendKitty"
)

var db *sql.DB

func init() {
	db, err := sql.Open("sqlite3", "./kotiki.db")
	if err != nil {
		log.Fatal(err)
	}
}

func (p *Processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("получена новая команда '%s' из '%s'", text, username)

	if isAddCmd(text) {
		return p.savePage(chatID, text, username)
	}

	switch text {
	case RndCmd:
		return p.sendRandom(chatID, username)
	case HelpCmd:
		return p.sendHelp(chatID)
	case StartCmd:
		return p.sendHello(chatID)
	case KittyCmd:
		return p.sendKitty(chatID)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}
}
func (p *Processor) savePage(chatID int, pageURL string, username string) (err error) {
	defer func() {
		err = e.WrapIfErr("низя использовать команду: сохранить страницу", err)
	}()
	page := &storage.Page{
		URL:      pageURL,
		UserName: username,
	}

	isExist, err := p.storage.IsExists(context.Background(), page)
	if err != nil {
		return err
	}
	if isExist {
		return p.tg.SendMessage(chatID, msgAlreadyExists)
	}
	if err := p.storage.Save(context.Background(), page); err != nil {
		return nil
	}
	if err := p.tg.SendMessage(chatID, msgSaved); err != nil {
		return err
	}
	return nil
}
func (p *Processor) sendRandom(chatID int, username string) (err error) {

	page, err := p.storage.PickRandom(context.Background(), username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return err
	}
	if errors.Is(err, storage.ErrNoSavedPages) {
		return p.tg.SendMessage(chatID, msgNoSavedPages)
	}

	if err := p.tg.SendMessage(chatID, page.URL); err != nil {
		return err
	}
	return p.storage.Remove(context.Background(), page)
}
func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}
func (p *Processor) sendHello(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}
func (p *Processor) sendKitty(chatID int) (err error) {
	defer func() { err = e.WrapIfErr("низя отправлять картинку", err) }()

	rows, err := db.Query("SELECT data FROM kotiki ORDER BY RANDOM() LIMIT 1")
	if err != nil {
		return err
	}
	defer func() { _ = rows.Close() }()

	var photo []byte
	for rows.Next() {
		err = rows.Scan(&photo)
		if err != nil {
			return err
		}
	}
	// Отправить картинку в Telegram чат
	msg := p.tg.NewPhoto(tgbotapi.NewPhotoConfig{
		File: tgbotapi.FileBytes{
			Name:  "kitty.jpg",
			Bytes: photo,
		},
	})
	err = p.tg.SendMessage(chatID, msg)
	return err
}
func isAddCmd(text string) bool {
	return isURL(text)
}
func isURL(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}
