# syntax=docker/dockerfile:1
ARG IMAGE=ghcr.io/zoomoid/assignment/cli:latest
FROM ${IMAGE}

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

COPY docker-entrypoint.sh /
ENTRYPOINT ["/docker-entrypoint.sh"]

ENV MIKTEX_USERCONFIG=/miktex/.miktex/texmfs/config
ENV MIKTEX_USERDATA=/miktex/.miktex/texmfs/data
ENV MIKTEX_USERINSTALL=/miktex/.miktex/texmfs/install

WORKDIR /miktex/work

RUN mkdir -p /miktex/user/texmf/tex/latex

# add assignments.cls to the TEXMF directory
ADD https://github.com/zoomoid/assignments/releases/latest/download/csassignments.ins /miktex/user/texmf/tex/latex/csassignments/csassignments.ins
ADD https://github.com/zoomoid/assignments/releases/latest/download/csassignments.dtx /miktex/user/texmf/tex/latex/csassignments/csassignments.dtx

WORKDIR /miktex/user/texmf/tex/latex/csassignments

RUN latex csassignments.ins

WORKDIR /miktex/work

RUN initexmf --register-root=/miktex/user/texmf && \
  initexmf --update-fndb

