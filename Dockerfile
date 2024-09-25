FROM golang:1.23-alpine AS builder

WORKDIR /

COPY go.mod go.sum .
RUN pwd
RUN go mod download && go mod verify

COPY . .
RUN CGO_ENABLE=0 go build -v -o healthd

FROM alpine:latest
COPY --from=builder /healthd ./
EXPOSE 80
