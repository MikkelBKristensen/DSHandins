package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/MikkelBKristensen/DSHandins/ChittyChat_service/gRPC"
	ChittyChat_service "github.com/MikkelBKristensen/DSHandins/ChittyChat_service/gRPC"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

type Client struct {
	id         int
	Timestamp  int
	portNumber int
	stream     ChittyChat_service.ChittyChat_ChatServiceServer
}

var (
	clientPort = flag.Int("cPort", 0, "client port number")
	serverPort = flag.Int("sPort", 0, "server port number (should match the port used for the server)")
)

func connectToServer() (ChittyChat_service.ChittyChatClient, error) {
	// Dial the server at the specified port.
	conn, err := grpc.Dial("localhost:"+strconv.Itoa(*serverPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Client could not connect to port %d", *serverPort)
	} else {
		log.Printf("Client connected to the server at port %d\n", *serverPort)
	}
	return ChittyChat_service.NewChittyChatClient(conn), nil
}

func ChatService(ctx context.Context, opts ...grpc.CallOption) (ChittyChat_ChatServiceClient, error) {

}
func sendMessage() {

}

func recieveMessage() {

}
func messageValidation(msg gRPC.Message) error {
	if len(msg.Message) > 128 {
		err := errors.New("!error! Your message is too long, the maximum is 128 characters")
		return err
	}
	if !utf8.ValidString(msg.Message) {
		err := errors.New("!error! Your message does not comply with UTF8 rules")
		return err
	}
	return nil
}

func main() {
	// Parse the flags to get the port for the client
	flag.Parse()

	// Set client port
	fmt.Println("Enter Server Port:")
	reader := bufio.NewReader(os.Stdin)
	serverPort, err := reader.ReadString('\n')

	if err != nil {
		log.Printf("Could not read input from terminal %v", err)
	}
	serverPort = strings.Trim(serverPort, "\r\n")

	log.Println("Connecting to port: " + serverPort)

	// Try to connect to the gRPC server
	conn, err := grpc.Dial(serverPort, grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Failed to connect to gRPC server %v", err)
	}
	defer conn.Close()

	// Eventuelt brug connectToServer i stedet for følgende
	client := ChittyChat_service.NewChittyChatClient(conn)

	stream, err := client.ChatService(context.Background())
	if err != nil {
		log.Fatalf("Could not connect to the ChittyChat service %v", err)
	}

	client := Client{}

}
