FROM golang:1.18.2-bullseye@sha256:a95776d414fbb293ca9095c2b616cba2d684120d7f22061fb8f4845bd273fae6 as builder

WORKDIR /go/src/voogle
COPY . .

RUN go build ./cmd/gray-server-transformer

FROM debian:11.3-slim@sha256:b771c35d1e6ecf2556718ad3c0f481b4a04c1fbc133c609643acc9dd6743ead2

RUN apt-get update && apt-get install --no-install-recommends -y ca-certificates=20210119 ffmpeg=7:4.3.4-0+deb11u1 && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /gray-server-transformer
COPY --from=builder /go/src/voogle/gray-server-transformer /gray-server-transformer

EXPOSE 50051 50051

CMD ["./gray-server-transformer"]
