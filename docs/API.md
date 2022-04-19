# Webapp & API communication

Webapp and API communicate with JSON, because it's simple and efficient.

# GET - all video

Route: `GET /api/v1/videos`

The json will be:

```json
{
  "status": "Success/Fail",
  "data": [
    {
      "id": "",
      "title": "..."
    },
    {
      "id": "",
      "title": "..."
    }
  ]
}
```

Directory storage video: `api/videos`

# GET - video master

Route: `GET /api/v1/videos/{id}/streams/master.m3u8`

Binary stream of the master file content

# GET - video sub part

Route: `GET /api/v1/videos/{id}/streams/{quality}/{filename}`

Binary stream of the requested file content

# GET - list of videos

Route: `GET api/v1/videos/list`

Json list of videos

# POST - upload video

Route: `POST /api/v1/videos/upload`

Json video uploaded informations and usable links

The json will be:

```json
{
  "Video":{
    "ID":"a-unique-id",
    "Title":"a-title",
    "Status":"VIDEO_STATUS_UPLOADED",
    "UploadedAt":"2022-04-15T14:19:59.23123497+02:00",
    "CreatedAt":"2022-04-15T12:19:58Z",
    "UpdatedAt":"2022-04-15T12:19:58Z"},
    "Links":[
      {
        "rel":"Status",
        "href":"api/v1/videos/upload/a-unique-id/status",
        "method":"GET"
      },
      {
        "rel":"Stream",
        "href":"/api/v1/videos/a-unique-id/streams/master.m3u8",
        "method":"GET"
      }
    ]
}
```
# GET POST - metrics

Route: `GET /metrics`
Route: `POST /metrics`

Metrics for prometheus

# GET - video status

Route: `GET /api/v1/videos/{id}/status`

Json status of the requested video
