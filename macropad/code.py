"""
Control the Collections and the VLC Player via HID commands
"""


import json
from adafruit_macropad import MacroPad
import displayio
import terminalio
from adafruit_display_shapes.rect import Rect
from adafruit_display_text import label
import usb_cdc

"""
Keypad Controls
Play: 		Key 0
Pause:		Key 1
Stop:		Key 2
Prev:		Key 3
Next:		Key 5
Vol Dwn:	Key 6
Vol Up:		Key 8
Clear:	    Key 9
Reset:	    Key 11Clear

Rotary Encoder: Select the Playlist

"""

# INITIALIZATION
macropad = MacroPad()
macropad.display.auto_refresh = False
macropad.pixels.auto_write = False

# serial comms
serial = usb_cdc.data
in_data=bytearray()

# configruation
configuration = None

# playlists for the encoder switch
playlists = None
playlistIndexes = []
encoder_last_position = None
pl_index = ""
encoder_switch_last = macropad.encoder_switch_debounced.pressed

# Display Layout
display_lines = macropad.display_text(title="Collection")
group = displayio.Group()
for key_index in range(12):
    x = key_index % 3
    y = key_index // 3
    group.append(label.Label(terminalio.FONT, text='', color=0xFFFFFF,
                             anchored_position=((macropad.display.width - 1) * x / 2,
                                                macropad.display.height - 1 -
                                                (3 - y) * 12),
                             anchor_point=(x / 2, 1.0)))
group.append(Rect(0, 0, macropad.display.width, 12, fill=0xFFFFFF))
group.append(label.Label(terminalio.FONT, text='', color=0x000000,
                         anchored_position=(0, 0),
                         anchor_point=(0.0, 0.0)))
macropad.display.show(group)

# Functions and macros
controlKeys = [
    ('Play', 'cmd|play'),
    ('Pause', 'cmd|pause'),
    ('Stop', 'cmd|stop'),
    ('Prev', 'cmd|prev'),
    ('', ''),
    ('Next', 'cmd|next'),
    ('Vol Dn', 'cmd|voldown'),
    ('', ''),
    ('Vol Up','cmd|volup'),
    ('Clear', 'cmd|clear'),
    ('', ''),
    ('Reset', 'cmd|reset'),
]
group[0].text = controlKeys[0][0]
group[1].text = controlKeys[1][0]
group[2].text = controlKeys[2][0]
group[3].text = controlKeys[3][0]
group[5].text = controlKeys[5][0]
group[6].text = controlKeys[6][0]
group[8].text = controlKeys[8][0]
group[9].text = controlKeys[9][0]
group[11].text = controlKeys[11][0]


def clearAll():
    macropad.keyboard.release_all()
    macropad.consumer_control.release()
    macropad.mouse.release_all()
    macropad.stop_tone()
    macropad.pixels.show()
    macropad.display.refresh()

def sendData(data):
    msg = f'{data}\r'
    serialBuff = bytearray(msg)
    serial.write(serialBuff)

def loadCollection(collection):
    global playlists
    global playlistIndexes
    playlists = collection
    playlistIndexes = list(playlists.keys())
    playlistIndexes.sort()
    group[13].text = "pl ready"
    print("collection loaded")


def processSerialInput(data):
    if "_collections_" in data.decode("utf-8"):
        print("load collection")
        loadCollection(json.loads(data))
    else:
        print("unknown input")


# Main LOOP
#reset
sendData('cmd|reset')
while True:
    clearAll()
    
    # Keypad Events
    key_event = macropad.keys.events.get()
    if key_event:
        if key_event.key_number < 12:
            if key_event.pressed:
                macropad.pixels[key_event.key_number] = 0x002EB8
                print(group[9].text)                
                print(controlKeys[key_event.key_number])
                # send command
                sendData(controlKeys[key_event.key_number][1])
                group[9].text = "Key: {}".format(key_event.key_number)
            else:
                macropad.pixels[key_event.key_number] = 0
                group[9].text = controlKeys[9][0]
    
    # Encoder Events
    # Read encoder position. If it's changed, switch PL
    position = macropad.encoder
    if position != encoder_last_position:
        print("encoder {}".format(position))
        if len(playlistIndexes)>0:
            plPos = position % len(playlistIndexes)
            pl_index=playlistIndexes[plPos]
            playlistName = playlists[pl_index].strip()
            print("{} - {}".format(pl_index, playlistName))
            group[13].text = "{}".format(playlistName[:24])
        encoder_last_position = position
    # encoder push switch
    macropad.encoder_switch_debounced.update()
    if macropad.encoder_switch_debounced.pressed:
        group[9].text = "PL: {}". format(pl_index)
        print(group[9].text)
        # send playlist change
        sendData('pl|{}'.format(pl_index))
    else:        
        group[9].text = controlKeys[9][0]

    # Read from serial line
    if serial.in_waiting > 0:
        guard=0
        while(serial.in_waiting>0):
            guard+=1
            if(guard>1024):
                break
            byte = serial.read(1)
            if byte == b'\r':
                print(in_data.decode("utf-8"))
                processSerialInput(in_data)
                in_data = bytearray()
            else:
                in_data += byte
