## raspberry pi
requirements:

```
sudo apt update && sudo apt install -y fswebcam
```

command:

```
fswebcam tmp.jpg && \
HOST=https://lonely-server.herokuapp.com/uiu ruby /home/pi/lonely-client/upload.rb tmp.jpg \
rm tmp.jpg
```

crontab -e :

```
* * * * * /home/pi/upload.sh
```
