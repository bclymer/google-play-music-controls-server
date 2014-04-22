package main

import (
	"fmt"
	socketio "github.com/googollee/go-socket.io"
	"log"
	"net/http"
)

func onConnect(ns *socketio.NameSpace) {
	fmt.Println("connected:", ns.Id(), " in channel ", ns.Endpoint())
}

func onDisconnect(ns *socketio.NameSpace) {
	fmt.Println("disconnected:", ns.Id(), " in channel ", ns.Endpoint())
}

func main() {
	sock_config := &socketio.Config{}
	sock_config.HeartbeatTimeout = 2
	sock_config.ClosingTimeout = 4

	sio := socketio.NewSocketIOServer(sock_config)

	// Handler for new connections, also adds socket.io event handlers
	sio.On("connect", onConnect)
	sio.On("disconnect", onDisconnect)
	sio.On("ping", func(ns *socketio.NameSpace){
		sio.Broadcast("pong", nil)
	})

	//this will serve a http static file server
	sio.Handle("/", http.FileServer(http.Dir("./public/")))
	sio.HandleFunc("/play", func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("play");
			sio.Broadcast("music", "play");
	});
	sio.HandleFunc("/next", func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("next");
			sio.Broadcast("music", "next");
	});
	sio.HandleFunc("/prev", func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("prev");
			sio.Broadcast("music", "prev");
	});
	//startup the server
	log.Fatal(http.ListenAndServe(":3000", sio))
}