FROM golang:1.26-alpine AS builder

WORKDIR /

COPY go.mod go.sum ./
RUN pwd
RUN go mod download && go mod verify

COPY . .
RUN CGO_ENABLED=0 go build -v -o healthd

FROM alpine:3.21
COPY --from=builder /healthd ./
EXPOSE 80
