package main

import (
	"fmt"
	"os"
)

func main() {
	_, err := ListenAndServe(":8080")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
