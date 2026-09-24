FROM hub.hamdocker.ir/library/golang:1.25-alpine AS builder

WORKDIR /src

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/code-reviewer .

FROM hub.hamdocker.ir/library/alpine:3.21

RUN apk add --no-cache ca-certificates tzdata \
    && adduser -D -H -u 1000 appuser

WORKDIR /app

COPY --from=builder /out/code-reviewer /app/code-reviewer
COPY ./docker/docker-entrypoint.sh /entrypoint.sh
COPY ./docker/docker-entrypoint.prod.sh /entrypoint.prod.sh
RUN sed -i 's/\r$//' /entrypoint.sh /entrypoint.prod.sh \
    && chmod +x /entrypoint.sh /entrypoint.prod.sh /app/code-reviewer

USER appuser

EXPOSE 8080

ENTRYPOINT ["/entrypoint.sh"]
