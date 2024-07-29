FROM ich777/debian-baseimage

LABEL org.opencontainers.image.authors="admin@minenet.at"
LABEL org.opencontainers.image.source="https://github.com/ich777/docker-beamng-mp-server"

RUN apt-get update && \
	apt-get -y install --no-install-recommends liblua5.3-0 libcurl4 jq tmux && \
	rm -rf /var/lib/apt/lists/*

RUN wget -O /tmp/gotty.tar.gz https://github.com/yudai/gotty/releases/download/v1.0.1/gotty_linux_amd64.tar.gz && \
	tar -C /usr/bin/ -xvf /tmp/gotty.tar.gz && \
	rm -rf /tmp/gotty.tar.gz

ENV DATA_DIR="/beamngmp"
ENV GAME_PARAMS=""
ENV PRERELEASE="false"
ENV ENABLE_WEBCONSOLE="true"
ENV GOTTY_PARAMS="-w --title-format BeamNG-MP"
ENV UMASK=000
ENV UID=99
ENV GID=100
ENV DATA_PERM=770
ENV USER="beamngmp"

RUN mkdir $DATA_DIR && \
	useradd -d $DATA_DIR -s /bin/bash $USER && \
	chown -R $USER $DATA_DIR && \
	ulimit -n 2048

ADD /scripts/ /opt/scripts/
ADD /config/ /opt/config/
RUN chmod -R 777 /opt/scripts/

#Server Start
ENTRYPOINT ["/opt/scripts/start.sh"]