FROM golang:1.23.4-bullseye

ARG SRCROOT=/usr/local/src/chainlink
WORKDIR ${SRCROOT}

ADD go.* ./
RUN go mod download
RUN mkdir -p tools/bin
