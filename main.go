package main

import (
	"fmt"
	"log"
)

func main() {
	store, err := newPostgressStore()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", store)
	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	server := newApiServer(":3001", store)
	server.Run()
}
