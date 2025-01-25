package loadPhotos

import (
	"database/sql"
	"io/ioutil"
	"log"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

const dirPhotos = "C:\\Users\\glowe\\kotyata"

func LoadPhotos(db *sql.DB) error {
	// Загрузить фотографии из папки
	files, err := ioutil.ReadDir(dirPhotos)
	if err != nil {
		return err
	}

	// Вставить фотографии в хранилище
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		// Открыть файл фотографии
		filePath := filepath.Join(dirPhotos, file.Name())
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			return err
		}

		// Вставить фотографию в хранилище
		_, err = db.Exec("INSERT INTO kotiki (data) VALUES (?)", data)
		if err != nil {
			return err
		}

		log.Printf("Загружена фотография: %s\n", file.Name())
	}

	return nil
}
