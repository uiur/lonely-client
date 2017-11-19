## raspberry pi
requirements:

```
sudo apt update && sudo apt install -y fswebcam
```

deploy:
```
make
scp build/linux_arm/lonely pi@pi.local:~/bin/
```

run on device:
```
LONELY_DEVICE_TOKEN=token LONELY_SERVER_HOST=https://lonely.example.com /home/pi/bin/lonely run
```
