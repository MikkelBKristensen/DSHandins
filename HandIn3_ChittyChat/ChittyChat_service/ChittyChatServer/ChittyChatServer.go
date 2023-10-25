package ChittyChatServer

//SRC: https://github.com/Mai-Sigurd/grpcTimeRequestExample#setting-up-the-server
import (
	"context"
	"errors"
	"flag"
	"fmt"
	ChittyChat_service "github.com/MikkelBKristensen/DSHandins/HandIn3_ChittyChat/ChittyChat_service/gRPC"
	"google.golang.org/grpc"
	"log"
	"net"
	"strconv"
)

// Struct that will be used to represent the Server.
type Server struct {
	ChittyChat_service.UnimplementedChittyChatServer // Necessary
	name                                             string
	port                                             int
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
	ChittyChat_service.RegisterChittyChatServer(grpcServer, server)
	serveError := grpcServer.Serve(listener)
	if serveError != nil {
		log.Fatalf("Could not serve listener")
	}
}

func (s *Server) Publish(ctx context.Context, req *ChittyChat_service.ChatReq) (*ChittyChat_service.ChatReply, error) {
	if len(req.Text) > 128 {
		return nil, errors.New("message is to long to be published")
	}
	return &ChittyChat_service.ChatReply{Text: fmt.Sprintf("Received: %s", req.Text)}, nil

}

func (s *Server) Broadcast(stream ChittyChat_service.ChittyChat_BroadcastServer) {
	// Implement your Broadcast logic
	// This method should handle streaming messages from the client

}

func (s *Server) Join(ctx context.Context, reply *ChittyChat_service.Participant) (*ChittyChat_service.JoinReply, error) {
	// Implement your Join logic
	// This method should handle joining participants
	return nil, nil
}

func (s *Server) Leave(ctx context.Context, reply *ChittyChat_service.Participant) (*ChittyChat_service.JoinReply, error) {
	// Implement your Leave logic
	// This method should handle leaving participants
	return nil, nil
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
