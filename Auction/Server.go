package main

import (
	"bufio"
	"context"
	"errors"
	_ "errors"
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
	isPrimaryServer bool
	lamportClock    int32
}

type ConsensusServer struct {
	Consensus.UnimplementedConsensusServer
	ConsensusServer Consensus.ConsensusServer
	ConsensusClient Consensus.ConsensusClient
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

func CreateServer() *Server {
	// TODO Maybe the ConsensesServer and AuctionServer needs to be initialised somehow
	s := &Server{
		Port:            "1",
		ConsensusServer: &ConsensusServer{},
		AuctionServer:   &AuctionServer{},
		isPrimaryServer: false,
		lamportClock:    0,
	}

	s.ConsensusServer.BackupList = make(map[string]Consensus.ConsensusClient)
	s.AuctionServer.Server = s
	s.ConsensusServer.Server = s

	err := s.ConsensusServer.updatePortList()
	if err != nil {
		return nil
	}

	if len(s.ConsensusServer.PortList) == 0 {
		s.Port = "5001"
		s.isPrimaryServer = true
	} else {
		portInt, _ := strconv.Atoi(s.ConsensusServer.PortList[len(s.ConsensusServer.PortList)-1])
		fmt.Println(portInt)
		s.Port = strconv.Itoa(portInt + 1)
	}
	s.ConsensusServer.updateFile()
	fmt.Printf("Port List: %v \n", s.ConsensusServer.PortList)
	// Find the next port that is available, so that the servers are hopefully sequentially numbered

	return s
}

func (s *Server) StartServer() error {
	//Listen on port
	lis, _ := net.Listen("tcp", ":"+s.Port)
	//Create server
	grpcServer := grpc.NewServer()
	//Register servers
	Consensus.RegisterConsensusServer(grpcServer, s.ConsensusServer)
	Auction.RegisterAuctionServer(grpcServer, s.AuctionServer)
	fmt.Println("Registered servers")
	//Start server
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	fmt.Println("Served listener")
	fmt.Println(len(s.ConsensusServer.BackupList))
	fmt.Println("Will now send connection")
	s.ConsensusServer.sendConnection()

	return nil
}

// ======================================= CONSENSUS PART OF THE SERVER ================================================

func (s *ConsensusServer) sendConnection() {

	var connectionAck = &Consensus.Ack{
		Port:   s.Server.Port,
		Status: "0",
	}
	fmt.Printf("I now want to send connections. My portList looks like this: %v \n", s.PortList)
	if len(s.PortList) == 0 {
		fmt.Println("There are no other servers to connect to")
	} else {
		fmt.Println("Will now send connection to ports in portList")
		for _, port := range s.PortList {
			if port == s.Server.Port {
				fmt.Println("Found own port in portList")
				continue
			}
			fmt.Printf("Will now send ConsensusConnect to port %s", port)
			err := s.ConsensusConnect(port)
			if err != nil {
				return
			}
			fmt.Printf("Server %s connected to server %s", s.Server.Port, port)
		}
		fmt.Printf(strconv.Itoa(len(s.BackupList)))
		fmt.Println("Will now send ConnectStatus to all in BackupList")

		for _, target := range s.BackupList {
			if target == nil {
				fmt.Println("Server " + s.Server.Port + " Found nil target in backupList")
			}
			ack, err := target.ConnectStatus(context.Background(), connectionAck)
			if err != nil {
				log.Printf("Could not connect to backup server: %s", target)

			}
			if ack.Status == "1" {
				log.Printf("Could not connect to backup server: %s", target)
			}
			fmt.Println("Got ConnectStatus response from port: " + ack.Port)
			log.Printf("Server %s connected to backup server: %s", s.Server.Port, ack.Port)
		}
	}

}

func (s *ConsensusServer) ConnectionStatus(_ context.Context, inComing *Consensus.Ack) (*Consensus.Ack, error) {

	// If the inComing port is not contained in the BackupList, this server will connect to the inComing port
	if _, contained := s.BackupList[inComing.Port]; !contained {
		time.Sleep(time.Second * 5)
		fmt.Println("Will now tro to call consensusConnect on port: " + inComing.Port)
		err := s.ConsensusConnect(inComing.Port)
		if err != nil {
			fmt.Println("Happened error when trying to call ConsensusConnect to port: " + inComing.Port)
			return &Consensus.Ack{
				Port:   s.Server.Port,
				Status: "1",
			}, nil
		}
		fmt.Printf("server %s conncted to server %s \n", s.Server.Port, inComing.Port)
	}

	log.Printf("Server %s connected to backup server: %s \n", s.Server.Port, inComing.Port)
	return &Consensus.Ack{
		Port:   s.Server.Port,
		Status: "0",
	}, nil

}

func (s *ConsensusServer) ConsensusConnect(port string) error {
	fmt.Printf("Server %s will now create a connection to %s \n", s.Server.Port, port)

	// Establish a connection to the server
	conn, err := grpc.Dial("localhost:"+port, grpc.WithInsecure())
	if err != nil {
		log.Printf("Failed to connect to server %s: %v", port, err)
		return err
	}
	defer conn.Close() // Close the connection when done

	// Create a ConsensusClient
	client := Consensus.NewConsensusClient(conn)
	if client == nil {
		log.Printf("Failed to create ConsensusClient for server %s", port)
		return errors.New("failed to create ConsensusClient")
	}

	// Add the client to the BackupList
	s.BackupList[port] = client

	log.Printf("Successfully added server %s to BackupList", port)
	return nil
}

func (s *ConsensusServer) updatePortList() (err error) {
	err = s.readFile("ServerPorts")
	if err != nil {
		return err
	}

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

// election lets all the other servers know that a primary is down, implicitly that makes the sender the primary server
func (s *ConsensusServer) election(deadServer *ConsensusServer) {

	s.Server.isPrimaryServer = true
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
	if !s.validateBid(bidRequest) {
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

func (s *AuctionServer) validateBid(bidRequest *Auction.BidRequest) bool {

	if bidRequest.Bid <= s.Auction.HighestBid {
		return false
	}

	return true

}

// =======================================Lamport clock=======================================================================

func (s *AuctionServer) UpdateAndIncrementClock(inComingClock int32) {
	s.Server.lamportClock = MaxL(inComingClock, s.Server.lamportClock) + 1
}

func (s *AuctionServer) startTimer() {

}

// ======================================= Helper Methods ===============================================================

func MaxL(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

// ======================================= FILE METHODS =======================================================================

func (s *ConsensusServer) readFile(fileName string) error {
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
		return err
	}

	if len(peerPortArray) == 0 {
		log.Printf("File: %s is empty", fileName)
	}

	s.PortList = peerPortArray
	return nil
}

func (s *ConsensusServer) updateFile() {
	peerPortFile, err := os.OpenFile("ServerPorts", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Panicf("Could not open file: %s", err)
	}
	defer func(peerPortFile *os.File) {
		err := peerPortFile.Close()
		if err != nil {
			_ = fmt.Errorf("error closing file: %v", err)
		}
	}(peerPortFile)

	fmt.Println("File opened successfully")

	// TODO CANNOT WRITE TO FILE FOR SOME REASON
	if _, err := peerPortFile.WriteString(s.Server.Port + "\n"); err != nil {
		log.Fatalf("Could not add Port: %s to file", s.Server.Port)
		fmt.Println("Could not add Port to file")
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
	server := CreateServer()
	err = server.StartServer()
	if err != nil {
		fmt.Println("Something went wrong in the startServer method")
	}
	time.Sleep(time.Second)
	//TODO Initate the ping process
	fmt.Println(server.Port)
	fmt.Println(server.isPrimaryServer)
	fmt.Println(server.ConsensusServer.PortList)
	for true {

	}

}
