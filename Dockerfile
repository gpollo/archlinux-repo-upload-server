FROM golang:latest AS builder
WORKDIR /build
ADD main.go .
RUN go build -o upload-server main.go

FROM archlinux:latest
WORKDIR /srv/repo
COPY --from=builder /build/upload-server ..
ENTRYPOINT ../upload-server
