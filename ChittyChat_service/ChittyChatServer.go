package main

//SRC: https://github.com/Mai-Sigurd/grpcTimeRequestExample#setting-up-the-server

import (
	"github.com/MikkelBKristensen/DSHandins/ChittyChat_service/gRPC"
	"google.golang.org/grpc"
	"log"
	"net"
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
		Message:   "Client: " + username + " has left the chat. \n",
		Timestamp: 1010,
	}
	s.Broadcast(&leaveMessage)

}

func (s *Server) Broadcast(msg *gRPC.Message) {
	//Should broadcast to all clients
	for _, client := range s.ClientStreams {
		if err := client.Send(msg); err != nil {
			//@TODO Implement logging
			//Log errormessage.
		}
	}
	//And log to log file.
}

// Should update the servers clock to be (time of a received message + 1)
func (s *Server) UpdateTime(clientTimestamp int32) int32 {
	s.Clock = max(clientTimestamp, s.Clock) + 1
	return s.Clock
}

// Trying to use homemade max()
func max(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

func main() {
	// Create listener
	listener, err := net.Listen("tcp", ":5001")
	if err != nil {
		log.Fatalf("Could not listen to port: 5001 %v \n", err)
	}
	log.Println("Now listening on port: 5001")

	// Create new server
	grpcServer := grpc.NewServer()

	// Register ChatService
	service := &Server{
		Usernames: make(map[gRPC.ChittyChat_ChatServiceServer]string),
	}

	gRPC.RegisterChittyChatServer(grpcServer, service)

	grpcServer.Serve(listener)
}
