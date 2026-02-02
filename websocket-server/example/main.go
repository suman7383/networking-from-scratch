package main

import (
	"fmt"
	"os"

	"github.com/suman7383/networking-from-scratch/websocket-server/internal/server"
	"github.com/suman7383/networking-from-scratch/websocket-server/internal/websocket"
)

func main() {

	s := server.NewServer(":8443", func(w websocket.DataWriter, data []byte) {
		fmt.Println("Received data", string(data))

		w.Send([]byte("Got it!"), websocket.DataTypeText)
	})

	err := s.ListenAndServe()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
