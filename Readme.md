# BeaconActor - WIP

This golang program tracks the performance of one or more validators on the ethereum mainnet beacon chain. It uses an Actor Design pattern to allow for the parallel processing of data for multiple validators. The program stores validator performance data in a postgresql database.

## Introduction

This program allows you to provide one or more validators whose performamce you wish to track. Each validators performance is tracked in almost real time - by updating the performance data for the designated validators each epoch (approx every 6.2 minutes)

## Features
1. HTTP endpoints for adding and removing validators to be tracked.
2. Each epoch an actor is spawned for each validator to update the performance data for that validator.
3. This data is stored in a postgresql database, with the performance data accessible via a HTTP endpoints to a Front End application.

## Getting Started

### Prerequisites
- To use this program you will need to have the following dependencies inplace:
  - Access to a remote or local postgresql database
  - A ethereum beaconchain node that exposes the beaconchain api.
- golang 1.21 or higher

### Installation
- Once these dependencies are in place you will need to copy the provided .env.example file to a .env file and update the values to match your environment - specifically, the urls for the beaconchain node and the postgresql database.
- Prior to running the program you will need to create the required database `beaconvalidators` and run the following command to create the required database schema.

```bash
make created-db
migrate -path ./migrations -database <database-connection-url> up
```


