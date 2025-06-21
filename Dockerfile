FROM golang:1.24-alpine AS builder
RUN apk add --no-cache git ca-certificates tzdata
RUN adduser -D -g '' appuser
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -trimpath \
    -tags netgo \
    -o main cmd/main.go

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /app/main .
COPY --from=builder /app/config/ ./config/
ENV TZ=Asia/Jakarta
USER appuser
EXPOSE 3000
ENTRYPOINT ["./main"]