FROM golang:1.27-alpine AS builder

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /build/habbit-graph-api ./cmd/main.go


FROM scratch

WORKDIR /app

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /build/habbit-graph-api /app/habbit-graph-api

EXPOSE 8080

ENTRYPOINT ["/app/habbit-graph-api"]