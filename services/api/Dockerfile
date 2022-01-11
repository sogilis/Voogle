FROM golang:1.17.6-buster@sha256:ad61e922fb389384f00ba98dba8dae56f1dcf9ff731506857ab8f1312c47c212 as builder

WORKDIR /go/src/api
COPY . .

RUN go build -o build/api

FROM debian:11.2-slim@sha256:b0d53c872fd640c2af2608ba1e693cfc7dedea30abcd8f584b23d583ec6dadc7

WORKDIR /api
COPY --from=builder /go/src/api/build/api /api

CMD ["./api"]