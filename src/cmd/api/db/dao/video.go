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
	CreateTableVideosReq VideosRequestName = iota
	CreateVideo
	UpdateVideo
	GetVideo
	GetVideoFromTitle
	GetVideosTitleAsc
	GetVideosTitleDesc
	GetVideosUploadedAtAsc
	GetVideosUploadedAtDesc
	GetTotalVideos
	DeleteVideo
)

var VideosRequests = map[VideosRequestName]string{
	CreateTableVideosReq: `CREATE TABLE IF NOT EXISTS videos (
			id              VARCHAR(36) NOT NULL,
			title           VARCHAR(64) NOT NULL,
			video_status    INT NOT NULL,
			uploaded_at     DATETIME,
			created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			source_path     VARCHAR(64) NOT NULL,
			cover_path      VARCHAR(64),

			CONSTRAINT pk PRIMARY KEY (id),
			CONSTRAINT unique_title UNIQUE (title)
		);`,

	CreateVideo:             "INSERT INTO videos (id, title, video_status, source_path, cover_path) VALUES (?, ? , ?, ?, ?)",
	UpdateVideo:             "UPDATE videos SET title = ?, video_status = ?, uploaded_at = ? WHERE id = ?",
	GetVideo:                "SELECT * FROM videos WHERE id = ?",
	GetVideoFromTitle:       "SELECT * FROM videos WHERE title = ?",
	GetVideosTitleAsc:       "SELECT * FROM videos WHERE video_status = ? ORDER BY title ASC LIMIT ?,?",
	GetVideosTitleDesc:      "SELECT * FROM videos WHERE video_status = ? ORDER BY title DESC LIMIT ?,?",
	GetVideosUploadedAtAsc:  "SELECT * FROM videos WHERE video_status = ? ORDER BY uploaded_at ASC LIMIT ?,?",
	GetVideosUploadedAtDesc: "SELECT * FROM videos WHERE video_status = ? ORDER BY uploaded_at DESC LIMIT ?,?",
	GetTotalVideos:          "SELECT COUNT(*) FROM videos WHERE video_status = ?",
	DeleteVideo:             "DELETE FROM videos WHERE id = ?",
}

type VideosDAO struct {
	DB                          *sql.DB
	stmtCreate                  *sql.Stmt
	stmtUpdate                  *sql.Stmt
	stmtGetVideo                *sql.Stmt
	stmtGetVideoFromTitle       *sql.Stmt
	stmtGetVideosTitleAsc       *sql.Stmt
	stmtGetVideosTitleDesc      *sql.Stmt
	stmtGetVideosUploadedAtAsc  *sql.Stmt
	stmtGetVideosUploadedAtDesc *sql.Stmt
	stmtGetTotalVideos          *sql.Stmt
	stmtDeleteVideo             *sql.Stmt
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

	// GetVideosTitleAsc
	stmts.stmtGetVideosTitleAsc, err = db.PrepareContext(ctx, VideosRequests[GetVideosTitleAsc])
	if err != nil {
		log.Error("Cannot prepare statement : ", err)
		return nil, err
	}

	// GetVideosTitleDesc
	stmts.stmtGetVideosTitleDesc, err = db.PrepareContext(ctx, VideosRequests[GetVideosTitleDesc])
	if err != nil {
		log.Error("Cannot prepare statement : ", err)
		return nil, err
	}

	// GetVideosUploadedAtAsc
	stmts.stmtGetVideosUploadedAtAsc, err = db.PrepareContext(ctx, VideosRequests[GetVideosUploadedAtAsc])
	if err != nil {
		log.Error("Cannot prepare statement : ", err)
		return nil, err
	}

	// GetVideosUploadedAtDesc
	stmts.stmtGetVideosUploadedAtDesc, err = db.PrepareContext(ctx, VideosRequests[GetVideosUploadedAtDesc])
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
	if _, err := db.ExecContext(ctx, VideosRequests[CreateTableVideosReq]); err != nil {
		log.Error("Cannot create table : ", err)
		return err
	}

	log.Debug("Table videos created (or existed already)")
	return nil
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

func (v VideosDAO) CreateVideo(ctx context.Context, ID, title string, status int, sourcePath string, coverPath string) (*models.Video, error) {
	res, err := v.stmtCreate.ExecContext(ctx, ID, title, status, sourcePath, coverPath)
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
		&video.CoverPath,
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
		&video.CoverPath,
	)
	if err != nil {
		log.Error("Error, video not found : ", err)
		return nil, err
	}

	return &video, nil
}

func (v VideosDAO) GetVideos(ctx context.Context, attribute interface{}, ascending bool, page, limit, status int) ([]models.Video, error) {

	var stmt *sql.Stmt
	switch attribute {
	case models.TITLE:
		if ascending {
			stmt = v.stmtGetVideosTitleAsc
		} else {
			stmt = v.stmtGetVideosTitleDesc
		}

	case models.UPLOADEDAT:
		if ascending {
			stmt = v.stmtGetVideosUploadedAtAsc
		} else {
			stmt = v.stmtGetVideosUploadedAtDesc
		}

	case models.CREATEDAT:
		err := fmt.Errorf("Request for create date not yet implemented")
		return nil, err

	case models.UPDATEDAT:
		err := fmt.Errorf("Request for update date not yet implemented")
		return nil, err

	default:
		err := fmt.Errorf("no such attribute")
		return nil, err
	}

	rows, err := stmt.QueryContext(ctx, status, (page-1)*limit, limit)
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
			&row.CoverPath,
		); err != nil {
			log.Error("Cannot read rows : ", err)
			return nil, err
		}
		videos = append(videos, row)
	}

	return videos, nil
}

func (v VideosDAO) GetTotalVideos(ctx context.Context, status int) (int, error) {
	var total int
	err := v.stmtGetTotalVideos.QueryRowContext(ctx, status).Scan(&total)
	if err != nil {
		log.Error("Cannot read rows : ", err)
		return -1, err
	}
	return total, nil
}

func (v VideosDAO) Close() {
	_ = v.stmtCreate.Close()
	_ = v.stmtUpdate.Close()
	_ = v.stmtGetVideo.Close()
	_ = v.stmtGetVideoFromTitle.Close()
	_ = v.stmtDeleteVideo.Close()
	_ = v.stmtGetTotalVideos.Close()
	_ = v.stmtGetVideosTitleAsc.Close()
	_ = v.stmtGetVideosTitleDesc.Close()
	_ = v.stmtGetVideosUploadedAtAsc.Close()
	_ = v.stmtGetVideosUploadedAtDesc.Close()

}
