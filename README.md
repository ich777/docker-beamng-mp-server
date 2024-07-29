# BeamNG.drive MP Server in Docker optimized for Unraid
This Docker will download and install BeamNG.drive-MP-Server.

**ATTENTION:** To get the server working please generate a Key over here: https://beammp.com/keymaster (you can get a full tutorial on how to obtain a key here: https://wiki.beammp.com/en/home/server-installation) and put it in your ServerConfig.toml at the entry "AuthKey".

**ServerConfig.toml:** Please head over to https://wiki.beammp.com/en/home/server-maintenance to see all available options and descriptions from the ServerConfig.toml

**Update Notice:** The container will check for a new version on each start/restart.

>**WEB CONSOLE:** You can connect to the BeamMP-Server console by opening your browser and go to HOSTIP:9045 (eg: 192.168.1.1:9045) or click on WebUI on the Docker page within Unraid.

## Env params
| Name | Value | Example |
| --- | --- | --- |
| DATA_DIR | Folder for gamefiles | /beamngmp |
| GAME_PARAMS | Values to start the server | empty |
| UID | User Identifier | 99 |
| GID | Group Identifier | 100 |

## Run example
```
docker run --name BeamNG-MP -d \
	-p 30814:30814 -p 30814:30814/udp -p 9045:8080 \
	--env 'UID=99' \
	--env 'GID=100' \
	--env 'UMASK=0000' \
	--volume /path/to/beamngmp:/beamngmp \
	ich777/beamng-mp-server:latest
```

This Docker was mainly edited for better use with Unraid, if you don't use Unraid you should definitely try it!

#### Support Thread: https://forums.unraid.net/topic/79530-support-ich777-gameserver-dockers/