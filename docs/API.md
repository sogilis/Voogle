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