FROM ubuntu:24.04

RUN apt update && \
    apt install -y \
        git \
        socat \
        gopls \
        golang &&\
    rm -rf /var/lib/apt/lists/*