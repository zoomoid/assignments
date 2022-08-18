# syntax=docker/dockerfile:1

FROM golang:1.17 AS builder

ARG VERSION="0.0.0-dev.0"
ARG EXTRA_LDFLAGS='-w -s'

WORKDIR /build/

COPY . .

RUN go mod download 

RUN CGO_ENABLED=0 go build -ldflags "${EXTRA_LDFLAGS} -X main.Version=${VERSION}" -o ./assignmentctl

FROM ubuntu:focal AS cli

WORKDIR /

COPY --from=builder /build/assignmentctl /usr/local/bin/assignmentctl

RUN chmod +x /usr/local/bin/assignmentctl

CMD ["assignmentctl"]