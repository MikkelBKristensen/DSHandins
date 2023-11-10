package main

import (
	"bufio"
	MeService "github.com/MikkelBKristensen/DSHandins/MutualExclusion/MeService"
	_ "github.com/MikkelBKristensen/DSHandins/MutualExclusion/MeService"
	"google.golang.org/grpc"
	"log"
	"os"
)

type Peer struct {
	lamportClock      int32
	port              string
	inCriticalSection bool
	pendingRequests   []MeService.Request
	peerList          []string
}

func (s *Peer) sendConnectionStatus(isJoin bool) {
	//Status: 0 for leaving and 1 for joining
	// Update peerList to be equal the PeerPorts file
	// and append own port to the file
	s.peerList, _ = readFile()
	s.updateFile()

	message := MeService.ConnectionMsg{
		Port:   s.port,
		isJoin: isJoin,
	}

	for i := 0; i < len(s.peerList); i++ {
		s2 := s.peerList[i]
		//Send to all peers
	}
}

func (s *Peer) receiveJoin(join *MeService.ConnectionMsg) {
	// TODO Update clock
	if join.IsJoin == true {
		append(s.peerList, join.GetPort())
	} else {
		//Make
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
	// Set filename
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

func (s *Peer) connect(port string) grpc.ClientConn {
	panic("Not Implemented yet")
}

func (s *Peer) sendMessage(msg MeService.ConnectionMsg) {
	panic("Not implemented yet")
}

func main() {

}
