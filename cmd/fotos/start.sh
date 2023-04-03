#! /bin/sh
sh libcamera.sh&               # start camera
DISPLAY=:0 /home/pi/hhpc/hhpc& # disable cursor
sh dpms.sh&                    # make it not blank the screen
sleep 5 # wait a bit for libcamera to be ready

sudo ./fotos
