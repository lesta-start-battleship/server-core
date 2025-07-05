package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"lesta-battleship/server-core/pkg/packets"
	"log"
	"net/url"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
)

var (
	pathFlag  = flag.String("path", "/random", "")
	queryFlag = flag.String("query", "", "")
)

func main() {
	flag.Parse()

	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: *pathFlag, RawQuery: *queryFlag}
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
			if err := conn.ReadJSON(&packet); err != nil {
				log.Println("read err: ", err)
				return
			}

			fmt.Println(packet.Body)
		}
	}()

	// id := rand.Text()

	go func() {
		for {
			select {
			default:
				var text string
				fmt.Scanln(&text)

				switch text {
				case "quit":
					SendPacket(conn, packets.NewDisconnect(""))

					cancel()

					return
				case "create":
					SendPacket(conn, packets.NewCreateRoom(""))
				case "join":
					var roomId string
					fmt.Scanln(&roomId)

					SendPacket(conn, packets.NewJoinRoom("", roomId))
				default:
					SendPacket(conn, packets.NewPlayerMessage("", text))
				}
			case <-done.Done():
				conn.WriteMessage(websocket.CloseMessage,
					websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				log.Println("Closed")
			}
		}
	}()

	<-done.Done()
}

func SendPacket(conn *websocket.Conn, packet packets.Packet) {
	// packet := packets.Packet{SenderId: senderId, Type: msg.String(), Body: msg}

	buffer := new(bytes.Buffer)
	json.NewEncoder(buffer).Encode(packet)

	test := packets.Packet{}
	json.Unmarshal(buffer.Bytes(), &test)

	err := conn.WriteMessage(websocket.TextMessage, buffer.Bytes())
	// err := conn.WriteJSON(packet)
	if err != nil {
		log.Println("write err: ", err)
		return
	}
}
