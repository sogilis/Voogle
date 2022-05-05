package dao

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"

	"github.com/Sogilis/Voogle/src/cmd/api/models"
)

func CreateTableVideos(ctx context.Context, db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS videos (
		id              VARCHAR(36) NOT NULL,
		title           VARCHAR(64) NOT NULL,
		video_status    INT NOT NULL,
		uploaded_at     DATETIME,
		created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	
		CONSTRAINT pk PRIMARY KEY (id),
		CONSTRAINT unique_title UNIQUE (title)
	);`

	if _, err := db.ExecContext(ctx, query); err != nil {
		log.Error("Cannot create table : ", err)
		return err
	}

	log.Info("Table videos created (or existed already)")

	return nil
}

func CreateVideo(ctx context.Context, db *sql.DB, ID, title string, status int) (*models.Video, error) {
	query := "INSERT INTO videos (id, title, video_status) VALUES ( ? , ?, ?)"
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		log.Error("Cannot prepare statement : ", err)
		return nil, err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, ID, title, status)
	if err != nil {
		log.Error("Error while insert into videos : ", err)
		return nil, err
	}

	nbRowAff, err := res.RowsAffected()
	if err != nil {
		log.Error("Error, can't know how many rows affected : ", err)
		return nil, err
	}

	// Check if one and only one rows has been affected
	if nbRowAff != 1 {
		err := fmt.Errorf("wrong number of row affected (%d) while creating video id : %v", nbRowAff, ID)
		log.Error(err)
		return nil, err
	}

	return GetVideo(ctx, db, ID)
}

func UpdateVideo(ctx context.Context, db *sql.DB, video *models.Video) error {
	query := "UPDATE videos SET title = ?, video_status = ?, uploaded_at = ? WHERE id = ?"
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		log.Error("Cannot prepare statement : ", err)
		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, video.Title, video.Status, video.UploadedAt, video.ID)
	if err != nil {
		log.Error("Error while update video status : ", err)
		return err
	}

	nbRowAff, err := res.RowsAffected()
	if err != nil {
		log.Error("Error, can't know how many rows affected : ", err)
		return err
	}

	// Check if one and only one rows has been affected
	if nbRowAff != 1 {
		err := fmt.Errorf("wrong number of row affected (%d) while update id : %v in table videos", nbRowAff, video.ID)
		log.Error(err)
		return err
	}
	return nil
}

func UpdateVideoTx(ctx context.Context, tx *sql.Tx, video *models.Video) error {
	query := "UPDATE videos SET title = ?, video_status = ?, uploaded_at = ? WHERE id = ?"
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		log.Error("Cannot prepare statement : ", err)
		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, video.Title, video.Status, video.UploadedAt, video.ID)
	if err != nil {
		log.Error("Error while update video status : ", err)
		return err
	}

	nbRowAff, err := res.RowsAffected()
	if err != nil {
		log.Error("Error, can't know how many rows affected : ", err)
		return err
	}

	// Check if one and only one rows has been affected
	if nbRowAff != 1 {
		err := fmt.Errorf("wrong number of row affected (%d) while update id : %v in table videos", nbRowAff, video.ID)
		log.Error(err)
		return err
	}

	return nil
}

func GetVideo(ctx context.Context, db *sql.DB, ID string) (*models.Video, error) {
	query := "SELECT * FROM videos v WHERE v.id = ?"
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		log.Error("Cannot prepare statement : ", err)
		return nil, err
	}
	defer stmt.Close()

	var video models.Video
	err = stmt.QueryRowContext(ctx, ID).Scan(
		&video.ID,
		&video.Title,
		&video.Status,
		&video.UploadedAt,
		&video.CreatedAt,
		&video.UpdatedAt,
	)
	if err != nil {
		log.Error("Error, video not found : ", err)
		return nil, err
	}

	return &video, nil
}

func GetVideoFromTitle(ctx context.Context, db *sql.DB, title string) (*models.Video, error) {
	query := "SELECT * FROM videos v WHERE v.title = ?"
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		log.Error("Cannot prepare statement : ", err)
		return nil, err
	}
	defer stmt.Close()

	var video models.Video
	err = stmt.QueryRowContext(ctx, title).Scan(
		&video.ID,
		&video.Title,
		&video.Status,
		&video.UploadedAt,
		&video.CreatedAt,
		&video.UpdatedAt,
	)
	if err != nil {
		log.Error("Error, video not found : ", err)
		return nil, err
	}

	return &video, nil
}

func GetVideos(ctx context.Context, db *sql.DB) ([]models.Video, error) {
	query := "SELECT * FROM videos v"
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		log.Error("Cannot prepare statement : ", err)
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		log.Error("Error, cannot query database : ", err)
		return nil, err
	}

	defer func() {
		if err = rows.Close(); err != nil {
			log.Error("Error while closing database Rows", err)
		}
	}()

	var videos []models.Video
	for rows.Next() {
		var row models.Video
		if err := rows.Scan(
			&row.ID,
			&row.Title,
			&row.Status,
			&row.UploadedAt,
			&row.CreatedAt,
			&row.UpdatedAt,
		); err != nil {
			log.Error("Cannot read rows : ", err)
			return nil, err
		}
		videos = append(videos, row)
	}

	return videos, nil
}
