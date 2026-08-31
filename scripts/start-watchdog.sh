#!/bin/bash
killpid="$(pidof BeamMP-Server)"
while true
do
  tail --pid=$killpid -f /dev/null
  kill "$(pidof beammp-webmanager)" 2>/dev/null
  kill "$(pidof tail)" 2>/dev/null
  exit 0
done
