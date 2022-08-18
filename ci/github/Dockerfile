# syntax=docker/dockerfile:1

ARG BASE_IMAGE=ghcr.io/zoomoid/assignments/cli:alpine

FROM ${BASE_IMAGE}

RUN apk add github-cli && gh version
