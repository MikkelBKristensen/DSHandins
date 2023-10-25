package main

import (
	"flag"
	ChittyChat_service "github.com/MikkelBKristensen/DSHandins/HandIn3_ChittyChat/ChittyChat_service/gRPC"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"strconv"
)

type Client struct {
	id         int
	portNumber int
}

var (
	clientPort = flag.Int("cPort", 0, "client port number")
	serverPort = flag.Int("sPort", 0, "server port number (should match the port used for the server)")
)

func main() {
	// Parse the flags to get the port for the client
	flag.Parse()

	// Create a client
	client := &Client{
		id:         1,
		portNumber: *clientPort,
	}

	// Wait for the client (user) to ask for the time
	go waitForTimeRequest(client)

	for {

	}
}

func connectToServer() (ChittyChat_service.ChittyChatClient, error) {
	// Dial the server at the specified port.
	conn, err := grpc.Dial("localhost:"+strconv.Itoa(*serverPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect to port %d", *serverPort)
	} else {
		log.Printf("Connected to the server at port %d\n", *serverPort)
	}
	return ChittyChat_service.NewChittyChatClient(conn), nil
}

func join()
