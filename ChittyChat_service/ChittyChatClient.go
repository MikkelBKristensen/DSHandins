package main

import (
	"C"
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/MikkelBKristensen/DSHandins/ChittyChat_service/gRPC"
	"google.golang.org/grpc"
	"log"
	"os"
	"strings"
	"unicode/utf8"
)

/*type Client struct {
	username   string
	Timestamp  int
	portNumber string
	stream     gRPC.ChittyChat_ChatServiceServer
}*/

// We found it easier to work with fields rather than structs as we don't have to send the entire client object around.
var username string
var time int32
var portNumber string
var stream gRPC.ChittyChat_ChatServiceServer

func sendMessage(text string, stream gRPC.ChittyChat_ChatServiceClient) {
	err := messageValidation(text)
	if err != nil {
		fmt.Print(err)
	}

	msg := gRPC.Message{
		Username:  username,
		Message:   text,
		Timestamp: 2,
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
			log.Fatalf("Could not receive message: %v", err)
			return
		}
		//@TODO Implement time and iteration
		fmt.Printf(recvMsg.Username + ": " + recvMsg.Message + "\n")
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
	fmt.Print("Choose your username: ")
	fmt.Scanln(&username)
	fmt.Printf("Username is: " + username + " \n")

	// Try to connect to the gRPC server
	conn, err := grpc.Dial("localhost:5001", grpc.WithInsecure())
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

	//Sending intial join message to the stream
	//@TODO Insert lamport time
	sendMessage("Participant "+username+" joined ChittyChat @ [INSERT LAMPORT TIME] \n", stream)

	//Listening to messages on the stream
	go receiveMessage(stream)

	messageReader := bufio.NewReader(os.Stdin)
	for {
		text, _ := messageReader.ReadString('\n')
		sendMessage(strings.TrimSpace(text), stream)
	}

}
