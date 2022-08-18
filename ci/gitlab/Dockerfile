# syntax=docker/dockerfile:1

ARG BASE_IMAGE=ghcr.io/zoomoid/assignments/cli:alpine

FROM ${BASE_IMAGE}

ADD https://release-cli-downloads.s3.amazonaws.com/latest/release-cli-linux-amd64 /usr/local/bin/release-cli

RUN chmod +x /usr/local/bin/release-cli && release-cli -v

