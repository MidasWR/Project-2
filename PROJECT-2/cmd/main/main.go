package main

import (
	"PROJECT-2/cmd/api-server"
	"log"
)

func main() {

	server := api.NewAPI()
	if err, er := server.NewRouter(); err != nil {
		log.Println("error creating server:", err)
		log.Fatal(err)
	} else if er != nil {
		log.Println("error closing db", er)
	}
}
