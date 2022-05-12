package dao

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"

	"github.com/Sogilis/Voogle/src/cmd/api/models"
)

type VideosRequestName int

const (
	createTableVideosReq VideosRequestName = iota
	CreateVideo
	UpdateVideo
	GetVideo
	GetVideoFromTitle
	GetTotalVideos
	DeleteVideo
)

var VideosRequests = map[VideosRequestName]string{
	createTableVideosReq: `CREATE TABLE IF NOT EXISTS videos (
			id              VARCHAR(36) NOT NULL,
			title           VARCHAR(64) NOT NULL,
			video_status    INT NOT NULL,
			uploaded_at     DATETIME,
			created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			source_path     VARCHAR(64) NOT NULL,

			CONSTRAINT pk PRIMARY KEY (id),
			CONSTRAINT unique_title UNIQUE (title)
		);`,

	CreateVideo:       "INSERT INTO videos (id, title, video_status, source_path) VALUES (?, ? , ?, ?)",
	UpdateVideo:       "UPDATE videos SET title = ?, video_status = ?, uploaded_at = ? WHERE id = ?",
	GetVideo:          "SELECT * FROM videos WHERE id = ?",
	GetVideoFromTitle: "SELECT * FROM videos WHERE title = ?",
	GetTotalVideos:    "SELECT COUNT(*) FROM videos",
	DeleteVideo:       "DELETE FROM videos WHERE id = ?",
}

type VideosDAO struct {
	DB                    *sql.DB
	stmtCreate            *sql.Stmt
	stmtUpdate            *sql.Stmt
	stmtGetVideo          *sql.Stmt
	stmtGetVideoFromTitle *sql.Stmt
	stmtGetTotalVideos    *sql.Stmt
	stmtDeleteVideo       *sql.Stmt
}

func CreateVideosDAO(ctx context.Context, db *sql.DB) (*VideosDAO, error) {
	if err := createTableVideos(ctx, db); err != nil {
		log.Error("Cannot create table videos : ", err)
		return nil, err
	}

	videoDAO, err := prepareVideoStmts(ctx, db)
	if err != nil {
		log.Error("Cannot prepare videos statements : ", err)
		return nil, err
	}

	videoDAO.DB = db

	return videoDAO, nil
}

func (v VideosDAO) CreateVideo(ctx context.Context, ID, title string, status int, sourcePath string) (*models.Video, error) {
	res, err := v.stmtCreate.ExecContext(ctx, ID, title, status, sourcePath)
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

	return v.GetVideo(ctx, ID)
}

func (v VideosDAO) DeleteVideo(ctx context.Context, ID string) error {
	res, err := v.stmtDeleteVideo.ExecContext(ctx, ID)
	if err != nil {
		log.Error("Error while delete from videos : ", err)
		return err
	}

	nbRowAff, err := res.RowsAffected()
	if err != nil {
		log.Error("Error, can't know how many rows affected : ", err)
		return err
	}

	// Check if one and only one rows has been affected
	if nbRowAff != 1 {
		err := fmt.Errorf("wrong number of row affected (%d) while deleting video id : %v", nbRowAff, ID)
		log.Error(err)
		return err
	}

	return nil
}

func (v VideosDAO) DeleteVideoTx(ctx context.Context, tx *sql.Tx, ID string) error {
	stmt := tx.StmtContext(ctx, v.stmtDeleteVideo)
	res, err := stmt.ExecContext(ctx, ID)
	if err != nil {
		log.Error("Error while delete from videos : ", err)
		return err
	}

	nbRowAff, err := res.RowsAffected()
	if err != nil {
		log.Error("Error, can't know how many rows affected : ", err)
		return err
	}

	// Check if one and only one rows has been affected
	if nbRowAff != 1 {
		err := fmt.Errorf("wrong number of row affected (%d) while deleting video id : %v", nbRowAff, ID)
		log.Error(err)
		return err
	}

	return nil
}

