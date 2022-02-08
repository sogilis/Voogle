# Service API
## Purpose

REST API called by the webapp. It serves the list of videos, each masters and segments of a video can be accessed, it
transfers uploaded video into a storage, and it emits events when a video processing can be done.

## Env vars

| Name          | Required   | Default value   | Description                                                        |
|---------------|------------|-----------------|--------------------------------------------------------------------|
| PORT          | false      | 4444            | Listening port of the API                                          |
| USER_AUTH     | true       | N/A             | Username (used by the webapp)                                      |
| PWD_AUTH      | true       | N/A             | User password (used by the webapp)                                 |
| DEV_MODE      | false      | false           | Enable debug logs                                                  |
| S3_HOST       | false      | ""              | Host address use by the S3 client (If empty, it connects to AWS)   |
| S3_AUTH_KEY   | true       | N/A             | S3 access token                                                    |
| S3_AUTH_PWD   | true       | N/A             | S3 password token                                                  |
| S3_BUCKET     | false      | voogle-video    | Bucket name used to store and access the videos                    |
| S3_REGION     | false      | eu-west-3       | Region used when the API connects to AWS                           |
