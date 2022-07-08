FROM golang:1.18.2-bullseye@sha256:a95776d414fbb293ca9095c2b616cba2d684120d7f22061fb8f4845bd273fae6 as builder

# WORKDIR /go/src/voogle
# COPY . .

# RUN go build ./cmd/api

# 1. Precompile the entire go standard library into the first Docker cache layer: useful for other projects too!
RUN CGO_ENABLED=0 GOOS=linux go install -v -installsuffix cgo -a std

# 2. Prepare and enter the src folder
WORKDIR /go/src/voogle

# 3. Download and precompile all third party libraries, ignoring errors (some have broken tests or whatever)
COPY go.mod .
COPY go.sum .
RUN go mod download -x
SHELL ["/bin/bash", "-o", "pipefail", "-c"]
RUN go list -m all | tail -n +2 | cut -f 1 -d " " | awk 'NF{print $0 "/..."}' | CGO_ENABLED=0 GOOS=linux xargs -n1 go build -v -installsuffix cgo -i; echo "done"

# 4. Add the sources
COPY . .

# 5. Compile! Should only compile our sources since everything else is precompiled
RUN CGO_ENABLED=0 GOOS=linux go build -v -installsuffix cgo -o api -ldflags "-s -w" ./cmd/api

FROM debian:11.3-slim@sha256:b771c35d1e6ecf2556718ad3c0f481b4a04c1fbc133c609643acc9dd6743ead2

RUN apt-get update && apt-get install --no-install-recommends -y ca-certificates=20210119 && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /api
COPY --from=builder /go/src/voogle/api /api

CMD ["./api"]
