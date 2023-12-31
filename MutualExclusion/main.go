// main.go

package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	MeService "github.com/MikkelBKristensen/DSHandins/MutualExclusion/MeService"
	"google.golang.org/grpc"
)

// Peer represents a participant in the mutual exclusion protocol.
type Peer struct {
	port             string
	meServer         MeService.MeServiceServer
	server           *grpc.Server
	lamportClock     int64
	allowedTimestamp int64
	state            int
	peerList         map[string]MeService.MeServiceClient
	PortList         []string
}
type MeServiceServer struct {
	MeService.UnimplementedMeServiceServer
	Peer *Peer
}

// NewPeer creates a new instance of the Peer struct.
func NewPeer(port string) *Peer {
	meServer := &MeServiceServer{
		Peer: &Peer{
			port:             port,
			lamportClock:     1,
			allowedTimestamp: 1,
			state:            0,
			peerList:         make(map[string]MeService.MeServiceClient),
		},
	}
	meServer.Peer.meServer = meServer
	return meServer.Peer
}

// Server methods

func (p *Peer) sendConnectionStatus(isJoin bool) {

	connectionMsg := MeService.ConnectionMsg{
		Port:      p.port,
		IsJoin:    isJoin,
		Timestamp: p.lamportClock,
	}

	if len(p.peerList) != 0 {
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

func (s *MeServiceServer) ConnectionStatus(_ context.Context, inComing *MeService.ConnectionMsg) (*MeService.Message, error) {
	p := s.Peer
	if p == nil {
		// Handle the case when Peer is nil
		return nil, fmt.Errorf("peer is nil")
	}
	if inComing.IsJoin == true {
		err := p.connect(inComing.GetPort())
		if err != nil {
			return nil, fmt.Errorf("peer %s could not connect to peer port %s", p.port, inComing.Port)
		}
		p.pickMaxAndUpdateClock(inComing.Timestamp)

		log.Printf("Peer %s noticed that Peer %s joined @ lamport time %d", p.port, inComing.GetPort(), p.lamportClock)
		return p.returnMessage(), nil
	} else {
		delete(p.peerList, inComing.GetPort())
		p.pickMaxAndUpdateClock(inComing.Timestamp)
		log.Printf("Peer %s noticed that Peer %s left @ lamport time %d", p.port, inComing.GetPort(), p.lamportClock)
		return nil, nil
	}
}

// File methods

func readFile(fileName string) ([]string, error) {
	peerPortFile, err := os.Open(fileName)
	if err != nil {
		log.Panicf("Could not read from data from file: %s", err)
	}
	defer func(peerPortFile *os.File) {
		err := peerPortFile.Close()
		if err != nil {
			_ = fmt.Errorf("error closing file: %v", err)
		}
	}(peerPortFile)

	var peerPortArray []string
	scanner := bufio.NewScanner(peerPortFile)
	for scanner.Scan() {
		port := scanner.Text()
		peerPortArray = append(peerPortArray, port)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return peerPortArray, nil
}

func (p *Peer) updateFile(fileName string) {
	peerPortFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Panicf("Could not open file: %s", err)
	}
	defer func(peerPortFile *os.File) {
		err := peerPortFile.Close()
		if err != nil {
			_ = fmt.Errorf("error closing file: %v", err)
		}
	}(peerPortFile)

	if _, err := peerPortFile.WriteString(p.port + "\n"); err != nil {
		log.Fatalf("Could not add Port: %s to file", p.port)
	}
	log.Printf("Port %s added to file", p.port)
}

func (p *Peer) deleteOwnPortFromFile(fileName string) error {
	peerPortFile, err := os.OpenFile(fileName, os.O_RDWR, 0644)
	if err != nil {
		log.Panicf("Could not open file: %s", err)
	}
	defer func(peerPortFile *os.File) {
		err := peerPortFile.Close()
		if err != nil {
			_ = fmt.Errorf("error closing file: %v", err)
		}
	}(peerPortFile)

	var stringsArray []string
	scanner := bufio.NewScanner(peerPortFile)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.Contains(line, p.port) {
			stringsArray = append(stringsArray, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	if err := peerPortFile.Truncate(0); err != nil {
		return err
	}

	if _, err := peerPortFile.Seek(0, 0); err != nil {
		return err
	}

	writer := bufio.NewWriter(peerPortFile)
	for _, line := range stringsArray {
		_, err := fmt.Fprintln(writer, line)
		if err != nil {
			return err
		}
	}

	if err := writer.Flush(); err != nil {
		return err
	}

	log.Printf("Port %s removed from file", p.port)
	return nil
}

// Peer methods

func (p *Peer) updatePortList() (err error) {
	p.PortList, err = readFile("PeerPorts.txt")
	if err != nil {
		log.Fatalf("Could not read from file: %s", err)
	}
	p.updateFile("PeerPorts.txt")

	return nil
}

func (p *Peer) connect(port string) (err error) {
	log.Printf("Peer %s is connecting to peer %s", p.port, port)
	conn, err := grpc.Dial(":"+port, grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("could not connect to peer on port %s: %v", port, err)
	}

	p.peerList[port] = MeService.NewMeServiceClient(conn)
	return nil
}

func (p *Peer) leave() {
	p.sendConnectionStatus(false)
	err := p.deleteOwnPortFromFile("PeerPorts.txt")
	if err != nil {
		return
	}
}

func (p *Peer) Start() error {

	err := p.StartServer()
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
	err := p.updatePortList()
	if err != nil {
		return err
	}

	if len(p.PortList) != 0 {
		for i := 0; i < len(p.PortList); i++ {
			if p.port == p.PortList[i] {
				continue
			}
			err := p.connect(p.PortList[i])
			if err != nil {
				return err
			}
		}
		p.sendConnectionStatus(true)
		return nil
	}
	log.Printf("Peer %s has no peers to connect to", p.port)
	return nil
}

func (p *Peer) StartServer() error {
	listener, err := net.Listen("tcp", ":"+p.port)
	if err != nil {
		log.Fatalf("Could not listen to port: %s: %v", p.port, err)
	}

	fmt.Println("Now listening on port: ", p.port)

	p.server = grpc.NewServer()
	MeService.RegisterMeServiceServer(p.server, p.meServer)

	go func() {
		if err := p.server.Serve(listener); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
	return nil
}

// MakeRequest This is the Client part of the peer, where it sends requests
func (p *Peer) MakeRequest() {
	p.state = 1
	//p.lamportClock = p.lamportClock + (int64(len(p.peerList)))
	p.lamportClock++
	p.allowedTimestamp = p.lamportClock
	// create an array that holds responses
	var responseArray []*MeService.Message = make([]*MeService.Message, len(p.peerList))

	var wg sync.WaitGroup

	log.Printf("Peer %s is requesting entry @ lamport time %d", p.port, p.lamportClock)
	for _, client := range p.peerList {
		wg.Add(1)
		go func(client MeService.MeServiceClient) {
			defer wg.Done()
			response, _ := client.RequestEntry(context.Background(), &MeService.Message{
				Timestamp: p.lamportClock,
				NodeId:    p.port,
			})
			responseArray = append(responseArray, response)
			p.pickMaxAndUpdateClock(response.Timestamp)
		}(client)

	}

	wg.Wait()
	p.enterCriticalSection()
	p.leaveCriticalSection()
}

// RequestEntry This is the server part of the peer, where it handles how to return the actual rpc method
func (s *MeServiceServer) RequestEntry(_ context.Context, entryRequest *MeService.Message) (*MeService.Message, error) {
	p := s.Peer
	if p.state == 2 || (p.state == 1 && p.allowedTimestamp < entryRequest.Timestamp) ||
		(p.state == 1 && p.allowedTimestamp == entryRequest.Timestamp && p.port < entryRequest.NodeId) {

		p.pickMaxAndUpdateClock(entryRequest.Timestamp)

		log.Printf("Peer %s received request from peer %s but has priority and queues request @ lamport time %d", p.port, entryRequest.NodeId, p.lamportClock)

		for p.state == 2 {

		}
		p.lamportClock++

		return p.returnMessage(), nil
	} else {
		p.pickMaxAndUpdateClock(entryRequest.Timestamp)
		p.lamportClock++
		log.Printf("Peer %s responds to request from peer %s @ lamport time %d", p.port, entryRequest.NodeId, p.lamportClock)
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
	f, err := os.OpenFile("Log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening log: %v", err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatalf("Error closing log: %v", err)
		}
	}(f)
	log.SetOutput(f)

	peer1 := NewPeer("5001")
	peer2 := NewPeer("5002")
	peer3 := NewPeer("5003")

	err = peer1.Start()
	if err != nil {
		_ = fmt.Errorf("error starting peer1: %v", err)
	}
	err = peer2.Start()
	if err != nil {
		_ = fmt.Errorf("error starting peer1: %v", err)
	}
	err = peer3.Start()
	if err != nil {
		_ = fmt.Errorf("error starting peer1: %v", err)
	}

	time.Sleep(2 * time.Second)

	go peer1.MakeRequest()
	time.Sleep(2 * time.Second)
	go peer2.MakeRequest()
	time.Sleep(2 * time.Second)
	go peer3.MakeRequest()

	time.Sleep(25 * time.Second)
	peerPortFile := "PeerPorts.txt"
	// truncate file
	err = os.Truncate(peerPortFile, 0)
}
