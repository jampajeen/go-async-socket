package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	_ "net/http/pprof"
	"os"
	"strconv"
	"strings"
	"syscall"

	Log "github.com/jampajeen/go-logger"
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
	go readCmdPipe(hub)   // TODO: remove if needed

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

func readCmdPipe(hub *Hub) {
	pipeName := "cmdpipe"
	//to create pipe: does not work in windows
	syscall.Mkfifo(pipeName, 0666)

	file, err := os.OpenFile(pipeName, os.O_RDONLY, os.ModeNamedPipe)
	if err != nil {
		log.Fatal("Open named pipe file error:", err)
	}
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err == nil {
			processCmd(hub, line)
		}
	}
}

func keyboardInput(hub *Hub) {
	fmt.Print("Enter command: \n")
	for {
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		processCmd(hub, text)
	}
}

func processCmd(hub *Hub, text string) {
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
