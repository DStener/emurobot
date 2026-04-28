FROM golang:alpine AS builder

# Install dependency 
RUN apk --no-cache add \
    git \
    ca-certificates

WORKDIR /app

COPY go.mod go.sum sim_dev/ ./

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /sim_dev .


FROM alpine:latest

RUN apk --no-cache add socat

COPY --from=builder /sim_dev /usr/local/bin/sim_dev