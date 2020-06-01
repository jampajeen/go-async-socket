# Basic TCP socket server

## How to try:

1. Connect to server by using netcat
```
nc 127.0.0.1 5555
```

2. Client send text to server by enter text on netcat console

3. Server send text to all clients by enter text on server console

4. Server send text to specific client by enter this command on server console
```
_to_ [peer addr] [text]
*** you can find the peer addr on server console ***
```
ex.
``` 
_to_ 127.0.0.1:51358 hello
```
or cmd pipe
```
echo "_to_ 127.0.0.1:51358 hello" > cmdpipe
```

5. Server close all connections by enter this command on server console
```
_close_
```
or pipe 
```
echo "_close_" > cmdpipe
```
