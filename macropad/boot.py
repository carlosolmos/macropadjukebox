import usb_cdc
try:
    usb_cdc.enable(console=True, data=True)    # Enable console and data# Write your code here :-)
except Exception as e:
    print(e)
