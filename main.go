package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/gorilla/websocket"
)

func main() {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	http.HandleFunc("/generate", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Printf("error occured serving %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Printf("Connection upgraded...")
		defer conn.Close()

		for {
			messageType, msg, err := conn.ReadMessage()
			if err != nil {
				fmt.Printf("error occured when reading message %s", err)
				return
			}

			log.Printf("MessageType is %d, msg is %s", messageType, string(msg))

			parameters := struct {
				Filename string `json:"filename"`
			}{}

			err = json.Unmarshal(msg, &parameters)
			if err != nil {
				fmt.Printf("error reading request parameters %s", err)
				return
			}

			fileBytes, err := os.ReadFile(path.Join("samples", parameters.Filename))
			if err != nil {
				fmt.Printf("unable to read requested file %s", err)
				return
			}

			err = conn.WriteMessage(messageType, fileBytes)
			if err != nil {
				fmt.Printf("unable to send file over the web socket connection %s", err)
				return
			}

		}
	})

	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		fmt.Printf("error occured serving %s", err)
	}
}
