package main

import (
	"context"
	"flag"
	ChittyChat_service "github.com/MikkelBKristensen/DSHandins/HandIn3_ChittyChat/ChittyChat_service/gRPC"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"strconv"
)

type Client struct {
	id         int
	Timestamp  int
	portNumber int
}

var (
	//Hardcoded
	clientPort = flag.Int("cPort", 0, "client port number")
	serverPort = flag.Int("sPort", 0, "server port number (should match the port used for the server)")
)

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

func Join(ctx context.Context) (reply ChittyChat_service.JoinReply, err error) {
	newClient, err := connectToServer()
	if err != nil {
		log.Fatalf("Client couldn't Join(*Wink wink*) the server.")
	}
	return
}
func Broadcast(ctx context.Context, in *ChittyChat_service.Participant, opts ...grpc.CallOption) (client ChittyChat_service.ChittyChat_BroadcastClient, err error) {

}

func Leave(client ChittyChat_service.ChittyChatClient) (reply ChittyChat_service.JoinReply, err error) {
	return nil, err
}

func Publish(ctx context.Context, req *ChittyChat_service.ChatReq, opts ...grpc.CallOption) (reply *ChittyChat_service.ChatReply, err error) {

	return nil, err
}

func main() {
	// Parse the flags to get the port for the client
	flag.Parse()

	// Create a client
	client, err := connectToServer()
	client.join()

}
