FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY cmd/ ./cmd/
COPY internal/ ./internal/
RUN CGO_ENABLED=0 GOOS=linux go build -o dev-tools ./cmd/.

FROM alpine:latest
RUN apk add --no-cache ca-certificates libc6-compat
WORKDIR /app
COPY --from=builder /app/dev-tools .
COPY templates/ ./templates/
COPY assets/ ./assets/

EXPOSE 8090
CMD ["./dev-tools"]