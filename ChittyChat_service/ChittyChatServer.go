package main

//SRC: https://github.com/Mai-Sigurd/grpcTimeRequestExample#setting-up-the-server

import (
	"flag"
	ChittyChat_service "github.com/MikkelBKristensen/DSHandins/ChittyChat_service/gRPC"
	"github.com/MikkelBKristensen/DSHandins/HandIn3_ChittyChat/ChittyChat_service/gRPC"
	"google.golang.org/grpc"
	"log"
	"net"
	"strconv"
)

// Server Struct that will be used to represent the Server.
type Server struct {
	ChittyChat_service.UnimplementedChittyChatServer // Necessary
	name                                             string
	port                                             int
	Clock                                            int32
	ClientStreams                                    []gRPC.ChittyChat_ChatServiceServer
	Usernames                                        map[gRPC.ChittyChat_ChatServiceServer]string
}

//@TODO ADD Clock functionality

func (s Server) mustEmbedUnimplementedChittyChatServer() {
	//TODO implement me
	panic("implement me")
}

// Used to get the user-defined port for the server from the command line
var port = flag.Int("port", 0, "server port number")

func startServer(server *Server) {

	// Create a new grpc server
	grpcServer := grpc.NewServer()

	// Make the server listen at the given port (convert int port to string)
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(server.port))

	if err != nil {
		log.Fatalf("Could not create the server %v", err)
	}
	log.Printf("Started server at port: %d\n", server.port)

	// Register the grpc server and serve its listener
	//@TODO Fix the weird server thing pls<3
	ChittyChat_service.RegisterChittyChatServer(grpcServer, server)
	serveError := grpcServer.Serve(listener)
	if serveError != nil {
		log.Fatalf("Could not serve listener")
	}
}

func (s *Server) ChatService(stream gRPC.ChittyChat_ChatServiceServer) error {
	s.ClientStreams = append(s.ClientStreams, stream)
	var RegisteredClient = false
	//The server should keep receiving messages and broadcasting them:
	for {

		clientmessage, err := stream.Recv()
		if err != nil {
			s.endStreamForClient(stream)
			return err
		}

		if !RegisteredClient {
			s.Usernames[stream] = clientmessage.GetUsername()
			RegisteredClient = true
		}
		s.Broadcast(clientmessage)
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

func main() {

	// Get the port from the command line when the server is run
	flag.Parse()

	// Create a server struct
	server := &Server{
		name: "serverName",
		port: *port,
	}

	// Start the server
	go startServer(server)

}
