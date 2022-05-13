package dao_test

import (
	"regexp"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Sogilis/Voogle/src/cmd/api/db/dao"
)

func ExpectVideosDAOCreation(mock sqlmock.Sqlmock) {
	mock.ExpectExec(regexp.QuoteMeta(dao.VideosRequests[dao.CreateTableVideosReq])).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectPrepare(regexp.QuoteMeta(dao.VideosRequests[dao.CreateVideo]))
	mock.ExpectPrepare(regexp.QuoteMeta(dao.VideosRequests[dao.UpdateVideo]))
	mock.ExpectPrepare(regexp.QuoteMeta(dao.VideosRequests[dao.GetVideo]))
	mock.ExpectPrepare(regexp.QuoteMeta(dao.VideosRequests[dao.GetVideoFromTitle]))
	mock.ExpectPrepare(regexp.QuoteMeta(dao.VideosRequests[dao.GetTotalVideos]))
	mock.ExpectPrepare(regexp.QuoteMeta(dao.VideosRequests[dao.DeleteVideo]))
}

func ExpectUploadsDAOCreation(mock sqlmock.Sqlmock) {
	mock.ExpectExec(regexp.QuoteMeta(dao.UploadsRequests[dao.CreateTableUploadsReq])).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectPrepare(regexp.QuoteMeta(dao.UploadsRequests[dao.CreateUpload]))
	mock.ExpectPrepare(regexp.QuoteMeta(dao.UploadsRequests[dao.UpdateUpload]))
	mock.ExpectPrepare(regexp.QuoteMeta(dao.UploadsRequests[dao.GetUpload]))
	mock.ExpectPrepare(regexp.QuoteMeta(dao.UploadsRequests[dao.GetUploads]))
	mock.ExpectPrepare(regexp.QuoteMeta(dao.UploadsRequests[(dao.DeleteUpload)]))
}
