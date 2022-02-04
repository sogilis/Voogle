# Voogle

Voogle is an application for broadcasting and sharing video streams, it's purpose is to be the demonstration medium for the SquareScale platform.

## Needed tools

- Docker
- Makefile
- Act `0.2.25` [Github action local](https://github.com/nektos/act)

- Go `1.17`
- Golangci-lint `1.29`(https://github.com/golangci/golangci-lint-action)

- Node `16.13.1`
- Npm `8.3.0`

## How to run the environment locally

To start Voogle on your machine, you need three services (for now): webapp, api and a S3-like and Redis.

/!\ If you want to use the local minio, you have to create the bucket with the UI first and it should be named `voogle-video`. No options exist to create it at launch with an env var.

You don't have to set manually `S3_HOST` unless you know what you are doing.

- You can start a MinIO that is a service that have the same API as S3, with `make start_s3`.
  The API will be available on the port `9000` and the console one the port `9001`. And it can be accessed with the credentials `admin` - `password` by default.
- You can start a Redis, with `make start_redis`. The Redis server will be available on the port `6379`. And it wan be accessed with password empty by default (yes, it's a strong password)
- Then you can start the api (`/services/api`) with `make run-dev` that uses the same credentials that the MinIO and the webapp by default.
- Finally, you can start the webapp (`/services/webapp`) with `npm run serve` to start the VueJS development server
- Credentials for Voogle account can be found in file (`/services/api/Makefile`) in `run-dev` or `run-dev-remote` command as USER_AUTH and PWD_AUTH environment variables.

## Visual Studio Code

### Configuration for multi-module workspaces

- `.vscode/setting.json`
  ```json
  {
    "go.useLanguageServer": true,
    "gopls": {
      "experimentalWorkspaceModule": true
    }
  }
  ```
