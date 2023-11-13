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
	"strings"
	"sync"
	"time"
)

type Peer struct {
	MeService.UnimplementedMeServiceServer
	port         string
	server       *grpc.Server
	client       MeService.MeServiceClient
	lamportClock int64

	// Is the value of the timestamp sent with an entryRequest
	allowedTimestamp int64

	// 0 = Released, 1 = Wanted, 2 = Held
	state           int
	pendingRequests []*MeService.Message
	peerList        map[string]MeService.MeServiceClient

	// TODO Do not use Portlist, and instead just read directly from file on startup
	PortList []string
}

func NewPeer(port string) *Peer {
	return &Peer{
		port:             port,
		lamportClock:     1,
		allowedTimestamp: 1,
		state:            0,
		peerList:         make(map[string]MeService.MeServiceClient),
	}
}

func (p *Peer) sendConnectionStatus(isJoin bool) {
	//Status: false for leave message and true for join message
	p.lamportClock++

	connectionMsg := MeService.ConnectionMsg{
		Port:      p.port,
		IsJoin:    isJoin,
		Timestamp: p.lamportClock,
	}

	// Send connectionMSG to all peers
	if len(p.peerList) != 0 {
		//Send connection status to all peers by iterating over map

		for k := range p.peerList {
			target := p.peerList[k]
			response, err := target.ConnectionStatus(context.Background(), &connectionMsg)
			if err != nil {
				log.Fatalf("Peer %s Could not send connection status to peer on port %s: %v", p.port, k, err)
			}

			p.pickMaxAndUpdateClock(response.Timestamp)
		}
	}

}

func (p *Peer) ConnectionStatus(ctx context.Context, inComming *MeService.ConnectionMsg) (msg *MeService.Message, err error) {
	// TODO Update clock

	if inComming.IsJoin == true {

		err := p.connect(inComming.GetPort())
		if err != nil {
			return nil, fmt.Errorf("Peer %s could not connect to peer port %s", p.port, inComming.Port)
		}
		p.pickMaxAndUpdateClock(inComming.Timestamp)

		log.Printf("Peer %s noticed that Peer %s joined @ lamport time %d", p.port, inComming.GetPort(), p.lamportClock)
		return p.returnMessage(), fmt.Errorf("")

	} else if inComming.IsJoin == false {
		//TODO Address that peer logs the time of notis and not the time of the actual leave msg

		delete(p.peerList, inComming.GetPort())
		p.pickMaxAndUpdateClock(inComming.Timestamp)
		log.Printf("Peer %s noticed that Peer %s left @ lamport time %d", p.port, inComming.GetPort(), p.lamportClock)
		return nil, err
	}
	return
}

func (p *Peer) leave() {

	p.sendConnectionStatus(false)
	err := p.deleteOwnPortFromFile()
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

func (p *Peer) updateFile() {
	// Choose filename
	fileName := "PeerPorts.txt"

	// Open file
	peerPortFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Panicf("Could not open file: %s", err)
	}
	defer peerPortFile.Close()

	// Append port to the file
	if _, err := peerPortFile.WriteString(p.port + "\n"); err != nil {
		log.Fatalf("Could not add Port: %s to file", p.port)
	}
	log.Printf("Port %s added to file", p.port)
}

