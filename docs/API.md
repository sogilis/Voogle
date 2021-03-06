# Webapp & API communication

Webapp and API communicate with JSON, because it's simple and efficient.

# GET - all video

Route: `GET /api/v1/videos/list`

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

# POST - upload video

Route: `POST /api/v1/videos/upload`

Json video uploaded informations and usable links

The json will be:

```json
{
  "video":{
    "id":"a-unique-id",
    "title":"title",
    "status":"Encoding",
    "uploadedAt":"2022-04-22T12:01:13.619636641+02:00",
    "createdAt":"2022-04-22T10:01:12Z",
    "updatedAt":"2022-04-22T10:01:12Z"
  },
  "links":[
    {
      "rel":"status",
      "href":"/api/v1/videos/a-unique-id/status",
      "method":"get"
    },
    {
      "rel":"stream",
      "href":"/api/v1/videos/a-unique-id/streams/master.m3u8",
      "method":"get"
    }
  ]
}
```

# GET - video informations

Route: `GET /api/v1/videos/{id}/info`

Json video informations

The json will be:

```json
{
  "title": "title",
  "uploadDateUnix": "date",
}
```
# GET POST - metrics

Route: `GET /metrics`
Route: `POST /metrics`

Metrics for prometheus

# GET - video status

Route: `GET /api/v1/videos/{id}/status`

Json status of the requested video

# DELETE - video

Route: `DELETE /api/v1/videos/{id}/delete`

# GET - video cover

Route: `GET /api/v1/videos/{id}/cover`

PNG Video cover in base64

# GET - websocket

Route: `GET /ws`
