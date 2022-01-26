# Voogle

Voogle is an application for broadcasting and sharing video streams, it's purpose is to be the demonstration medium for the SquareScale platform.

## Needed tools

- Docker
- Makefile
- Act `0.2.25` [Github action local](https://github.com/golangci/golangci-lint-action)

- Go `1.17`
- Golangci-lint `1.29`

- Node `16.13.1`
- Npm `8.3.0`

## How to run the environment locally
To start Voogle on your machine, you need three services (for now): webapp, api and a S3-like.

/!\ If you want to use the local minio, you have to create the bucket with the UI first. No options exist to create it at launch with an env var.

You don't have to set manually `S3_HOST` unless you know what you are doing.

* You can start a MinIO that is a service that have the same API as S3, with `make start_s3_local`. 
The API will be available on the port `9000` and the console one the port `9001`. And it can be accessed with the credentials `admin` - `password` by default.
* Then you can start the api (`/services/api`) with `make run-dev` that uses the same credentials that the MinIO and the webapp by default.
* Finally, you can start the webapp (`/services/webapp`) with `npm run serve` to start the VueJS development server