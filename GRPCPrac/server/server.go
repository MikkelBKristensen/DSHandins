package main

import (
	gRPC "github.com/MikkelBKristensen/DSHandins/GRPCPrac/proto"
)

type Server struct {
	gRPC.UnimplementedHelloServiceServer
}

func main() {

}
