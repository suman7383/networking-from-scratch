package main

import (
	"fmt"
	"os"

	"github.com/suman7383/networking-from-scratch/websocket-server/cmd/server"
)

func main() {
	s := server.NewServer(":8080")

	err := s.ListenAndServe()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
