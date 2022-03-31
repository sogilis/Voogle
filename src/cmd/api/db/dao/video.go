package dao

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"

	"github.com/Sogilis/Voogle/src/cmd/api/db/models"
)

func CreateVideo(db *sql.DB, video *models.VideoUpload) error {
	query := "INSERT INTO videos (id, public_id, title) VALUES ( ? , ?, ?);"
	res, err := db.Exec(query, video.Id, video.PublicId, video.Title)
	if err != nil {
		log.Error("Error while insert into videos : ", err)
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

func GetVideos(db *sql.DB) ([]models.Video, error) {
	query := `SELECT v.id, public_id, title, state_name, last_update
			  FROM videos v
			  INNER JOIN video_state vs ON v.v_state = vs.id;`

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
		if err := rows.Scan(&row.Id, &row.PublicId, &row.Title, &row.VState, &row.LastUpdate); err != nil {
			log.Error("Cannot read rows : ", err)
			return nil, err
		}
		videos = append(videos, row)
	}

	return videos, nil
}
