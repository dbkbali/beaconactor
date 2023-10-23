package main

import (
	"context"
	"fmt"
	"os"

	"beacon-actor/server"

	"beacon-actor/datastore"

	eth2client "github.com/attestantio/go-eth2-client"
	"github.com/attestantio/go-eth2-client/http"
	"github.com/rs/zerolog"
)

var db *datastore.DB
var eth2HttpClient *eth2client.Service

func main() {
	// startup database
	ctx := context.Background()
	var err error
	db, err = datastore.NewDb(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Println("unable to connect to database: %w", err)
		os.Exit(1)
	}
	db.Ping(ctx)
	dbName := "validators"
	res, err := db.DbExists(ctx, dbName)
	if err != nil {
		fmt.Println("unable to query database: %w", err)
		os.Exit(1)
	}
	if !res {
		fmt.Printf("database %s does not exist.. creating", dbName)
		createDb := fmt.Sprintf("CREATE DATABASE %s;", dbName)
		_, err = db.Exec(ctx, createDb)
		if err != nil {
			fmt.Println("unable to create database: %w", err)
			os.Exit(1)
		}
	}
	beaconUrl := os.Getenv("BEACON_NODE_URL")
	client := connectToBeaconNode(beaconUrl)
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
		http.WithAddress(beaconNodeUrl),
		http.WithLogLevel(zerolog.DebugLevel),
	)
	if err != nil {
		fmt.Println("unable to connect to beacon node: %w", err)
		os.Exit(1)
	}

	fmt.Printf("Connected to %s\n", client.Name())
	return client
	// get genesis
	// if provider, isProvider := client.(eth2client.GenesisProvider); isProvider {
	// 	genesisResponse, err := provider.Genesis(ctx)
	// 	if err != nil {
	// 		// Errors may be API errors, in which case they will have more detail
	// 		// about the failure.
	// 		// var apiErr *api.Error
	// 		// if errors.As(err, &apiErr) {
	// 		// 	switch apiErr.StatusCode {
	// 		// 	case 404:
	// 		// 		panic("genesis not found")
	// 		// 	case 503:
	// 		// 		panic("node is syncing")
	// 		// 	}
	// 		// }
	// 		panic(err)
	// 	}
	// 	fmt.Printf("Genesis time is %v\n", genesisResponse.GenesisTime)
	// 	cancel()
}
