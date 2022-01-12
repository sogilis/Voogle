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
      "title": "..."
    },
    {
      "title": "..."
    }
  ]
}
```

