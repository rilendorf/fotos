#! /bin/sh
cd /home/pi/fotos/cmd/fotos

echo "deleting tmpimage to ensure no wrong images are uploaded"
rm tmpimage.jpg

echo "starting libcamera"
sh /home/pi/fotos/cmd/fotos/libcamera.sh&               # start camera

echo "cursor be gone"
DISPLAY=:0 /home/pi/hhpc/hhpc& # disable cursor

echo "no more screen goto sleep"
sh /home/pi/fotos/cmd/fotos/dpms.sh&                    # make it not blank the screen
sleep 5 # wait a bit for libcamera to be ready

echo "starting fotos"
sudo /home/pi/fotos/cmd/fotos/fotos
