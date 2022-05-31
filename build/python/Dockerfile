FROM ubuntu:20.04

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

COPY build/python/docker-entrypoint.sh /
ENTRYPOINT ["/docker-entrypoint.sh"]

ENV MIKTEX_USERCONFIG=/miktex/.miktex/texmfs/config
ENV MIKTEX_USERDATA=/miktex/.miktex/texmfs/data
ENV MIKTEX_USERINSTALL=/miktex/.miktex/texmfs/install

WORKDIR /miktex/work

CMD ["bash"]

RUN apt install -y software-properties-common && \
  add-apt-repository -y ppa:deadsnakes/ppa && \
  apt-get install -y python3.10 python3-pip zip && \
  pip install requests

RUN update-alternatives --install /usr/local/bin/python3 python3 /usr/bin/python3.10 4
