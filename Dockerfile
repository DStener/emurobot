FROM golang:alpine AS builder

# Install dependency 
RUN apk --no-cache add \
    git \
    ca-certificates

WORKDIR /app

COPY . . 

# Download go-package dependency 
RUN go mod download

# Build sim_dev
# RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /sim_dev ./sim_dev
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /emu_record ./emu_record
# RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /emu_play ./emu_play

# # Copy to simple minimalistic environment
# FROM alpine:latest

# RUN apk --no-cache add socat

# # COPY --from=builder /sim_dev /usr/local/bin/sim_dev
# COPY --from=builder /emu_record /usr/local/bin/emu_record
# # COPY --from=builder /emu_play /usr/local/bin/emu_play

FROM ubuntu:24.04

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        ffmpeg \
        v4l-utils \
        ca-certificates \
        kmod  \
        socat && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /emu_record /usr/local/bin/emu_record

# ENTRYPOINT ["/app/camera-service"]
