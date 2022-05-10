package dao

import (
	"regexp"

	"github.com/DATA-DOG/go-sqlmock"
)

func ExpectVideosDAOCreation(mock sqlmock.Sqlmock) {
	mock.ExpectExec(regexp.QuoteMeta(VideosRequests[createTableVideosReq])).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectPrepare(regexp.QuoteMeta(VideosRequests[CreateVideo]))
	mock.ExpectPrepare(regexp.QuoteMeta(VideosRequests[UpdateVideo]))
	mock.ExpectPrepare(regexp.QuoteMeta(VideosRequests[GetVideo]))
	mock.ExpectPrepare(regexp.QuoteMeta(VideosRequests[GetVideoFromTitle]))
	mock.ExpectPrepare(regexp.QuoteMeta(VideosRequests[GetVideos]))
}

func ExpectUplaodsDAOCreation(mock sqlmock.Sqlmock) {
	mock.ExpectExec(regexp.QuoteMeta(UploadsRequests[createTableUploadsReq])).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectPrepare(regexp.QuoteMeta(UploadsRequests[CreateUpload]))
	mock.ExpectPrepare(regexp.QuoteMeta(UploadsRequests[UpdateUpload]))
	mock.ExpectPrepare(regexp.QuoteMeta(UploadsRequests[GetUpload]))
	mock.ExpectPrepare(regexp.QuoteMeta(UploadsRequests[GetUploads]))
}
