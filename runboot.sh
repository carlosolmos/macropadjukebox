#!/bin/bash
set -x
set -e

sleep 5
cd /home/pi
cvlc --rc-host "localhost:8888" -I rc > /dev/null 2>&1 &
sleep 5
./hotkeys-pi > /home/pi/hotkeys.log 2>&1 &
