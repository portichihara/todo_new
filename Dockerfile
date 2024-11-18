FROM golang:1.22-alpine AS builder

WORKDIR /build

RUN apk add --no-cache postgresql-client

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o app

FROM alpine:latest

RUN apk add --no-cache postgresql-client

WORKDIR /app

COPY --from=builder /build/app ./
COPY --from=builder /build/templates ./templates/

COPY entrypoint.sh ./
RUN chmod +x /app/app && \
    chmod +x /app/entrypoint.sh

RUN ls -la /app && ls -la /app/templates

EXPOSE 8080

CMD ["./entrypoint.sh"]