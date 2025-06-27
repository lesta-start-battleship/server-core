package main

import (
	"context"
	"fmt"
	"lesta-battleship/server-core/pkg/packets"
	"log"
	"net/url"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
)

func main() {
	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws"}
	log.Printf("connecting to %s", u.String())
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return
	}
	defer conn.Close()

	done, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	go func() {
		for {
			var packet packets.Packet
			conn.ReadJSON(&packet)
			if err != nil {
				log.Println("read err: ", err)
				return
			}
			log.Println(packet)
		}
	}()

	for {
		var text string
		fmt.Scanln(&text)
		msg := packets.PlayerMessage{Msg: text}
		packet := packets.Packet{SenderId: "", Body: msg}

		// err := conn.WriteMessage(websocket.TextMessage, []byte(text))
		err := conn.WriteJSON(packet)
		if err != nil {
			log.Println("write err: ", err)
			return
		}

		select {
		default:
			continue
		case <-done.Done():
			return
		}
	}
}
