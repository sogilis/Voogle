FROM rust:1.60.0-slim-buster as builder

WORKDIR /usr/src/env-generator
COPY . .
RUN cargo build --release

FROM debian:11.3-slim

LABEL org.opencontainers.image.source="https://github.com/sogilis/Voogle"

WORKDIR /env
COPY --from=builder /usr/src/env-generator/target/release/env_generator /env_generator
WORKDIR /

CMD ["sh", "-c", "cd /env && /env_generator"]
