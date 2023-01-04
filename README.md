# macropadjukebox
Service to capture custom hotkeys and control VLC


Used to control VLC in Raspberry Pi from an [Adafruit Macropad RP2040](https://learn.adafruit.com/adafruit-macropad-rp2040/overview)


Circuit Python libraries

https://learn.adafruit.com/adafruit-macropad-rp2040/circuitpython

https://learn.adafruit.com/adafruit-macropad-rp2040/macropad-circuitpython-library



Based on the cdc usb serial comms

Build:
```
mkdir -p bin
go build -o bin/hotkeys-pi main.go
```

## Control Keys

```
Play: 		Key 0
Pause:		Key 1
Stop:		Key 2
Prev:		Key 3
Next:		Key 5
Vol Dwn:	Key 6
Vol Up:		Key 8
Clear:	    Key 9
Reset:	    Key 11


cmd|{keycommand}
```
## Playlist Keys
```
pl|{PL index}
```






#start service in rpi

@reboot cd /home/pi && ./runboot.sh >> /dev/null 2> /home/pi/cron.log &



