package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func main() {

	if err := run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	run(os.Args)
}

func run(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("Expected 2 arguments got %d\n", len(args))
	}
	// Format the port properly(:8080)
	port := fmt.Sprintf(":%s", args[1])

	router := NewRouter()

	router.HandleRoute("/health", func(w ResponseWriter, r *Request) {
		w.Header().Set("Content-Type", "application/json")

		data := ExampleBody{
			Message: "Operation successful",
			Data:    "Gibrish gibrish",
		}

		json.NewEncoder(w).Encode(data)
	})

	s := Server{
		Addr:   port,
		router: router,
	}

	return s.ListenAndServe()

}

type ExampleBody struct {
	Message string `json:"message"`
	Data    string `json:"data"`
}
