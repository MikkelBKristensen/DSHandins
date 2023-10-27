package main

//SRC: https://github.com/Mai-Sigurd/grpcTimeRequestExample#setting-up-the-server

import (
	"container/list"
	"flag"
	ChittyChat_service "github.com/MikkelBKristensen/DSHandins/HandIn3_ChittyChat/ChittyChat_service/gRPC"
	"google.golang.org/grpc"
	"log"
	"net"
	"strconv"
)

// Server Struct that will be used to represent the Server.
type Server struct {
	/*
		ChittyChat_service.UnimplementedChittyChatServer // Necessary
		name                                             string
		port                                             int
		streamList 										 list
	*/

}

func (s Server) ChatService(server ChittyChat_service.ChittyChat_ChatServiceServer) error {
	//TODO implement me
	panic("implement me")
}

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
	streamList = list.New()

	// Register the grpc server and serve its listener
	ChittyChat_service.RegisterChittyChatServer(grpcServer, server)
	serveError := grpcServer.Serve(listener)
	if serveError != nil {
		log.Fatalf("Could not serve listener")
	}
}

func ChatService(stream ChittyChat_service.ChittyChat_ChatServiceServer) error {
	for {
		// Receive a message from the client
		clientmessage, err := stream.Recv()
		if err != nil {
			return err
		}

		/*---- Process the received message (msg) ----*/

		//Leave stream functionality
		if clientmessage.Message == "Leave" {
			response := &ChittyChat_service.ServerMessage{
				Username:  clientmessage.Username,
				Message:   clientmessage.Username + ": Left broadcast @ " + strconv.Itoa(int(clientmessage.Timestamp)),
				Timestamp: 1010,
			}
			if err := stream.Send(response); err != nil {
				return err
			}
		}
		//Join stream functionality
		if clientmessage.Message == "Join" {
			response := &ChittyChat_service.ServerMessage{
				Username:  clientmessage.Username,
				Message:   clientmessage.Username + ": Joined broadcast @ " + strconv.Itoa(int(clientmessage.Timestamp)),
				Timestamp: 1010,
			}
			streamList.PushFront(stream)
			if err := stream.Send(response); err != nil {
				return err
			}
		}

		// Send a response back to the client
		response := &ChittyChat_service.ServerMessage{
			Username:  clientmessage.Username,
			Message:   clientmessage.Username + ": " + clientmessage.Message + " @ " + strconv.Itoa(int(clientmessage.Timestamp)),
			Timestamp: 1011,
		}

		if err := stream.Send(response); err != nil {
			return err
		}
	}
}

func Send(server ChittyChat_service.ChittyChat_ChatServiceServer) {

}

func Recv() {

}

/*
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
*/
