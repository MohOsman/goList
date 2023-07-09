package main

import (
	"goList/api"
	"log"
	
)

func main() {

    server := api.NewServer(":8080")
	log.Fatal(server.Start())


}

// create read update and delete tasks


