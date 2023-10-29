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

//@TODO ADD Clock functionality

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

		if !RegisteredClient {
			s.Usernames[stream] = clientMessage.GetUsername()
			RegisteredClient = true
		}
		clientMessage.Timestamp = s.UpdateTime(clientMessage.Timestamp)
		s.Broadcast(clientMessage)
		//Check if Client already has a stream
	}
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
		Message:   "Client: " + username + " has left the chat.",
		Timestamp: s.Clock,
	}
	s.Broadcast(&leaveMessage)

}

func (s *Server) Broadcast(msg *gRPC.Message) {
	log.Printf("[SERVER - OUTGOING] "+msg.Username+" has left the chat @ Lamport time ", s.Clock)

	//Should broadcast to all clients
	for _, client := range s.ClientStreams {
		if err := client.Send(msg); err != nil {
			log.Printf("Could not broadcast message: %v", err)
		}
	}
}

// UpdateTime Update server clock if msg timestamp is greater than servers time, otherwise dont update
func (s *Server) UpdateTime(clientTimestamp int32) int32 {
	if clientTimestamp > s.Clock {
		s.Clock = clientTimestamp + 1
	}
	return s.Clock
}

func main() {
	// Set up the log
	f, err := os.OpenFile("logfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	// Create listener
	listener, err := net.Listen("tcp", ":5001")
	if err != nil {
		log.Fatalf("Could not listen to port: 5001 %v \n", err)
	}
	log.Println("Now listening on port: 5001")
	fmt.Println("Now listening on port: 5001")

	// Create new server
	grpcServer := grpc.NewServer()

	// Register ChatService
	service := &Server{
		Usernames: make(map[gRPC.ChittyChat_ChatServiceServer]string),
	}

	gRPC.RegisterChittyChatServer(grpcServer, service)

	grpcServer.Serve(listener)
}
