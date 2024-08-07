#!/bin/bash
if [ "${PRERELEASE}" == "true" ]; then
  LAT_V="$(wget -qO- https://api.github.com/repos/BeamMP/BeamMP-Server/releases | jq -r 'map(select(.prerelease)) | first | .tag_name')"
else
  LAT_V="$(wget -qO- https://api.github.com/repos/BeamMP/BeamMP-Server/releases/latest | jq -r '.tag_name')"
fi
CUR_V="$(find ${DATA_DIR} -maxdepth 1 -type f -name "beamngmp_v*" 2>/dev/null | cut -d '_' -f2)"

if [ -z "${LAT_V}" ]; then
  if [ -z "${CUR_V}" ]; then
    echo "---Can't get latest version from BeamNG-MP-Server and found no local installed version!---"
        sleep infinity
  else
    echo "---Can't get latest version from BeamNG-MP-Server, falling back to installed version ${CUR_V}---"
        LAT_V="${CUR_V}"
  fi
fi

echo "---Version Check---"
if [ -z "${CUR_V}" ]; then
  echo "---BeamNG-MP-Server not installed, installing...---"
  DL_URL="$(wget -qO- https://api.github.com/repos/BeamMP/BeamMP-Server/releases | jq -r --arg LAT_V "$LAT_V" '.[] | select(.tag_name == $LAT_V) | .assets[] | .browser_download_url | match("^.*BeamMP-Server\\.debian.12.x86_64*$") | .string')"
  cd ${DATA_DIR}
  rm -rf ${DATA_DIR}/BeamMP-Server 2>/dev/null
  if wget -q -nc --show-progress --progress=bar:force:noscroll -O ${DATA_DIR}/BeamMP-Server "${DL_URL}" ; then
    echo "---Sucessfully downloaded BeamNG-MP-Server ${LAT_V}---"
  else
    echo "---Something went wrong, can't download BeamNG-MP-Server ${LAT_V}, putting container in sleep mode---"
    sleep infinity
  fi
  chmod +x ${DATA_DIR}/BeamMP-Server
  touch ${DATA_DIR}/beamngmp_${LAT_V}
elif [ "${CUR_V}" != "${LAT_V}" ]; then
  echo "---Version missmatch, installed ${CUR_V}, downloading and installing latest ${LAT_V}...---"
  cd ${DATA_DIR}
  rm -rf ${DATA_DIR}/beamngmp_v* ${DATA_DIR}/BeamMP-Server
  DL_URL="$(wget -qO- https://api.github.com/repos/BeamMP/BeamMP-Server/releases | jq -r --arg LAT_V "$LAT_V" '.[] | select(.tag_name == $LAT_V) | .assets[] | .browser_download_url | match("^.*BeamMP-Server\\.debian.12.x86_64*$") | .string')"
  cd ${DATA_DIR}
  if wget -q -nc --show-progress --progress=bar:force:noscroll -O ${DATA_DIR}/BeamMP-Server "${DL_URL}" ; then
    echo "---Sucessfully downloaded BeamNG-MP-Server ${LAT_V}---"
  else
    echo "---Something went wrong, can't download BeamNG-MP-Server ${LAT_V}, putting container in sleep mode---"
    sleep infinity
  fi
  chmod +x ${DATA_DIR}/BeamMP-Server
  touch ${DATA_DIR}/beamngmp_${LAT_V}
elif [ "${CUR_V}" == "${LAT_V}" ]; then
        echo "---BeamNG-MP-Server $CUR_V up-to-date---"
fi

echo "---Prepare Server---"
if [ ! -f ~/.tmux.conf ]; then
  echo "set-option -g status off
set-option -g display-time 5000
unbind -n C-c
bind-key -n C-c display-message 'Blocked. Please use to command \"exit\" to shutdown the server or close this window to exit the terminal.'" > ~/.tmux.conf
fi
if [ ! -f "${DATA_DIR}/ServerConfig.toml" ]; then
  echo "---No ServerConfig.toml found, copying...---"
  cp -f /opt/config/ServerConfig.toml ${DATA_DIR}
fi
chmod -R ${DATA_PERM} ${DATA_DIR}

echo "---Start Server---"
if [ "${ENABLE_WEBCONSOLE}" == "true" ]; then
    /opt/scripts/start-gotty.sh 2>/dev/null &
fi
cd ${DATA_DIR}
if [ ! -f ${DATA_DIR}/BeamMP-Server ]; then
  echo "---Something went wrong, can't find the executable, putting container into sleep mode!---"
  sleep infinity
else
  tmux new-session -d -s BeamMP-Server /beamngmp/BeamMP-Server ${GAME_PARAMS}
  if [ ! -f ${DATA_DIR}/Server.log ]; then
    ${DATA_DIR}/Server.log
  fi
  sleep 2
  /opt/scripts/start-watchdog.sh &
  tail -n 9999 -f ${DATA_DIR}/Server.log
fi