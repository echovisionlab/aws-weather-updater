FROM golang:1.22.1-alpine AS builder
RUN apk add git

WORKDIR /build
COPY . ./
RUN go mod tidy
RUN go mod verify
RUN go mod download

RUN go build -o app .
WORKDIR /dist
RUN cp /build/app .

FROM alpine:latest
RUN apk add chromium

WORKDIR /app
LABEL org.opencontainers.image.source=https://github.com/echovisionlab/aws-weather-updater
LABEL org.opencontainers.image.licenses=MIT
COPY --from=builder /dist/app /app/

ENTRYPOINT ["./app"]
