package db

import (
	"appContract/pkg/db"
	"appContract/pkg/models"
	"context"
	"errors"

	"log"

	"github.com/jackc/pgx"
)

func DBgetPhoto(id_user int) (*models.Photo, error) {
	conn := db.GetDB()
	if conn == nil {
		return nil, errors.New("failed to connect to database")
	}

	var photo models.Photo
	// Добавляем LIMIT 1 для гарантии одной записи
	err := conn.QueryRow(context.Background(),
		`SELECT id_photo, data, type 
         FROM user_photos 
         WHERE id_user = $1 
         LIMIT 1`, // Важно для безопасности
		id_user,
	).Scan(&photo.Id_photo, &photo.Data_photo, &photo.Type_photo)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // Фото не найдено
		}
		log.Printf("Error getting photo: %v", err)
		return nil, err
	}

	return &photo, nil
}
func DBaddPhoto(photo models.Photo) error {
	conn := db.GetDB()
	if conn == nil {
		return errors.New("failed to connect to database")
	}

	// Начинаем транзакцию
	tx, err := conn.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	// Добавляем новое фото
	_, err = tx.Exec(context.Background(),
		`INSERT INTO user_photos (data, type, id_user)
		VALUES ($1, $2, $3)`,
		photo.Data_photo,
		photo.Type_photo,
		photo.Id_user)
	if err != nil {
		log.Printf("Error inserting new photo: %v", err)
		return err
	}

	// Коммитим транзакцию
	if err = tx.Commit(context.Background()); err != nil {
		return err
	}

	return nil
}

func DBChangePhoto(photo models.Photo) error {
	conn := db.GetDB()
	if conn == nil {
		return errors.New("failed to connect to database")
	}

	// Начинаем транзакцию
	tx, err := conn.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	// Удаляем старые фото пользователя
	_, err = tx.Exec(context.Background(),
		"DELETE FROM user_photos WHERE id_user = $1",
		photo.Id_user)
	if err != nil {
		log.Printf("Error deleting old photos: %v", err)
		return err
	}

	// Добавляем новое фото
	_, err = tx.Exec(context.Background(),
		`INSERT INTO user_photos (data, type, id_user)
        VALUES ($1, $2, $3)`,
		photo.Data_photo,
		photo.Type_photo,
		photo.Id_user)
	if err != nil {
		log.Printf("Error inserting new photo: %v", err)
		return err
	}

	// Коммитим транзакцию
	if err = tx.Commit(context.Background()); err != nil {
		return err
	}

	return nil
}

func DBDeletePhoto(id_user int) error {
	conn := db.GetDB()
	if conn == nil {
		return errors.New("failed to connect to database")
	}

	// Начинаем транзакцию
	tx, err := conn.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	// Удаляем старые фото пользователя
	_, err = tx.Exec(context.Background(),
		"DELETE FROM user_photos WHERE id_user = $1",
		id_user)
	if err != nil {
		log.Printf("Error deleting old photos: %v", err)
		return err
	}

	// Коммитим транзакцию
	if err = tx.Commit(context.Background()); err != nil {
		return err
	}

	return nil
}
