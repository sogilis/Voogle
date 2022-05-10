package dao

import (
	"regexp"

	"github.com/DATA-DOG/go-sqlmock"
)

func ExpectVideosDAOCreation(mock sqlmock.Sqlmock) {
	queryCreateTAble := regexp.QuoteMeta(
		`CREATE TABLE IF NOT EXISTS videos (
			id              VARCHAR(36) NOT NULL,
			title           VARCHAR(64) NOT NULL,
			video_status    INT NOT NULL,
			uploaded_at     DATETIME,
			created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			source_path     VARCHAR(64) NOT NULL,

			CONSTRAINT pk PRIMARY KEY (id),
			CONSTRAINT unique_title UNIQUE (title)
		);`)

	stmtCreate := regexp.QuoteMeta("INSERT INTO videos (id, title, video_status) VALUES ( ? , ?, ?)")
	stmtUpdate := regexp.QuoteMeta("UPDATE videos SET title = ?, video_status = ?, uploaded_at = ? WHERE id = ?")
	stmtGetVideo := regexp.QuoteMeta("SELECT * FROM videos WHERE id = ?")
	stmtGetVideoFromTitle := regexp.QuoteMeta("SELECT * FROM videos WHERE title = ?")
	stmtGetVideos := regexp.QuoteMeta("SELECT * FROM videos")

	mock.ExpectExec(queryCreateTAble).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectPrepare(stmtCreate)
	mock.ExpectPrepare(stmtUpdate)
	mock.ExpectPrepare(stmtGetVideo)
	mock.ExpectPrepare(stmtGetVideoFromTitle)
	mock.ExpectPrepare(stmtGetVideos)
}

func ExpectUplaodsDAOCreation(mock sqlmock.Sqlmock) {
	queryCreateTAble := regexp.QuoteMeta(
		`CREATE TABLE IF NOT EXISTS uploads (
			id              VARCHAR(36) NOT NULL,
			video_id        VARCHAR(36) NOT NULL,
			upload_status   INT NOT NULL,
			uploaded_at     DATETIME,
			created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		
			CONSTRAINT pk PRIMARY KEY (id),
			CONSTRAINT fk_v_id FOREIGN KEY (video_id) REFERENCES videos (id)
		);`)

	stmtCreate := regexp.QuoteMeta("INSERT INTO uploads (id, video_id, upload_status) VALUES ( ? , ?, ?)")
	stmtUpdate := regexp.QuoteMeta("UPDATE uploads SET video_id = ?, upload_status = ?, uploaded_at = ? WHERE id = ?")
	stmtGetUpload := regexp.QuoteMeta("SELECT * FROM uploads WHERE id = ?")
	stmtGetUploads := regexp.QuoteMeta("SELECT * FROM uploads")

	mock.ExpectExec(queryCreateTAble).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectPrepare(stmtCreate)
	mock.ExpectPrepare(stmtUpdate)
	mock.ExpectPrepare(stmtGetUpload)
	mock.ExpectPrepare(stmtGetUploads)
}
