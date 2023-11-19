package main
import (
	Consensus "github.com/MikkelBKristensen/DSHandins/Handin5/Consensus"
	_ "google.golang.org/grpc"
)

Server struct{
	consensusServer Consensus.ConsensusServer
	consensusClients []Consensus.ConsensusClient

}