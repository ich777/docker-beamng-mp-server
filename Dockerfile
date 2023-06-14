FROM ich777/debian-baseimage:bullseye_amd64

LABEL org.opencontainers.image.authors="admin@minenet.at"
LABEL org.opencontainers.image.source="https://github.com/ich777/docker-beamng-mp-server"

RUN apt-get update && \
	apt-get -y install --no-install-recommends liblua5.3-0 libcurl4 jq && \
	rm -rf /var/lib/apt/lists/*

ENV DATA_DIR="/beamngmp"
ENV GAME_PARAMS=""
ENV PRERELEASE="false"
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