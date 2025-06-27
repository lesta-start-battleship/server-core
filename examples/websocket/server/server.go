package main

import (
	"lesta-battleship/server-core/internal/multiplayer"
	"net/http"
)

func main() {
	hub := multiplayer.NewHub()
	go hub.Run()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		hub.ConnectPlayer(w, r)
	})
	http.ListenAndServe(":8080", nil)
}
