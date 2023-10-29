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

// We found it easier to work with fields rather than structs as we don't have to send the entire client object around.
var username string
var time int32 = 0
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
			log.Fatalf("Could not receive message: %v", err)
			return
		}

		// write received message to log
		log.Printf("[TO %s, FROM %s] : %s @ lamport time %d", username, recvMsg.Username, recvMsg.Message, recvMsg.GetTimestamp())

		// print in the client's terminal
		fmt.Printf("%s: %s\n", recvMsg.Username, recvMsg.Message)

		// Update time
		if recvMsg.Timestamp > time {
			time = recvMsg.Timestamp + 1
		}
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
	// Choose displayed username
	fmt.Print("Choose your username: ")
	fmt.Scanln(&username)
	fmt.Printf("Username is: " + username + " \n")

	// Set up the log
	f, err := os.OpenFile("logfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	// Try to connect to the gRPC server
	conn, err := grpc.Dial("localhost:5001", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server %v", err)
	}
	defer conn.Close()

	// Create client on server
	client := gRPC.NewChittyChatClient(conn)

	stream, err := client.ChatService(context.Background())
	if err != nil {
		log.Fatalf("Could not connect to the ChittyChat service %v", err)
	}

	//Sending intial join message to the stream
	joinString := fmt.Sprintf("Partipant %s joined the chat", username)
	sendMessage(joinString, stream)

	//Listening to messages on the stream
	go receiveMessage(stream)

	// Lets client write messages in terminal
	messageReader := bufio.NewReader(os.Stdin)
	for {
		text, _ := messageReader.ReadString('\n')
		sendMessage(strings.TrimSpace(text), stream)
	}

}
