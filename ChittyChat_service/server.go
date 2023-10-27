package main

import (
	"github.com/MikkelBKristensen/DSHandins/HandIn3_ChittyChat/ChittyChat_service/ChittyChatServer"
	ChittyChat_service "github.com/MikkelBKristensen/DSHandins/HandIn3_ChittyChat/ChittyChat_service/gRPC"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {

	// Set port number
	Port := "5000"

	// Create listener
	listener, err := net.Listen("tcp", ":"+Port)
	if err != nil {
		log.Fatalf("Could not listen to port: ", Port, err)
	}
	log.Println("Now listening on port:" + Port)

	grpcServer := grpc.NewServer()

	chatServer := ChittyChatServer.Server{}
	ChittyChat_service.RegisterChittyChatServer(grpcServer, chatServer)

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("Could not start saerver", err)
	}
}
