package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	Consensus "github.com/MikkelBKristensen/DSHandins/Auction/Consensus"
	Auction "github.com/MikkelBKristensen/DSHandins/Auction/Proto"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Server struct {
	Port            string
	ConsensusServer *ConsensusServer
	AuctionServer   *AuctionServer
	lamportClock    int32
	isPrimaryServer bool
}

type ConsensusServer struct {
	Consensus.UnimplementedConsensusServer
	ConsensusServer *Consensus.ConsensusServer
	ConsensusClient *Consensus.ConsensusClient
	BackupList      map[string]Consensus.ConsensusClient
	HasHadElection  bool
	Server          *Server
	PortList        []string
}

// AuctionServer This struct is used in the Client/Server interaction happening in the Auction service
type AuctionServer struct {
	Auction.UnimplementedAuctionServer
	AuctionServer  Auction.AuctionServer
	AuctionClients []Auction.AuctionClient
	Auction        *ActiveAuction
	Server         *Server
}

// ActiveAuction This struct is used as a container to collect the relevant auction data
type ActiveAuction struct {
	HighestBid    int32
	HighestBidder int32
	time.Duration
	isActive  bool
	isStarted bool
}

// ======================================= CENTRAL SERVER PART =========================================================

func (s *Server) CreateServer(port string) error {
	//Create server
	grpcServer := grpc.NewServer()

	//Register servers
	Consensus.RegisterConsensusServer(grpcServer, s.ConsensusServer)
	Auction.RegisterAuctionServer(grpcServer, s.AuctionServer)

	//Listen on port
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return err
	}

	//Start server
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
		return err
	}
	s.ConsensusServer.BackupList = make(map[string]Consensus.ConsensusClient)
	return nil
}

func (s *Server) StartServer() {

	//TODO Find port to create server on
	for i := 1; i < 4; i++ {
		port := "500" + strconv.Itoa(i)
		//Try to create server on port
		err := s.CreateServer(port)

		//If failed, connect to that port
		if err != nil {
			connectionError := s.
			if connectionError != nil {
				log.Panicf("Could not connect to port: %v", port)
			}
		} else {
			//Succesful creation of port
			if i == 1 {
				//Potential cause of splitbrain
				log.Printf("Server on port: %v is primary server(StartServer method)", s.Port)
				s.isPrimaryServer = true
			}
			break
		}

	}
}

// ======================================= CONSENSUS PART OF THE SERVER ================================================

func (s *ConsensusServer) sendConnection() {

	var connectionAck = &Consensus.Ack{
		Port:   s.Server.Port,
		Status: "0",
	}

	if len(s.BackupList) != 0 {
		for port := range s.BackupList {
			target := s.BackupList[port]
			ack, err := target.ConnectStatus(context.Background(), connectionAck)
			if err != nil {
				log.Printf("Could not connect to backup server: %s", port)
			}
			if ack.Status == "1" {
				log.Printf("Could not connect to backup server: %s", port)
			}
			log.Printf("Server %s connected to backup server: %s", s.Server.Port, ack.Port)
		}
	}
	log.Printf("There are no backup servers to connect to")
}

func (s *ConsensusServer) ConnectionStatus(_ context.Context, inComing *Consensus.Ack) (*Consensus.Ack, error) {
	err := s.Server.ConsensusConnect(inComing.Port)
	if err != nil {
		return &Consensus.Ack{
			Status: "1",
		}, nil
	}
	log.Printf("Server %s connected to backup server: %s", s.Server.Port, inComing.Port)
	return &Consensus.Ack{
		Port:   s.Server.Port,
		Status: "0",
	}, nil
}

