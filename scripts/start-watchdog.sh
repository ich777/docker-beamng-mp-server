#!/bin/bash
killpid="$(pidof BeamMP-Server)"
while true
do
  tail --pid=$killpid -f /dev/null
  kill "$(pidof tail)"
  exit 0
done