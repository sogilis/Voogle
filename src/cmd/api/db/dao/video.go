package dao

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"

	"github.com/Sogilis/Voogle/src/cmd/api/models"
)

func CreateVideo(db *sql.DB, ID, title string, status int) (*models.Video, error) {
	query := "INSERT INTO videos (id, title, v_status) VALUES ( ? , ?, ?);"
	res, err := db.Exec(query, ID, title, status)
	if err != nil {
		log.Error("Error while insert into videos : ", err)
		return nil, err
	}

	nbRowAff, err := res.RowsAffected()
	if err != nil {
		log.Error("Error, can't know how many rows affected : ", err)
		return nil, err
	}

	log.Infof("%d row inserted", nbRowAff)
	return GetVideo(db, ID)
}

func GetVideo(db *sql.DB, ID string) (*models.Video, error) {
	query := `SELECT * FROM videos v WHERE v.id = ?`

	rows, err := db.Query(query, ID)
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
			&row.VideoStatus,
			&row.UploadedAt,
			&row.CreatedAt,
			&row.UpdatedAt,
		); err != nil {
			log.Error("Cannot read rows : ", err)
			return nil, err
		}
		videos = append(videos, row)
	}

	if len(videos) != 1 {
		err := fmt.Errorf("wrong number of results for unique id : %v in table videos", ID)
		log.Error(err)
		return nil, err
	}

	return &videos[0], nil

}

func GetVideos(db *sql.DB) ([]models.Video, error) {
	query := `SELECT *
			  FROM videos v`

	rows, err := db.Query(query)
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
			&row.VideoStatus,
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

func UpdateVideoStatus(db *sql.DB, ID string, status int) error {
	query := "UPDATE videos SET v_status = ? WHERE id = ?;"
	res, err := db.Exec(query, status, ID)
	if err != nil {
		log.Error("Error while update video status : ", err)
		return err
	}

	nbRowAff, err := res.RowsAffected()
	if err != nil {
		log.Error("Error, can't know how many rows affected : ", err)
		return err
	}

	log.Infof("%d row inserted", nbRowAff)
	return nil
}
