#!/bin/bash
echo "---Checking if UID: ${UID} matches user---"
usermod -u ${UID} ${USER}
echo "---Checking if GID: ${GID} matches user---"
usermod -g ${GID} ${USER}
echo "---Setting umask to ${UMASK}---"
umask ${UMASK}

echo "---Checking for optional scripts---"
cp -f /opt/custom/user.sh /opt/scripts/start-user.sh > /dev/null 2>&1 ||:
cp -f /opt/scripts/user.sh /opt/scripts/start-user.sh > /dev/null 2>&1 ||:

if [ -f /opt/scripts/start-user.sh ]; then
  echo "---Found optional script, executing---"
  chmod -f +x /opt/scripts/start-user.sh ||:
  /opt/scripts/start-user.sh || echo "---Optional Script has thrown an Error---"
else
  echo "---No optional script found, continuing---"
fi

echo "---Starting...---"
chown -R root:${GID} /opt/scripts
chmod -R 750 /opt/scripts
chown -R ${UID}:${GID} ${DATA_DIR}

term_handler() {
  su $USER -c "tmux send-keys -t BeamMP-Server \"exit\" C-m"
  tail --pid=$(pidof BeamMP-Server) -f 2>/dev/null
  sleep 0.5
  exit 143;
}

trap 'kill ${!}; term_handler' SIGTERM
su ${USER} -c "/opt/scripts/start-server.sh" &
killpid="$!"
while true
do
	wait $killpid
	exit 0;
done