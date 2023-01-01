#!/bin/bash

cp -af bin/hotkeys-pi /home/pi/
chmod +x /home/pi/hotkeys-pi

cp -af runboot.sh /home/pi/
chmod +x /home/pi/*.sh

cp -af config.json /home/pi/
