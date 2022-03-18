package dao

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"

	"github.com/Sogilis/Voogle/src/cmd/api/db/models"
)

func PutVideo(db *sql.DB, video models.VideoModelUpload) error {
	query := "INSERT INTO videos (id, title) VALUES ('" + video.Title + "');"
	res, err := db.Exec(query)
	if err != nil {
		log.Error("Error while insert into videos:", err)
		return err
	}

	nbRowAff, err := res.RowsAffected()
	if err != nil {
		log.Error("Error, can't know how many rows affected:", err)
		return err
	}

	log.Infof("%d row(s) inserted", nbRowAff)
	return nil
}

func GetVideos(db *sql.DB) ([]models.VideoModel, error) {
	query := `SELECT v.id, client_id, title, vs.state_name, last_update
			  FROM videos v
			  INNER JOIN video_state vs ON v.v_state = vs.id;`

	rows, err := db.Query(query)
	if err != nil {
		log.Error("Error : ", err)
		return nil, err
	}

	defer func() {
		if err = rows.Close(); err != nil {
			log.Error("Error while closing database Rows", err)
		}
	}()

	var videos []models.VideoModel
	for rows.Next() {
		var row models.VideoModel
		if err := rows.Scan(&row.Id, &row.ClientId, &row.Title, &row.VState, &row.LastUpdate); err != nil {
			log.Error("Cannot read rows : ", err)
			return nil, err
		}
		videos = append(videos, row)
	}

	return videos, nil
}