func (v VideosDAO) UpdateVideo(ctx context.Context, video *models.Video) error {
	res, err := v.stmtUpdate.ExecContext(ctx, video.Title, video.Status, video.UploadedAt, video.ID)
	if err != nil {
		log.Error("Error while update video : ", err)
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

func (v VideosDAO) UpdateVideoTx(ctx context.Context, tx *sql.Tx, video *models.Video) error {
	stmt := tx.StmtContext(ctx, v.stmtUpdate)
	res, err := stmt.ExecContext(ctx, video.Title, video.Status, video.UploadedAt, video.ID)
	if err != nil {
		log.Error("Error while update video : ", err)
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

func (v VideosDAO) GetVideo(ctx context.Context, ID string) (*models.Video, error) {
	var video models.Video
	err := v.stmtGetVideo.QueryRowContext(ctx, ID).Scan(
		&video.ID,
		&video.Title,
		&video.Status,
		&video.UploadedAt,
		&video.CreatedAt,
		&video.UpdatedAt,
		&video.SourcePath,
	)
	if err != nil {
		log.Error("Error, video not found : ", err)
		return nil, err
	}

	return &video, nil
}

func (v VideosDAO) GetVideoFromTitle(ctx context.Context, title string) (*models.Video, error) {
	var video models.Video
	err := v.stmtGetVideoFromTitle.QueryRowContext(ctx, title).Scan(
		&video.ID,
		&video.Title,
		&video.Status,
		&video.UploadedAt,
		&video.CreatedAt,
		&video.UpdatedAt,
		&video.SourcePath,
	)
	if err != nil {
		log.Error("Error, video not found : ", err)
		return nil, err
	}

	return &video, nil
}

func (v VideosDAO) GetVideos(ctx context.Context, attribute interface{}, ascending bool, page int, limit int) ([]models.Video, error) {

	switch attribute {
	case models.TITLE:
		attribute = "title"
	case models.UPLOADEDAT:
		attribute = "uploaded_at"
	case models.CREATEDAT:
		attribute = "created_at"
	case models.UPDATEDAT:
		attribute = "updated_at"
	default:
		err := fmt.Errorf("no such attribute")
		return nil, err
	}

	direction := "DESC"
	if ascending {
		direction = "ASC"
	}

	query := fmt.Sprintf("SELECT * FROM videos ORDER BY %v %v LIMIT %d,%d", attribute, direction, (page-1)*limit, limit)
	rows, err := v.DB.QueryContext(ctx, query)

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
			&row.SourcePath,
		); err != nil {
			log.Error("Cannot read rows : ", err)
			return nil, err
		}
		videos = append(videos, row)
	}

	return videos, nil
}

func (v VideosDAO) GetTotalVideos(ctx context.Context) (int, error) {
	var total int
	err := v.stmtGetTotalVideos.QueryRowContext(ctx).Scan(&total)
	if err != nil {
		log.Error("Cannot read rows : ", err)
		return -1, err
	}
	return total, nil
}

func prepareVideoStmts(ctx context.Context, db *sql.DB) (*VideosDAO, error) {
	stmts := VideosDAO{}

	// CreateVideo
	var err error
	stmts.stmtCreate, err = db.PrepareContext(ctx, VideosRequests[CreateVideo])
	if err != nil {
		log.Error("Cannot prepare statement : ", err)
		return nil, err
	}

	// UpdateVideo
	stmts.stmtUpdate, err = db.PrepareContext(ctx, VideosRequests[UpdateVideo])
	if err != nil {
		log.Error("Cannot prepare statement : ", err)
		return nil, err
	}

	// GetVideo
	stmts.stmtGetVideo, err = db.PrepareContext(ctx, VideosRequests[GetVideo])
	if err != nil {
		log.Error("Cannot prepare statement : ", err)
		return nil, err
	}

	// GetVideoFromTitle
	stmts.stmtGetVideoFromTitle, err = db.PrepareContext(ctx, VideosRequests[GetVideoFromTitle])
	if err != nil {
		log.Error("Cannot prepare statement : ", err)
		return nil, err
	}

	// GetTotalVideos
	stmts.stmtGetTotalVideos, err = db.PrepareContext(ctx, VideosRequests[GetTotalVideos])
	if err != nil {
		log.Error("Cannot prepare statement : ", err)
		return nil, err
	}

	// DeleteVideo
	stmts.stmtDeleteVideo, err = db.PrepareContext(ctx, VideosRequests[DeleteVideo])
	if err != nil {
		log.Error("Cannot prepare statement : ", err)
		return nil, err
	}

	return &stmts, nil
}

func createTableVideos(ctx context.Context, db *sql.DB) error {
	if _, err := db.ExecContext(ctx, VideosRequests[createTableVideosReq]); err != nil {
		log.Error("Cannot create table : ", err)
		return err
	}

	log.Info("Table videos created (or existed already)")
	return nil
}

func (v VideosDAO) Close() {
	_ = v.stmtCreate.Close()
	_ = v.stmtUpdate.Close()
	_ = v.stmtGetVideo.Close()
	_ = v.stmtGetVideoFromTitle.Close()
	_ = v.stmtDeleteVideo.Close()
	_ = v.stmtGetTotalVideos.Close()

}
