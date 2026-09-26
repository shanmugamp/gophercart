FROM golang:1.26.2-alpine AS build

WORKDIR /src
ENV GOWORK=off

COPY services/product-service/go.mod services/product-service/go.sum ./
RUN go mod download

COPY services/product-service/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/product-service ./cmd/server

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=build /out/product-service /product-service

EXPOSE 8080
ENTRYPOINT ["/product-service"]