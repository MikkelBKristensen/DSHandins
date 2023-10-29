package main

//SRC: https://github.com/Mai-Sigurd/grpcTimeRequestExample#setting-up-the-server

import (
	"fmt"
	"github.com/MikkelBKristensen/DSHandins/ChittyChat_service/gRPC"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

// Server Struct that will be used to represent the Server.
type Server struct {
	gRPC.UnimplementedChittyChatServer // Necessary
	Clock                              int32
	ClientStreams                      []gRPC.ChittyChat_ChatServiceServer
	Usernames                          map[gRPC.ChittyChat_ChatServiceServer]string
}

func (s *Server) ChatService(stream gRPC.ChittyChat_ChatServiceServer) error {
	s.ClientStreams = append(s.ClientStreams, stream)

	var RegisteredClient = false

	//The server should keep receiving messages and broadcasting them:
	for {

		clientMessage, err := stream.Recv()
		if err != nil {
			s.endStreamForClient(stream)
			return err
		}

		//Check if Client already has a stream
		if !RegisteredClient {
			s.Usernames[stream] = clientMessage.GetUsername()
			RegisteredClient = true
			clientMessage.Username = "Server"
		}

		// If the received timestamp is less than the servers clock,
		// then set the clients clock equal to the servers clock to maintain order
		if clientMessage.Timestamp < s.Clock {
			clientMessage.Timestamp = s.Clock
		}

		// Increment server clock
		s.Clock = max(clientMessage.Timestamp, s.Clock) + 1

		s.Broadcast(clientMessage)
	}
}

func (s *Server) RegisterUsername(stream gRPC.ChittyChat_ChatServiceServer, username string) string {
	return username
}

func (s *Server) endStreamForClient(targetClient gRPC.ChittyChat_ChatServiceServer) {
	//Locate  the specific stream that needs to be closed:
	for i, client := range s.ClientStreams {
		if client == targetClient {
			s.ClientStreams = append(s.ClientStreams[:i], s.ClientStreams[i+1:]...)
			break
		}
	}

	username := s.Usernames[targetClient]

	var leaveMessage = gRPC.Message{
		Username:  "Server",
		Message:   username + " has left the chat.",
		Timestamp: s.Clock,
	}
	s.Broadcast(&leaveMessage)
	//log.Printf("[SERVER]: %s has left the chat @ Lamport time %d", username, s.Clock)
}

func (s *Server) Broadcast(msg *gRPC.Message) {
	//Should broadcast to all clients
	for _, client := range s.ClientStreams {
		if err := client.Send(msg); err != nil {
			log.Printf("[SERVER]: Could not broadcast message: %v", err)
		}
	}
}

// Join message
func (s *Server) GetJoinMessage(username string) gRPC.Message {
	var message = gRPC.Message{
		Username:  "Server",
		Message:   "Participant " + username + " joined the chat!",
		Timestamp: s.Clock,
	}
	return message
}

func main() {
	// Set up the log
	f, err := os.OpenFile("logfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("[SERVER]: error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	// Create listener
	listener, err := net.Listen("tcp", ":5001")
	if err != nil {
		log.Fatalf("[SERVER]: Could not listen to port: 5001 %v \n", err)
	}
	log.Println("[SERVER]: Now listening on port: 5001")
	fmt.Println("Now listening on port: 5001")

	// Create new server
	grpcServer := grpc.NewServer()

	// Register ChatService
	service := &Server{
		Usernames: make(map[gRPC.ChittyChat_ChatServiceServer]string),
		Clock:     0,
	}

	gRPC.RegisterChittyChatServer(grpcServer, service)

	grpcServer.Serve(listener)
}
