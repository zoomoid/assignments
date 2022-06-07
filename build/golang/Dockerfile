# syntax=docker/dockerfile:1

FROM golang:1.17 AS builder

WORKDIR /build/

COPY . .

WORKDIR /build/cli/

RUN go mod download && \
  go build -o ./assignment

FROM ubuntu:20.04 AS runner

ARG DEBIAN_FRONTEND=noninteractive

RUN apt-get update && \
  apt-get install -y --no-install-recommends \
  apt-transport-https \
  ca-certificates \
  dirmngr \
  ghostscript \
  gnupg \
  gosu \
  make \
  perl

RUN apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys D6BC243565B2087BC3F897C9277A7293F59E4889 && \
  echo "deb http://miktex.org/download/ubuntu focal universe" | tee /etc/apt/sources.list.d/miktex.list

RUN apt-get update -y && \
  apt-get autoremove -y && \
  apt-get install -f -y miktex

RUN miktexsetup finish \
  && initexmf --admin --set-config-value=[MPM]AutoInstall=1 \
  && mpm --admin --update-db \
  && mpm --admin \
  --install amsfonts \
  --install biber-linux-x86_64 \
  && initexmf --admin --update-fndb

COPY build/golang/docker-entrypoint.sh /
ENTRYPOINT ["/docker-entrypoint.sh"]

ENV MIKTEX_USERCONFIG=/miktex/.miktex/texmfs/config
ENV MIKTEX_USERDATA=/miktex/.miktex/texmfs/data
ENV MIKTEX_USERINSTALL=/miktex/.miktex/texmfs/install

WORKDIR /miktex/work

RUN mkdir -p /miktex/user/texmf/tex/latex

# add assignments.cls to the TEXMF directory
ADD https://github.com/zoomoid/assignments/releases/latest/download/assignments.ins /miktex/user/texmf/tex/latex/assignments/assignments.ins
ADD https://github.com/zoomoid/assignments/releases/latest/download/assignments.dtx /miktex/user/texmf/tex/latex/assignments/assignments.dtx

RUN latex /miktex/user/texmf/tex/latex/assignments/assignments.ins

RUN initexmf --register-root=/miktex/user/texmf && \
  initexmf --update-fndb

CMD ["bash"]