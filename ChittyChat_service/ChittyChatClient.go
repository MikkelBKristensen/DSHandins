package main

import (
	"C"
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	gRPC "github.com/MikkelBKristensen/DSHandins/ChittyChat_service/gRPC"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

/*type Client struct {
	username   string
	Timestamp  int
	portNumber int
	stream     gRPC.ChittyChat_ChatServiceServer
}*/

// We found it easier to work with fields rather than structs as we don't have to send the entire client object around.
var username string
var time int32
var portNumber int
var stream gRPC.ChittyChat_ChatServiceServer

var (
	clientPort = flag.Int("cPort", 0, "client port number")
	serverPort = flag.Int("sPort", 0, "server port number (should match the port used for the server)")
)

func connectToServer() (gRPC.ChittyChatClient, error) {
	// Dial the server at the specified port.
	conn, err := grpc.Dial("localhost:"+strconv.Itoa(*serverPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Client could not connect to port %d", *serverPort)
	} else {
		log.Printf("Client connected to the server at port %d\n", *serverPort)
	}
	return gRPC.NewChittyChatClient(conn), nil
}

func sendMessage(text string, stream gRPC.ChittyChat_ChatServiceClient) {
	err := messageValidation(text)
	if err != nil {
		fmt.Print(err)
	}

	msg := gRPC.Message{
		Username:  username,
		Message:   text,
		Timestamp: time,
	}
	err1 := stream.Send(&msg)
	if err != nil {
		fmt.Print("Could not send, error: ", err1)
	}

}

func receiveMessage(stream gRPC.ChittyChat_ChatServiceClient) {
	for {
		recvMsg, err := stream.Recv()
		if err != nil {
			log.Fatalf("Could not recieve message", err)
			return
		}
		//@TODO Implement time and iteration
		fmt.Print(recvMsg.Username, recvMsg.Message)
	}

}

func messageValidation(msg string) (err error) {
	if len(msg) > 128 {
		err := errors.New("!ERROR! Your message is too long, the maximum is 128 characters")
		return err
	}
	if !utf8.ValidString(msg) {
		err := errors.New("!ERROR! Your message does not comply with UTF8 rules")
		return err
	}
	return
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

	//@TODO Eventuelt brug connectToServer i stedet for følgende
	client := gRPC.NewChittyChatClient(conn)

	stream, err := client.ChatService(context.Background())
	if err != nil {
		log.Fatalf("Could not connect to the ChittyChat service %v", err)
	}

	fmt.Print("Choose your username: ")
	fmt.Scanln(&username)
	fmt.Printf("Username is: " + username)

	//Sending intial join message to the stream
	//@TODO Insert lamport time
	sendMessage("Participant "+username+" joined ChittyChat @ [INSERT LAMPORT TIME]", stream)

	//Listening to messages on the stream
	go receiveMessage(stream)

}
