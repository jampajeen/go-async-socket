# Basic TCP socket server

## How to try:

Connect to server by using netcat
```
nc 127.0.0.1 5555
```

Client send text to server by enter text on netcat console

Server send text to all clients by enter text on server console

Server send text to specific client by enter this command on server console
```
_to_ [peer addr] [text]
ex. 
_to_ 127.0.0.1:51358 hello

*** you can find the peer addr on server console ***
```

Server close all connections by enter this command on server console
```
_close_
```
