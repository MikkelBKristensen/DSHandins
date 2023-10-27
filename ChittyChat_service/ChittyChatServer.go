package main

//SRC: https://github.com/Mai-Sigurd/grpcTimeRequestExample#setting-up-the-server

import (
	"flag"
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

func (s Server) mustEmbedUnimplementedChittyChatServer() {
	//TODO implement me
	panic("implement me")
}

// Used to get the user-defined port for the server from the command line
var port = flag.Int("port", 0, "server port number")

/*func startServer(server *Server) {

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
	//chatServer :=
	gRPC.RegisterChittyChatServer(grpcServer, chatServer)
	serveError := grpcServer.Serve(listener)
	if serveError != nil {
		log.Fatalf("Could not serve listener")
	}
}*/

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

	// Set port number
	Port := "5001"

	// Create listener
	listener, err := net.Listen("tcp", ":5001")
	if err != nil {
		log.Fatalf("Could not listen to port: %d %v", Port, err)
	}
	log.Println("Now listening on port:" + Port)

	// Create new server
	grpcServer := grpc.NewServer()

	// Register ChatService
	service := &Server{
		Usernames: make(map[gRPC.ChittyChat_ChatServiceServer]string),
	}

	gRPC.RegisterChittyChatServer(grpcServer, service)

	grpcServer.Serve(listener)
	/*err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("Could not start server", err)
	}*/
}
