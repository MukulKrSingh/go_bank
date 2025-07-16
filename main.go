package main

import (
	"fmt"
	"log"
)

func main() {
	store, err := NewPostgresStore()

	if err != nil {
		log.Fatal(err)
	}
	err = store.Init()
	if err != nil {
		log.Fatal(err)
	}

	server := NewApiServer(":3000", store)
	server.Run()

	fmt.Print("Bank is alive")

}
