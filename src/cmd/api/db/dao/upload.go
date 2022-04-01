package dao

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"

	"github.com/Sogilis/Voogle/src/cmd/api/models"
)

func CreateUpload(db *sql.DB, ID, videoID string, status int) (*models.Upload, error) {
	query := "INSERT INTO uploads (id, video_id, upload_status) VALUES ( ? , ?, ?)"
	res, err := db.Exec(query, ID, videoID, status)
	if err != nil {
		log.Error("Error while insert into uploads : ", err)
		return nil, err
	}

	nbRowAff, err := res.RowsAffected()
	if err != nil {
		log.Error("Error, can't know how many rows affected : ", err)
		return nil, err
	}

	// Check if one and only one rows has been affected
	if nbRowAff != 1 {
		err := fmt.Errorf("wrong number of row affected while creating upload id : %v", ID)
		log.Error(err)
		return nil, err
	}

	log.Debugf("%d row inserted", nbRowAff)
	return GetUpload(db, ID)
}

func UpdateUpload(db *sql.DB, upload *models.Upload) error {
	query := "UPDATE uploads SET video_id = ?, upload_status = ?, uploaded_at = ?, updated_at = ? WHERE id = ?"
	res, err := db.Exec(query, upload.VideoId, upload.Status, upload.UploadedAt, upload.UpdatedAt, upload.ID)
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
		err := fmt.Errorf("wrong number of row affected while update id : %v in table uploads", upload.ID)
		log.Error(err)
		return err
	}

	log.Debugf("%d row updated", nbRowAff)
	return nil
}

func GetUpload(db *sql.DB, id string) (*models.Upload, error) {
	query := "SELECT * FROM uploads u WHERE u.id = ?"

	rows, err := db.Query(query, id)
	if err != nil {
		log.Error("Error, cannot query database : ", err)
		return nil, err
	}

	defer func() {
		if err = rows.Close(); err != nil {
			log.Error("Error while closing database Rows", err)
		}
	}()

	var uploads []models.Upload
	for rows.Next() {
		var row models.Upload
		if err := rows.Scan(
			&row.ID,
			&row.VideoId,
			&row.Status,
			&row.UploadedAt,
			&row.CreatedAt,
			&row.UpdatedAt,
		); err != nil {
			log.Error("Cannot read rows : ", err)
			return nil, err
		}
		uploads = append(uploads, row)
	}

	if len(uploads) != 1 {
		err := fmt.Errorf("wrong number of results for unique id : %v in table uploads", id)
		log.Error(err)
		return nil, err
	}

	return &uploads[0], nil
}

func GetUploads(db *sql.DB) ([]models.Upload, error) {
	query := "SELECT * FROM uploads v"

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

	var uploads []models.Upload
	for rows.Next() {
		var row models.Upload
		if err := rows.Scan(
			&row.ID,
			&row.VideoId,
			&row.Status,
			&row.UploadedAt,
			&row.CreatedAt,
			&row.UpdatedAt,
		); err != nil {
			log.Error("Cannot read rows : ", err)
			return nil, err
		}
		uploads = append(uploads, row)
	}

	return uploads, nil
}
