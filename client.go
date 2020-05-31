package main

import (
	"bufio"
	"errors"
	"net"
	"time"

	Log "github.com/jampajeen/go-async-socket/logger"
)

// Client ...
type Client struct {
	hub    *Hub
	conn   net.Conn
	sendCH chan []byte
	IDUser string
}

func (c *Client) writePump() {
	pingTicker := time.NewTicker(pingPeriod)

	defer func() {
		pingTicker.Stop()
		c.hub.unregisterCH <- c
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.sendCH:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				Log.Debug("Channel(sendCH) is closed: %s\n", c.conn.RemoteAddr().String())
				return
			}

			if _, err := c.conn.Write([]byte(message)); err != nil {
				Log.Error(errors.New("Error write"))
				return
			}
			// go c.hub.onSent(c, message)
			c.hub.onSent(c, message)

		case <-pingTicker.C:
			// c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			Log.Debug("Ping interval: %s \n", c.conn.RemoteAddr().String())
		}
	}
}

func (c *Client) readPump() {
	buf := bufio.NewReader(c.conn)

	defer func() {
		c.hub.unregisterCH <- c
		c.conn.Close()
	}()

	c.conn.SetReadDeadline(time.Now().Add(pongWait))

	for {
		data, err := buf.ReadString('\n')
		if err != nil {
			Log.Debug("Pong timeout: %s\n", c.conn.RemoteAddr().String())
			break
		}

		// go c.hub.onReceived(c, []byte(data))
		c.hub.onReceived(c, []byte(data))
		// c.sendCH <- []byte(data)
	}
}
