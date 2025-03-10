package storage

import (
	"awesomeProject3/lib/e"
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
)

var ErrNoSavedPages = errors.New("низя найти сохраненные страницы")

type Storage interface {
	Save(cts context.Context, p *Page) error
	PickRandom(ctx context.Context, UserName string) (*Page, error)
	Remove(ctx context.Context, p *Page) error
	IsExists(ctx context.Context, p *Page) (bool, error)
}
type Page struct {
	URL      string
	UserName string
}

func (p Page) Hash() (string, error) {
	h := sha1.New()

	if _, err := io.WriteString(h, p.URL); err != nil {
		return "", e.Wrap("низя хэшировать", err)
	}
	if _, err := io.WriteString(h, p.UserName); err != nil {
		return "", e.Wrap("низя хэшировать", err)
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
