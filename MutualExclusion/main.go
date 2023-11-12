package main

import (
	"bufio"
	"context"
	"fmt"
	MeService "github.com/MikkelBKristensen/DSHandins/MutualExclusion/MeService"
	_ "github.com/MikkelBKristensen/DSHandins/MutualExclusion/MeService"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

type Peer struct {
	MeService.UnimplementedMeServiceServer
	port             string
	server           *grpc.Server
	client           MeService.MeServiceClient
	lamportClock     int64
	allowedTimestamp int64
	requestInCS      bool

	// 0 = Released, 1 = Wanted, 2 = Held
	state           int
	pendingRequests []MeService.Request
	peerList        map[string]MeService.MeServiceClient

	// TODO Do not use Portlist, and instead just read directly from file on startup
	PortList []string
}

func NewPeer(port string) *Peer {
	return &Peer{port: port}
}

func (s *Peer) sendConnectionStatus(isJoin bool) {
	//Status: false for leaving and true for joining

	// create a new connectionMsg
	connectionMsg := MeService.ConnectionMsg{
		Port:   s.port,
		IsJoin: isJoin,
	}

	// Send connectionMSG to all peers
	for i := 0; i < len(s.peerList); i++ {
		s.peerList[strconv.Itoa(i)].ConnectionStatus(context.Background(), &connectionMsg)
	}
}

func (s *Peer) receiveConnectionStatus(message *MeService.ConnectionMsg) {
	// TODO Update clock

	if message.IsJoin == true {
		s.connect(message.GetPort())
	} else if message.IsJoin == false {
		delete(s.peerList, message.GetPort())
	}
}

func (s *Peer) leave() {

	s.sendConnectionStatus(false)
	err := s.deleteOwnPortFromFile()
	if err != nil {
		return
	}

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
	log.Printf("Port %s added to file", s.port)
}

func (s *Peer) deleteOwnPortFromFile() error {
	// Choose filename
	fileName := "PeerPorts.txt"

	// Open file
	peerPortFile, err := os.OpenFile(fileName, os.O_RDWR, 0644)
	if err != nil {
		log.Panicf("Could not open file: %s", err)
	}
	defer func(peerPortFile *os.File) {
		err := peerPortFile.Close()
		if err != nil {

		}
	}(peerPortFile)

	// Initialize array to hold the modified content
	var stringsArray []string

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(peerPortFile)
	for scanner.Scan() {
		line := scanner.Text()

		// Check if the line contains the own port
		if !strings.Contains(line, s.port) {
			// If not, append it to the modified array
			stringsArray = append(stringsArray, line)
		}
	}

	// Check for errors during scanning
	if err := scanner.Err(); err != nil {
		return err
	}

	// Truncate the file and write the modified content back
	if err := peerPortFile.Truncate(0); err != nil {
		return err
	}
	// Go to the beginning of the file
	if _, err := peerPortFile.Seek(0, 0); err != nil {
		return err
	}
	// Write the modified content to the file
	writer := bufio.NewWriter(peerPortFile)
	for _, line := range stringsArray {
		fmt.Fprintln(writer, line)
	}
	// Flush write to ensure that all lines are written
	if err := writer.Flush(); err != nil {
		return err
	}
	log.Printf("Port %s removed from file", s.port)
	return nil
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

func (s *Peer) Start() error {
	//This is where the peer starts, first with the server part and then the client part

	// Create listener
	listener, err := net.Listen("tcp", ":"+s.port)
	if err != nil {
		log.Fatalf("Could not listen to port: %s: %v", s.port, err)
	}

	fmt.Println("Now listening on port: ", s.port)

	// Create new server
	s.server = grpc.NewServer()
	MeService.RegisterMeServiceServer(s.server, s.UnimplementedMeServiceServer)

	// Start server
	go func() {
		if err := s.server.Serve(listener); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
	s.updatePortList()
	s.updateFile()
	s.StartClient()

	return nil
}

func (s *Peer) StartClient() error {
	// Connect to all peers
	for i := 0; i < len(s.PortList); i++ {
		err := s.connect(s.PortList[i])
		if err != nil {
			return err
		}
	}
	s.sendConnectionStatus(true)
	return nil
}

func (s *Peer) RequestEntry() {

}

func (s *Peer) RequestEntryRPC(ctx context.Context, req *MeService.Request) (*MeService.Response, error) {
	panic("")
}

func (s *Peer) canEnterCriticalSection() bool {
	panic("")
}

func (s *Peer) enterCriticalSection() {

}

func (s *Peer) ReleaseCriticalSection(ctx context.Context, req *MeService.Request) (*MeService.Empty, error) {
	panic("")
}

func (s *Peer) leaveCriticalSection() {
}

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func main() {
	// Set up the log
	f, err := os.OpenFile("Log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening log: %v", err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)
	log.SetOutput(f)

	// Create three peers
	peer1 := NewPeer("5001")
	peer2 := NewPeer("5002")
	peer3 := NewPeer("5003")

	// Start peers
	go func() {
		err := peer1.Start()
		if err != nil {
			log.Fatalf("Error starting peer1: %v", err)
		}
	}()
	go func() {
		err := peer2.Start()
		if err != nil {
			log.Fatalf("Error starting peer2: %v", err)
		}
	}()
	go func() {
		err := peer3.Start()
		if err != nil {
			log.Fatalf("Error starting peer3: %v", err)
		}
	}()

	// Wait for the demonstration to complete
	for {
	}
}
