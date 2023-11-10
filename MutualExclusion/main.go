package main

import (
	"bufio"
	MeService "github.com/MikkelBKristensen/DSHandins/MutualExclusion/MeService"
	_ "github.com/MikkelBKristensen/DSHandins/MutualExclusion/MeService"
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

func (s *Peer) sendJoin(join *MeService.JoinMessage) {
	s.peerList, _ = readFile()

	message := MeService.JoinMessage{
		Port: s.port,
	}

	for i := 0; i < len(s.peerList); i++ {
		s2 := s.peerList[i]
		//Send to all peers
	}
}

func (s *Peer) receiveJoin(join *MeService.JoinMessage) {
	append(s.peerList, join.GetPort())
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
	if _, err := peerPortFile.WriteString(s.port + "\n") {

	}

}

func main() {

}