func (s *Server) ConsensusConnect(port string) error {
	conn, err := grpc.Dial(port, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	client := Consensus.NewConsensusClient(conn)
	s.ConsensusServer.BackupList[port] = client
	return nil
}

func (s *ConsensusServer) updatePortList() (err error) {
	s.PortList, err = readFile("ServerPorts")
	if err != nil {
		log.Fatalf("Could not read from file: %s", err)
	}
	s.updateFile("ServerPorts")

	return nil
}

func (s *ConsensusServer) sendSync(bidReq *Auction.BidRequest) error {
	clientBid := &Consensus.ClientBid{
		Id:        bidReq.Id,
		Bid:       bidReq.Bid,
		Timestamp: bidReq.Timestamp,
	}

	var wg = sync.WaitGroup{}

	//TODO Send result to client

	//TODO Send result to backup servers

	for _, target := range s.BackupList {
		wg.Add(1)
		ack, err := target.Sync(context.Background(), clientBid)
		if err != nil {
			log.Printf("Could not sync with backup server: %v", err)
			return err
		}

		// create switch case for ack.Status
		switch ack.Status {
		case "0":
			//Success
			log.Printf("Synced with backup server: %v", target)
			wg.Done()
		case "1":
			//Fail
			//TODO Handle fail - Is the server still alive?
			log.Printf("Could not sync with backup server: %v", target)
			return err
		case "2":
			// Exception
			// TODO Was would happen in this case?
			log.Printf("Exception when syncing with backup server: %v", target)
			return err
		}
	}

	wg.Wait()
	log.Printf("All backup servers was synced successfully")
	return nil
}

func (s *ConsensusServer) Sync(_ context.Context, clientBid *Consensus.ClientBid) (ack *Consensus.Ack, err error) {

	if clientBid.Bid > s.Server.AuctionServer.Auction.HighestBid {
		s.Server.AuctionServer.Auction.HighestBid = clientBid.Bid
		s.Server.AuctionServer.Auction.HighestBidder = clientBid.Id
		log.Printf("Server %s was successfulyl synced with primary server", s.Server.Port)
		return &Consensus.Ack{
			Status: "0",
		}, nil
	} else {
		log.Fatalf("The bid is not valid for backup server %s, and could not sync", s.Server.Port)
		return &Consensus.Ack{
			Status: "1",
		}, nil
	}

}

// ======================================= CONSENSUS: FAIL DETECTION ================================================
func (s *ConsensusServer) Ping(_ context.Context, ack *Consensus.Ack) (*Consensus.Ack, error) {
	//Respond to ping
	return &Consensus.Ack{
		Status: "0",
	}, nil
}

func (s *ConsensusServer) PingServer(targetServer *ConsensusServer) {
	ack := &Consensus.Ack{
		Status: "4",
	}
	errCh := make(chan error)

	for {
		time.Sleep(10 * time.Second) // Every 5 seconds the server should ping the other servers

		go func() {
			_, err := targetServer.Ping(context.Background(), ack)
			if err != nil {
				errCh <- err
			} // Send the error to the channel
		}()

		// Wait for either the goroutine to complete or the timeout

		select {
		case err := <-errCh:
			//remove server from backuplist
			// Handle the error if needed
			if err != nil {
				delete(s.BackupList, targetServer.Server.Port)
				if targetServer.Server.isPrimaryServer && !s.HasHadElection {

					//TODO CAll for election
					s.election(targetServer)
				} else if !targetServer.Server.isPrimaryServer {

					log.Printf("Could not ping server: %v", targetServer.Server.Port)
				}
			}
		case <-time.After(10 * time.Second):
			// Timeout after 10 seconds
			fmt.Println("Timeout reached, the goroutine did not complete in time.")
		}
	}

}

// ======================================= ELECTION  ===================================================================



//election lets all the other servers know that a primary is down, implicitly that makes the sender the primary server
func (s *ConsensusServer) election(deadServer *ConsensusServer) {

	s.Server.isPrimaryServer = true;
	election := &Consensus.Command{
		Port: deadServer.Server.Port,
	}

	//Send election to all servers
	for deadServerPort, deadServer := range s.BackupList {
		_, err := deadServer.ElectionCommand(context.Background(), election)
		if err != nil {
			log.Printf("Server: %v could not send election command to server: %v (Election)", s.Server.Port, deadServerPort)
		}
	}

}

// ======================================= AUCTION PART OF THE SERVER ==================================================

func (s *AuctionServer) Bid(ctx context.Context, bidRequest *Auction.BidRequest) (resp *Auction.BidResponse, err error) {
	// Ack : 0 = success, 1 = fail, 2 = exception
	// 0 : bid accepted and synced between servers
	// 1 : bid not accepted, either too low, could not sync (maybe), (auction is over?)
	// 2 : Some exception happened, not primary server
	if !s.Server.isPrimaryServer {
		resp = &Auction.BidResponse{
			Status:    "exception",
			Timestamp: s.Server.lamportClock,
		}
		return resp, nil
	}
	//Step 1: Sync clock, Update and increment
	s.UpdateAndIncrementClock(bidRequest.Timestamp) //Sync clock and Increment

	//Step 2: Is auction complete
	if s.Auction.isStarted && !s.Auction.isActive {

		resp = &Auction.BidResponse{
			Status:    "fail",
			Timestamp: s.Server.lamportClock,
		}
		return resp, nil
	}

	//Step 3: Check if auction has begun, if not then start it
	if !s.Auction.isStarted {
		s.StartAuction()
	}

	//Step 4: Validate bid
	if bidRequest.Bid <= s.Auction.HighestBid {
		resp = &Auction.BidResponse{
			Status:    "fail",
			Timestamp: s.Server.lamportClock,
		}
		return resp, nil
	}

	//Step 5: Update the highest bid and bidder, sync, send success response
	s.Auction.HighestBid = bidRequest.Bid
	s.Auction.HighestBidder = bidRequest.Id

	//TODO Send Sync with replica servers
	err = s.Server.ConsensusServer.sendSync(bidRequest)
	if err != nil {
		log.Fatalf("Could not sync with backup server: %v", err)
	}

	resp = &Auction.BidResponse{
		Status:    "success",
		Timestamp: s.Server.lamportClock,
	}

	return resp, nil
}

func (s *AuctionServer) StartAuction() {

	//Skrive til de andre servere for at prøve at synkronisere auction clocks

	//Initiate auction
	s.Auction.isStarted = true
	s.Auction.isActive = true
	//TODO Start timer here and continue - Maybe it should be its own method call
	s.startTimer()

}

func (s *AuctionServer) Result(ctx context.Context, resultRequest *Auction.ResultRequest) (resp *Auction.ResultResponse, err error) {
	//States for an auction: result || the highest bid so far
	if !s.Server.isPrimaryServer {
		resp = &Auction.ResultResponse{
			Status:    "exception",
			Timestamp: s.Server.lamportClock,
		}
		return resp, nil
	}

	status := ""
	if !s.Auction.isActive {
		status = "EndResult"
		return resp, nil
	} else if !s.Auction.isStarted && !s.Auction.isActive {
		status = "NotStarted"
	} else {
		status = "Result"
	}

	resp = &Auction.ResultResponse{
		Id:        s.Auction.HighestBidder,
		Bid:       s.Auction.HighestBid,
		Status:    status,
		Timestamp: s.Server.lamportClock,
	}

	return resp, nil
}

func (s *AuctionServer) validateBid(bidRequest *Auction.BidRequest) (bool, error) {
	//Step 2: Is auction still active
	if !s.Auction.isActive {
		err := errors.New("auction is not active")
		return false, err
	}

	//Step 3: Check if the bidder is registered
	//If there are no other bidders, start Auction timer
	if len(s.Auction.Bidders) == 0 {

		// Initiate auction
		s.Auction.isActive = true
		//TODO Start timer here and continue - Maybe it should be its own method call

		//Register bidder
		s.Auction.Bidders[bidRequest.Id] = true

	} else if !s.Auction.Bidders[bidRequest.Id] { // If the bidder is not in the map
		//If not registered, then register bidder
		s.Auction.Bidders[bidRequest.Id] = true
	}

	//Step 4: Validate bid
	if bidRequest.Bid < s.Auction.HighestBid {
		err := errors.New("bid is too low")
		return false, err
	}

	//Step 5: Update the highest bid and bidder, sync, send success response
	s.Auction.HighestBid = bidRequest.Bid
	s.Auction.HighestBidder = bidRequest.Id

	return true, nil
}

// =======================================Lamport clock=======================================================================

func (s *AuctionServer) UpdateAndIncrementClock(inComingClock int32) {
	s.Server.lamportClock = MaxL(inComingClock, s.Server.lamportClock) + 1
}

func (s *AuctionServer) startTimer() {

}

// ======================================= FILE METHODS =======================================================================

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

func (s *ConsensusServer) updateFile(fileName string) {
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

	if _, err := peerPortFile.WriteString(s.Server.Port + "\n"); err != nil {
		log.Fatalf("Could not add Port: %s to file", s.Server.Port)
	}
	log.Printf("Port %s added to file", s.Server.Port)
}

func (s *ConsensusServer) deletePortFromFile(fileName string) error {
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
		if !strings.Contains(line, s.Server.Port) {
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
	log.Printf("Port %s removed from file", s.Server.Port)
	return nil
}

func main() {
	// Set up the log
	f, err := os.OpenFile("logfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("[SERVER]: error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	// Create and Start server
	Server := new(Server)
	Server.StartServer()

	//TODO Initate the ping process

}
