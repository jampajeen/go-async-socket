package main

import Log "github.com/jampajeen/go-async-socket/logger"

// Hub ...
type Hub struct {
	clients       map[*Client]bool
	clientUserMap map[string]*Client
	broadcastCH   chan []byte
	registerCH    chan *Client
	unregisterCH  chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcastCH:   make(chan []byte),
		registerCH:    make(chan *Client),
		unregisterCH:  make(chan *Client),
		clients:       make(map[*Client]bool),
		clientUserMap: make(map[string]*Client),
	}
}

func (hub *Hub) run() {
	for {
		select {
		case client := <-hub.registerCH:
			hub.clients[client] = true
			hub.clientUserMap[client.IDUser] = client
			// go hub.onConnected(client)
			hub.onConnected(client)

		case client := <-hub.unregisterCH:
			if _, ok := hub.clients[client]; ok {
				delete(hub.clientUserMap, client.IDUser)
				delete(hub.clients, client)
				close(client.sendCH)
				// go hub.onDisconnected(client)
				hub.onDisconnected(client)
			}

		case message := <-hub.broadcastCH:
			for client := range hub.clients {
				select {
				case client.sendCH <- message:
				default:
				}
			}
		}
	}
}

func (hub *Hub) closeAll() {
	for c := range hub.clients {
		hub.unregisterCH <- c
	}
}

func (hub *Hub) send(data []byte, client *Client) {
	client.sendCH <- data
}

func (hub *Hub) sendBroadcastCH(data []byte) {
	hub.broadcastCH <- data
}

func (hub *Hub) sendToIDUser(data []byte, idUser string) {
	if client, ok := hub.clientUserMap[idUser]; ok {
		client.sendCH <- data
	}
}

func (hub *Hub) onReceived(client *Client, data []byte) {
	Log.Debug("Recv: %s <- %s", client.conn.RemoteAddr().String(), string(data))
}

func (hub *Hub) onSent(client *Client, data []byte) {
	Log.Debug("Sent: %s -> %s", client.conn.RemoteAddr().String(), string(data))
}

func (hub *Hub) onConnected(client *Client) {
	Log.Debug("Client connected: %s", client.conn.RemoteAddr().String())
}

func (hub *Hub) onDisconnected(client *Client) {
	Log.Debug("Client disconnected: %s", client.conn.RemoteAddr().String())
}
