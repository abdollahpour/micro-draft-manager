ARG GOLANG_VERSION=1.16
FROM golang:${GOLANG_VERSION}-alpine AS builder
RUN apk add ca-certificates git
WORKDIR /app
ADD go.mod ./
ADD go.sum ./
RUN go mod download
ADD cmd ./cmd
ADD internal ./internal
ARG APP_VERSION=developement
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-X main.Version=${APP_VERSION}" -a -installsuffix cgo -o app /app/cmd/server/main.go

FROM alpine:3.13
WORKDIR /app
COPY --from=builder /app/app .
ENTRYPOINT [ "./app" ]