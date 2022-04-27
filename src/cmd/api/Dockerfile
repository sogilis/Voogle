FROM golang:1.17.9-bullseye@sha256:a4bd6256540bb5283b13dbc8aa9fdfd3440b009ce65771d293a97f1da6af5955 as builder

WORKDIR /go/src/voogle
COPY . .

RUN go build ./cmd/api

FROM debian:11.3-slim@sha256:b771c35d1e6ecf2556718ad3c0f481b4a04c1fbc133c609643acc9dd6743ead2

RUN apt-get update && apt-get install --no-install-recommends -y ca-certificates=20210119 && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /api
COPY --from=builder /go/src/voogle/api /api

CMD ["./api"]
