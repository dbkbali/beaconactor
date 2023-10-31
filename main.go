package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	server "beacon-actor/api"
	"beacon-actor/config"
	"beacon-actor/datastore/postgresql"

	eth2client "github.com/dbkbali/go-eth2-client"
	"github.com/dbkbali/go-eth2-client/http"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

var db *postgresql.DB
var eth2HttpClient *eth2client.Service

func init() {
	// load .env file
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
		os.Exit(1)
	}
}

func main() {
	config := config.New()
	// startup database
	ctx := context.Background()
	var err error
	db, err = postgresql.NewDb(ctx, config.DatabaseUrl)
	if err != nil {
		fmt.Println("unable to connect to database: %w", err)
		os.Exit(1)
	}
	db.Ping(ctx)
	client := connectToBeaconNode(config.BeaconNodeUrl)
	eth2HttpClient = &client
	fmt.Println("starting server")
	server := server.NewServer(&server.ServerConfig{
		ListenAddress: ":8080",
	}, db, eth2HttpClient)
	server.Start()
	fmt.Printf("Server started listening on port %s\n", server.ListenAddress)

}

func connectToBeaconNode(beaconNodeUrl string) (client eth2client.Service) {
	fmt.Printf("connecting to beacon node at %s\n", beaconNodeUrl)

	client, err := http.New(context.Background(),
		http.WithTimeout(15*time.Second),
		http.WithAddress(beaconNodeUrl),
		http.WithLogLevel(zerolog.DebugLevel),
	)
	if err != nil {
		fmt.Println("unable to connect to beacon node: %w", err)
		os.Exit(1)
	}

	fmt.Printf("Connected to %s\n", client.Name())
	return client
}
