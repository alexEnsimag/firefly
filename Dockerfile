
FROM golang:1.24 AS builder

WORKDIR /app
    
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/essayanalyzer
    

FROM alpine:latest

WORKDIR /usr/local/bin

RUN adduser -D appuser
COPY --from=builder /app/app /usr/local/bin/app
COPY --from=builder /app/data /usr/local/bin/data
USER appuser

ENTRYPOINT ["/usr/local/bin/app"]