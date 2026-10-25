# syntax=docker/dockerfile:1

FROM golang:1.26 AS builder

WORKDIR /src

COPY contracts/go.mod contracts/go.sum* ./contracts/
COPY services/product-service/go.mod services/product-service/go.sum* ./services/product-service/

WORKDIR /src/contracts
RUN go mod download

WORKDIR /src/services/product-service
RUN go mod download

WORKDIR /src

COPY contracts ./contracts
COPY services/product-service ./services/product-service

WORKDIR /src/services/product-service

RUN CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/product-service \
    ./cmd/server


FROM alpine:3.22 AS runtime

RUN addgroup -S appgroup && \
    adduser -S appuser -G appgroup && \
    apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /out/product-service /app/product-service

USER appuser

EXPOSE 8080
EXPOSE 50051

ENTRYPOINT ["/app/product-service"]