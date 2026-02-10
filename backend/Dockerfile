# Build stage
FROM golang:1.24.11-alpine AS build
WORKDIR /app
COPY . .
# 使用 vendor 目录
RUN go build -mod=vendor -o main

# Run stage
FROM alpine:latest
WORKDIR /app
COPY --from=build /app/main .
COPY --from=build /app/config/config.yaml ./config/config.yaml
COPY --from=build /app/template ./template
EXPOSE 2400
CMD ["/app/main"]