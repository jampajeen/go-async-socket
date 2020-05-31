package main

import (
	"bufio"
	"fmt"
	"net"
	_ "net/http/pprof"
	"os"
	"strconv"
	"strings"

	Log "github.com/jampajeen/go-async-socket/logger"
)

func main() {
	ln, err := net.Listen("tcp", bindAddr+":"+strconv.Itoa(bindPort))
	if err != nil {
		panic(err)
	}
	Log.Info("Server is ready on [ %s ]\n", bindAddr+":"+strconv.Itoa(bindPort))

	hub := newHub()
	go hub.run()

	go keyboardInput(hub) // TODO: remove if needed

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			break
		}

		client := &Client{hub: hub, conn: conn, sendCH: make(chan []byte), IDUser: conn.RemoteAddr().String()} // use IP addr as user id
		hub.registerCH <- client

		go client.writePump()
		go client.readPump()
	}
}

func keyboardInput(hub *Hub) {
	fmt.Print("Enter text: \n")
	for {
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')

		if strings.HasPrefix(text, "_close_") {
			hub.closeAll()
		} else if strings.HasPrefix(text, "_to_") {
			s := strings.Split(text, " ")
			if len(s) > 2 {
				hub.sendToIDUser([]byte(s[2]), s[1])
			}
		} else {
			hub.sendBroadcastCH([]byte(text))
		}
	}
}
