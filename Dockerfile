FROM alpine:3.20 AS root-certs
RUN apk add -U --no-cache ca-certificates
RUN addgroup -g 1001 app
RUN adduser app -u 1001 -D -G app /home/app

FROM golang:1.22 AS builder
WORKDIR /app
COPY --from=root-certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs
RUN git clone https://github.com/vindosVP/snauth.git .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X main.buildCommit=$(git rev-list -1 HEAD) -X main.buildTime=$(date -u '+%Y-%m-%d_%I:%M:%S%p') -X main.version=$(git describe --tags --abbrev=0)" -o ./snauth ./cmd/main.go

FROM scratch AS final
COPY --from=root-certs /etc/passwd /etc/passwd
COPY --from=root-certs /etc/group /etc/group
COPY --chown=1001:1001 --from=root-certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs
COPY --chown=1001:1001 --from=builder /app/snauth /snauth
USER app
ENTRYPOINT ["/snauth"]