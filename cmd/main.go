package main

import (
	"log"
	"os"

	"github.com/JigmeTenzinChogyel/go-net-http-server/cmd/api"
	"github.com/JigmeTenzinChogyel/go-net-http-server/database"
	_ "github.com/lib/pq"
)

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	db, err := database.NewPostgresStorage()
	if err != nil {
		log.Fatal(err)
	}

	server := api.NewAPIServer(":"+port, db)

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}

}
