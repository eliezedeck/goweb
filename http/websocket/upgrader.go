package websocket

import (
	gws "github.com/gorilla/websocket"
)

// U is a global gorilla/websocket Upgrader
var U gws.Upgrader

// SetupWebsocketUpgrader sets up the global U Upgrader
func SetupWebsocketUpgrader() {
	U = gws.Upgrader{
		ReadBufferSize:  2048,
		WriteBufferSize: 2048,
	}
}