func (p *Peer) deleteOwnPortFromFile() error {
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
		if !strings.Contains(line, p.port) {
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
	log.Printf("Port %s removed from file", p.port)
	return nil
}

func (p *Peer) updatePortList() (err error) {
	p.PortList, err = readFile()
	if err != nil {
		log.Fatalf("Could not read from file: %s", err)
	}
	p.updateFile()

	return nil
}

func (p *Peer) connect(port string) (err error) {
	//Connect to peers
	log.Printf("Peer %s is connecting to peer %s", p.port, port)
	conn, err := grpc.Dial(":"+port, grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("could not connect to peer on port %s: %v", port, err)
	}
	// add connection to peerList
	p.peerList[port] = MeService.NewMeServiceClient(conn)

	return nil
}

func (p *Peer) Start() error {
	//This is where the peer starts, first with the server part and then the client part

	// Create listener
	listener, err := net.Listen("tcp", ":"+p.port)
	if err != nil {
		log.Fatalf("Could not listen to port: %s: %v", p.port, err)
	}

	fmt.Println("Now listening on port: ", p.port)

	// Create new server
	p.server = grpc.NewServer()
	MeService.RegisterMeServiceServer(p.server, p.UnimplementedMeServiceServer)

	// Start server
	go func() {
		if err := p.server.Serve(listener); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
	err = p.updatePortList()
	if err != nil {
		return err
	}

	err = p.StartClient()
	if err != nil {
		return err
	}

	return nil
}

func (p *Peer) StartClient() error {
	//Setting up the peers client
	conn, err := grpc.Dial(":"+p.port, grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("could not connect to peer on port %s: %v", p.port, err)
	}
	p.client = MeService.NewMeServiceClient(conn)

	// Connect to all peers
	for i := 0; i < len(p.PortList); i++ {
		err := p.connect(p.PortList[i])
		if err != nil {
			return err
		}
	}
	p.sendConnectionStatus(true)
	return nil
}

// MakeRequest This is the Client part of the peer, where it sends requests
func (p *Peer) MakeRequest() {
	p.state = 1
	p.lamportClock++
	p.allowedTimestamp = p.lamportClock

	var wg sync.WaitGroup

	for _, client := range p.peerList {
		// Increment the WaitGroup counter for each goroutine
		wg.Add(1)

		// Make gRPC call in a goroutine
		go func(client MeService.MeServiceClient) {
			// Decrement the WaitGroup counter when the goroutine finishes
			defer wg.Done()

			// Make the gRPC call
			_, _ = client.RequestEntry(context.Background(), &MeService.Message{
				Timestamp: p.lamportClock,
				NodeId:    p.port,
			})
		}(client)
	}

	// Wait for all goroutines to finish before returning
	wg.Wait()

	// Enter critical section
	p.enterCriticalSection()
	// Leave critical section
	p.leaveCriticalSection()
}

// RequestEntry This is the server part of the peer, where it handles how to return the actual rpc method
func (p *Peer) RequestEntry(ctx context.Context, entryRequest *MeService.Message) (*MeService.Message, error) {

	if p.state == 2 || (p.state == 1 && p.allowedTimestamp < entryRequest.Timestamp) ||
		(p.state == 1 && p.allowedTimestamp == entryRequest.Timestamp && p.port < entryRequest.NodeId) {

		p.pickMaxAndUpdateClock(entryRequest.Timestamp)
		p.pendingRequests = append(p.pendingRequests, entryRequest)

		// Infinite loop while waiting for the critical section to be released
		for p.state == 2 {
		}

		return p.returnMessage(), nil

	} else if (p.state == 1 && entryRequest.Timestamp < p.allowedTimestamp) ||
		(p.state == 1 && p.allowedTimestamp == entryRequest.Timestamp && p.port >= entryRequest.NodeId) {

		p.pickMaxAndUpdateClock(entryRequest.Timestamp)

		return p.returnMessage(), nil

	} else {
		p.pickMaxAndUpdateClock(entryRequest.Timestamp)
		return p.returnMessage(), nil
	}
}
func (p *Peer) returnMessage() *MeService.Message {
	return &MeService.Message{
		Timestamp: p.lamportClock,
		NodeId:    p.port,
	}
}
func (p *Peer) pickMaxAndUpdateClock(requestTimeStamp int64) {
	p.lamportClock = maxL(p.lamportClock, requestTimeStamp)
	p.lamportClock++
}
func (p *Peer) enterCriticalSection() {
	p.lamportClock++
	p.state = 2
	log.Printf("Peer %s enters critical section at lamport time %d", p.port, p.lamportClock)

	// wait 5 seconds before leaving
	time.Sleep(5 * time.Second)
}
func (p *Peer) leaveCriticalSection() {
	p.lamportClock++
	p.state = 0
	log.Printf("Peer %s leaves critical section at lamport time %d", p.port, p.lamportClock)
}
func maxL(a, b int64) int64 {
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
	err = peer1.Start()
	if err != nil {
		log.Fatalf("Error starting peer1: %v", err)
	}
	err = peer2.Start()
	if err != nil {
		log.Fatalf("Error starting peer2: %v", err)
	}
	err = peer3.Start()
	if err != nil {
		log.Fatalf("Error starting peer3: %v", err)
	}

	time.Sleep(2)

	peer1.MakeRequest()
	time.Sleep(1)
	peer2.MakeRequest()
	time.Sleep(1)
	peer3.MakeRequest()
	time.Sleep(1)
	peer1.leave()
	time.Sleep(1)
	peer2.leave()
	time.Sleep(1)
	peer3.leave()

	// Wait for the demonstration to complete
	for {
	}
}
