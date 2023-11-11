package main

import (
	"bufio"
	"fmt"
	MeService "github.com/MikkelBKristensen/DSHandins/MutualExclusion/MeService"
	_ "github.com/MikkelBKristensen/DSHandins/MutualExclusion/MeService"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

type Peer struct {
	MeService.UnimplementedMeServiceServer
	port              string
	server            *grpc.Server
	client            MeService.MeServiceClient
	lamportClock      int32
	inCriticalSection bool
	pendingRequests   []MeService.Request
	peerList          map[string]MeService.MeServiceClient
	PortList          []string
}

func NewPeer(port string) *Peer {
	return &Peer{port: port}
}

func (s *Peer) sendJoin(status bool) {
	//Status: 0 for leaving and 1 for joining
	// Update peerList to be equal the PeerPorts file
	// and append own port to the file
	s.peerList, _ = readFile()
	s.updateFile()

	// create a new connectionMsg
	connectionMsg := MeService.ConnectionMsg{
		Port:   s.port,
		IsJoin: status,
	}

	for i := 0; i < len(s.peerList); i++ {
		port := s.peerList[i]
		// TODO Connect to peer
		// TODO Send connectionMsg to peer
	}
}

func (s *Peer) receiveConnectionStatus(message *MeService.ConnectionMsg) {
	// TODO Update clock
	if message.IsJoin == true {
		s.connect(message.GetPort())
		append(s.PortList, message.GetPort())

	} else {
		//Make it delete the port from the PeerList

	}

}

func (s *Peer) leave() {
	s.deleteFromFile()
	panic("Not implemented yet")
}

func readFile() ([]string, error) {
	// Set filename
	fileName := "PeerPorts.txt"

	// Open file
	peerPortFile, err := os.Open(fileName)
	if err != nil {
		log.Panicf("Could not read from data from file: %s", err)
	}
	defer peerPortFile.Close()

	var peerPortArray []string

	// Create scanner to read all lines from the file
	scanner := bufio.NewScanner(peerPortFile)
	for scanner.Scan() {
		port := scanner.Text()
		peerPortArray = append(peerPortArray, port)
	}

	// Check for errors during scanning
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return peerPortArray, nil
}

func (s *Peer) updateFile() {
	// Choose filename
	fileName := "PeerPorts.txt"

	// Open file
	peerPortFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Panicf("Could not open file: %s", err)
	}
	defer peerPortFile.Close()

	// Append port to the file
	if _, err := peerPortFile.WriteString(s.port + "\n"); err != nil {
		log.Fatalf("Could not add Port: %s to file", s.port)
	}

}

func (s *Peer) deleteFromFile() {
	panic("Not implemented yet")
}

func (s *Peer) updatePortList() (err error) {
	s.PortList, err = readFile()
	if err != nil {
		log.Fatalf("Could not read from file: %s", err)
	}
	s.updateFile()

	return nil
}

func (s *Peer) connect(port string) (err error) {
	//Connect to peers
	conn, err := grpc.Dial(":"+port, grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("could not connect to peer on port %s: %v", port, err)
	}
	// add connection to peerList
	s.peerList[port] = MeService.NewMeServiceClient(conn)

	return nil
}

func (s *Peer) StartServer() error {
	// Create listener
	listener, err := net.Listen("tcp", ":"+s.port)
	if err != nil {
		log.Fatalf("Could not listen to port: %s: %v", s.port, err)
	}

	fmt.Println("Now listening on port: ", s.port)

	// Create new server
	s.server = grpc.NewServer()
	MeService.RegisterMeServiceServer(s.server, s)

	// Start server
	go func() {
		if err := s.server.Serve(listener); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	s.updatePortList()
	s.updateFile()

	return nil
}

func main() {
	Peer1 := NewPeer("5001")
	Peer1.StartServer()

}
