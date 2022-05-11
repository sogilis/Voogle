package dao

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"

	"github.com/Sogilis/Voogle/src/cmd/api/models"
)

func CreateTableUploads(ctx context.Context, db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS uploads (
		id              VARCHAR(36) NOT NULL,
		video_id        VARCHAR(36) NOT NULL,
		upload_status   INT NOT NULL,
		uploaded_at     DATETIME,
		created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	
		CONSTRAINT pk PRIMARY KEY (id),
		CONSTRAINT fk_v_id FOREIGN KEY (video_id) REFERENCES videos (id)
	);`

	if _, err := db.ExecContext(ctx, query); err != nil {
		log.Error("Cannot create table : ", err)
		return err
	}

	log.Info("Table uploads created (or existed already)")

	return nil
}

func CreateUpload(ctx context.Context, db *sql.DB, ID, videoID string, status int) (*models.Upload, error) {
	query := "INSERT INTO uploads (id, video_id, upload_status) VALUES ( ? , ?, ?)"
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		log.Error("Cannot prepare statement : ", err)
		return nil, err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, ID, videoID, status)
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
		err := fmt.Errorf("wrong number of row affected (%d) while creating upload id : %v", nbRowAff, ID)
		log.Error(err)
		return nil, err
	}

	return GetUpload(ctx, db, ID)
}

func DeleteUpload(ctx context.Context, db *sql.DB, ID string) error {
	query := "DELETE FROM uploads WHERE video_id = ?"
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		log.Error("Cannot prepare statement : ", err)
		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, ID)
	if err != nil {
		log.Error("Error while delete from uploads : ", err)
		return err
	}

	nbRowAff, err := res.RowsAffected()
	if err != nil {
		log.Error("Error, can't know how many rows affected : ", err)
		return err
	}

	// Check if at least one row has been affected
	if nbRowAff == 0 {
		err := fmt.Errorf("wrong number of row affected (%d) while deleting uploads with video id : %v", nbRowAff, ID)
		log.Error(err)
		return err
	}

	return nil
}

func DeleteUploadTx(ctx context.Context, tx *sql.Tx, ID string) error {
	query := "DELETE FROM uploads WHERE video_id = ?"
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		log.Error("Cannot prepare statement : ", err)
		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, ID)
	if err != nil {
		log.Error("Error while delete from uploads : ", err)
		return err
	}

	nbRowAff, err := res.RowsAffected()
	if err != nil {
		log.Error("Error, can't know how many rows affected : ", err)
		return err
	}

	// Check if at least one row has been affected
	if nbRowAff == 0 {
		err := fmt.Errorf("wrong number of row affected (%d) while deleting uploads with video id : %v", nbRowAff, ID)
		log.Error(err)
		return err
	}

	return nil
}

func UpdateUpload(ctx context.Context, db *sql.DB, upload *models.Upload) error {
	query := "UPDATE uploads SET video_id = ?, upload_status = ?, uploaded_at = ? WHERE id = ?"
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		log.Error("Cannot prepare statement : ", err)
		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, upload.VideoId, upload.Status, upload.UploadedAt, upload.ID)
	if err != nil {
		log.Error("Error while update upload : ", err)
		return err
	}

	nbRowAff, err := res.RowsAffected()
	if err != nil {
		log.Error("Error, can't know how many rows affected : ", err)
		return err
	}

	// Check if one and only one rows has been affected
	if nbRowAff != 1 {
		err := fmt.Errorf("wrong number of row affected (%d) while update id : %v in table uploads", nbRowAff, upload.ID)
		log.Error(err)
		return err
	}

	return nil
}

func UpdateUploadTx(ctx context.Context, tx *sql.Tx, upload *models.Upload) error {
	query := "UPDATE uploads SET video_id = ?, upload_status = ?, uploaded_at = ? WHERE id = ?"
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		log.Error("Cannot prepare statement : ", err)
		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, upload.VideoId, upload.Status, upload.UploadedAt, upload.ID)
	if err != nil {
		log.Error("Error while update upload : ", err)
		return err
	}

	nbRowAff, err := res.RowsAffected()
	if err != nil {
		log.Error("Error, can't know how many rows affected : ", err)
		return err
	}

	// Check if one and only one rows has been affected
	if nbRowAff != 1 {
		err := fmt.Errorf("wrong number of row affected (%d) while update id : %v in table uploads", nbRowAff, upload.ID)
		log.Error(err)
		return err
	}

	return nil
}

func GetUpload(ctx context.Context, db *sql.DB, id string) (*models.Upload, error) {
	query := "SELECT * FROM uploads u WHERE u.id = ?"
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		log.Error("Cannot prepare statement : ", err)
		return nil, err
	}
	defer stmt.Close()

	var upload models.Upload
	err = stmt.QueryRowContext(ctx, id).Scan(
		&upload.ID,
		&upload.VideoId,
		&upload.Status,
		&upload.UploadedAt,
		&upload.CreatedAt,
		&upload.UpdatedAt,
	)
	if err != nil {
		log.Error("Error, upload not found : ", err)
		return nil, err
	}

	return &upload, nil
}

func GetUploads(ctx context.Context, db *sql.DB) ([]models.Upload, error) {
	query := "SELECT * FROM uploads v"
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
