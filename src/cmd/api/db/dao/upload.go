package dao

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"

	"github.com/Sogilis/Voogle/src/cmd/api/models"
)

type UploadsRequestName int

const (
	CreateTableUploadsReq UploadsRequestName = iota
	CreateUpload
	UpdateUpload
	GetUpload
	GetUploads
	DeleteUpload
)

var UploadsRequests = map[UploadsRequestName]string{
	CreateTableUploadsReq: `CREATE TABLE IF NOT EXISTS uploads (
			id              VARCHAR(36) NOT NULL,
			video_id        VARCHAR(36) NOT NULL,
			upload_status   INT NOT NULL,
			uploaded_at     DATETIME,
			created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		
			CONSTRAINT pk PRIMARY KEY (id),
			CONSTRAINT fk_v_id FOREIGN KEY (video_id) REFERENCES videos (id)
		);`,

	CreateUpload: "INSERT INTO uploads (id, video_id, upload_status) VALUES ( ? , ?, ?)",
	UpdateUpload: "UPDATE uploads SET video_id = ?, upload_status = ?, uploaded_at = ? WHERE id = ?",
	GetUpload:    "SELECT * FROM uploads WHERE id = ?",
	GetUploads:   "SELECT * FROM uploads",
	DeleteUpload: "DELETE FROM uploads WHERE video_id = ?",
}

type UploadsDAO struct {
	DB               *sql.DB
	stmtCreateUpload *sql.Stmt
	stmtUpdateUpload *sql.Stmt
	stmtGetUpload    *sql.Stmt
	stmtGetUploads   *sql.Stmt
	stmtDeleteUpload *sql.Stmt
}

func CreateUploadsDAO(ctx context.Context, db *sql.DB) (*UploadsDAO, error) {
	if err := createTableUploads(ctx, db); err != nil {
		log.Error("Cannot create table uploads : ", err)
		return nil, err
	}

	uploadDAO, err := prepareUploadStmts(ctx, db)
	if err != nil {
		log.Error("Cannot prepare uploads statements : ", err)
		return nil, err
	}

	uploadDAO.DB = db

	return uploadDAO, nil
}

func (u UploadsDAO) CreateUpload(ctx context.Context, ID, videoID string, status int) (*models.Upload, error) {
	res, err := u.stmtCreateUpload.ExecContext(ctx, ID, videoID, status)
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

	return u.GetUpload(ctx, ID)
}

func (u UploadsDAO) DeleteUpload(ctx context.Context, ID string) error {
	res, err := u.stmtDeleteUpload.ExecContext(ctx, ID)
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

func (u UploadsDAO) DeleteUploadTx(ctx context.Context, tx *sql.Tx, ID string) error {
	stmt := tx.StmtContext(ctx, u.stmtDeleteUpload)
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

func (u UploadsDAO) UpdateUpload(ctx context.Context, upload *models.Upload) error {
	res, err := u.stmtUpdateUpload.ExecContext(ctx, upload.VideoId, upload.Status, upload.UploadedAt, upload.ID)
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

func (u UploadsDAO) UpdateUploadTx(ctx context.Context, tx *sql.Tx, upload *models.Upload) error {
	stmt := tx.StmtContext(ctx, u.stmtUpdateUpload)
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

func (u UploadsDAO) GetUpload(ctx context.Context, id string) (*models.Upload, error) {
	var upload models.Upload
	err := u.stmtGetUpload.QueryRowContext(ctx, id).Scan(
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

func (u UploadsDAO) GetUploads(ctx context.Context, db *sql.DB) ([]models.Upload, error) {
	rows, err := u.stmtGetUploads.QueryContext(ctx)
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

func createTableUploads(ctx context.Context, db *sql.DB) error {
	if _, err := db.ExecContext(ctx, UploadsRequests[CreateTableUploadsReq]); err != nil {
		log.Error("Cannot create table : ", err)
		return err
	}

	log.Info("Table uploads created (or existed already)")
	return nil
}

func prepareUploadStmts(ctx context.Context, db *sql.DB) (*UploadsDAO, error) {
	stmts := UploadsDAO{}

	// CreateUpload
	var err error
	stmts.stmtCreateUpload, err = db.PrepareContext(ctx, UploadsRequests[CreateUpload])
	if err != nil {
		log.Error("Cannot prepare statement : ", err)
		return nil, err
	}

	// UpdateUpload
	stmts.stmtUpdateUpload, err = db.PrepareContext(ctx, UploadsRequests[UpdateUpload])
	if err != nil {
		log.Error("Cannot prepare statement : ", err)
		return nil, err
	}

	// GetUpload
	stmts.stmtGetUpload, err = db.PrepareContext(ctx, UploadsRequests[GetUpload])
	if err != nil {
		log.Error("Cannot prepare statement : ", err)
		return nil, err
	}

	// GetUploads
	stmts.stmtGetUploads, err = db.PrepareContext(ctx, UploadsRequests[GetUploads])
	if err != nil {
		log.Error("Cannot prepare statement : ", err)
		return nil, err
	}

	// DeleteUpload
	stmts.stmtDeleteUpload, err = db.PrepareContext(ctx, UploadsRequests[DeleteUpload])
	if err != nil {
		log.Error("Cannot prepare statement : ", err)
		return nil, err
	}

	return &stmts, nil
}

func (u UploadsDAO) Close() {
	_ = u.stmtCreateUpload.Close()
	_ = u.stmtUpdateUpload.Close()
	_ = u.stmtGetUpload.Close()
	_ = u.stmtGetUploads.Close()
	_ = u.stmtDeleteUpload.Close()
}
