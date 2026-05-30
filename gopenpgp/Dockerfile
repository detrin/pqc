FROM golang:1.24-bookworm AS builder

ENV GOTOOLCHAIN=auto
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY cmd/ cmd/
COPY main.go .
RUN CGO_ENABLED=0 go build -o /pqcrypt ./cmd/pqcrypt/

FROM debian:bookworm-slim

COPY --from=builder /pqcrypt /usr/local/bin/pqcrypt

WORKDIR /data
ENTRYPOINT ["pqcrypt"]
CMD ["--help"]
