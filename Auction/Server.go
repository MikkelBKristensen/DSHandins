package main

import (
	"bufio"
	"context"
	"errors"
	_ "errors"
	"fmt"
	"log"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	Consensus "github.com/MikkelBKristensen/DSHandins/Auction/Consensus"
	Auction "github.com/MikkelBKristensen/DSHandins/Auction/Proto"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc"
)

type Server struct {
	Port            string
	ConsensusServer *ConsensusServer
	AuctionServer   *AuctionServer
	isPrimaryServer bool
	lamportClock    int32
	grpcServer      *grpc.Server
}

type ConsensusServer struct {
	Consensus.UnimplementedConsensusServer
	ConsensusServer Consensus.ConsensusServer
	ConsensusClient Consensus.ConsensusClient
	BackupList      map[string]Consensus.ConsensusClient
	Server          *Server
	PortList        []string
	PortOfPrimary   string
}

// AuctionServer This struct is used in the Client/Server interaction happening in the Auction service
type AuctionServer struct {
	Auction.UnimplementedAuctionServer
	AuctionServer Auction.AuctionServer
	Auction       *ActiveAuction
	Server        *Server
}

// ActiveAuction This struct is used as a container to collect the relevant auction data
type ActiveAuction struct {
	HighestBid    int32
	HighestBidder int32
	timeofStart   time.Time
	isActive      bool
	isStarted     bool
	Duration      float64
}

// ======================================= CENTRAL SERVER PART =========================================================

func CreateServer() *Server {
	// TODO Maybe the ConsensusServer and AuctionServer needs to be initialised somehow
	s := &Server{
		Port: "1",
		ConsensusServer: &ConsensusServer{
			BackupList: make(map[string]Consensus.ConsensusClient),
			PortList:   make([]string, 0),
		},
		AuctionServer: &AuctionServer{
			Auction: &ActiveAuction{
				HighestBid:    0,
				HighestBidder: 0,
				isActive:      false,
				isStarted:     false,
				Duration:      60,
			},
		},
		isPrimaryServer: false,
	}

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
		s.Port = strconv.Itoa(portInt + 1)
	}
	s.ConsensusServer.updateFile()
	fmt.Printf("My portList: %v \n", s.ConsensusServer.PortList)
	// Find the next port that is available, so that the servers are hopefully sequentially numbered

	return s
}

func (s *Server) StartServer() error {
	//Listen on port
	lis, _ := net.Listen("tcp", ":"+s.Port)
	//Create server
	s.grpcServer = grpc.NewServer()
	//Register servers
	Consensus.RegisterConsensusServer(s.grpcServer, s.ConsensusServer)
	Auction.RegisterAuctionServer(s.grpcServer, s.AuctionServer)
	//Start server
	go func() {
		if err := s.grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	fmt.Println("Now served listener")
	fmt.Println("The length of my backupList is: ", len(s.ConsensusServer.BackupList))
	s.ConsensusServer.sendConnection()

	//
	if len(s.ConsensusServer.PortList) != 0 {
		port, _ := s.ConsensusServer.FindSuccessor()
		target := s.ConsensusServer.BackupList[port]
		s.ConsensusServer.PingServer(target)
	}

	return nil
}

// ======================================= CONSENSUS PART OF THE SERVER ================================================

func (s *ConsensusServer) sendConnection() {

	var connectionAck = &Consensus.Ack{
		Port:   s.Server.Port,
		Status: "0",
	}
	if len(s.PortList) == 0 {
		fmt.Println("There are no other servers to connect to")
	} else {
		fmt.Println("Will now send connection to ports in portList")
		for _, port := range s.PortList {
			if port == s.Server.Port {
				fmt.Println("Found own port in portList")
				continue
			}
			fmt.Printf("Will now send ConsensusConnect to port %s \n", port)
			err := s.ConsensusConnect(port)
			if err != nil {
				return
			}
			fmt.Printf("Server %s connected to server %s \n", s.Server.Port, port)
		}
		fmt.Println("This is the length of my backupList now: ", strconv.Itoa(len(s.BackupList)))
		fmt.Println("Will now send ConnectStatus to all in BackupList")

		if len(s.BackupList) > 1 {
			for port, target := range s.BackupList {
				if target == nil {
					fmt.Println("Found nil target in backupList for port:", port)
					continue
				}
				fmt.Println("Will now call ConnectStatus on target")
				ack, err := target.ConnectStatus(context.Background(), connectionAck)
				if err != nil {
					log.Printf("Could not connect to backup server %s: %v\n", target, err)
				}
				fmt.Println("Have now called ConnectStatus on target")

				if ack.Status == "1" {
					log.Printf("Received ack status of 1 from server \n")
				}

				fmt.Println("Got ConnectStatus response from port: " + ack.Port)
				log.Printf("Server %s connected to backup server: %s \n", s.Server.Port, ack.Port)
			}
		} else {
			_, err := s.BackupList[s.PortList[0]].ConnectStatus(context.Background(), connectionAck)
			if err != nil {
				log.Printf("Primary %v Could not connect to backup server %s: %v\n", s.Server.Port, s.PortList[0], err)
			}
		}

	}

}

func (s *ConsensusServer) ConnectStatus(_ context.Context, inComing *Consensus.Ack) (*Consensus.Ack, error) {

	// If the inComing port is not contained in the BackupList, this server will connect to the inComing port
	if _, contained := s.BackupList[inComing.Port]; !contained {
		// Add inComing port to portlist
		s.PortList = append(s.PortList, inComing.Port)

		// Let other server initialise fully
		time.Sleep(time.Second * 1)

		fmt.Println("Will now try to call consensusConnect on port: " + inComing.Port)

		//Adds server to backuplist
		err := s.ConsensusConnect(inComing.Port)
		if err != nil {
			fmt.Println("Error happened when trying to call ConsensusConnect to port: " + inComing.Port)
			resp := &Consensus.Ack{
				Port:   s.Server.Port,
				Status: "1",
			}
			return resp, nil
		}
		fmt.Printf("server %s conncted to server %s \n", s.Server.Port, inComing.Port)

		// Set new server as successor
		if len(s.PortList) >= 1 {

			targetPort, _ := s.FindSuccessor()
			target := s.BackupList[targetPort]
			s.PingServer(target)
			log.Printf("Server %s has set server %s as successor", s.Server.Port, inComing.Port)
		}

	}

	log.Printf("Server %s connected to backup server: %s \n", s.Server.Port, inComing.Port)
	resp := &Consensus.Ack{
		Port:   s.Server.Port,
		Status: "0",
	}
	return resp, nil

}

func (s *ConsensusServer) ConsensusConnect(port string) error {
	fmt.Printf("Server %s will now create a connection to %s \n", s.Server.Port, port)

	// Establish a connection to the server
	conn, err := grpc.Dial("localhost:"+port, grpc.WithInsecure())
	if err != nil {
		log.Printf("Failed to connect to server %s: %v", port, err)
		return err
	}

	// Create a ConsensusClient
	client := Consensus.NewConsensusClient(conn)
	if client == nil {
		log.Printf("Failed to create ConsensusClient for server %s", port)
		return errors.New("failed to create ConsensusClient")
	}

	// Add the client to the BackupList
	s.BackupList[port] = client

	log.Printf("%v: Successfully added server %s to BackupList", s.Server.Port, port)
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
		Id:           bidReq.Id,
		Bid:          bidReq.Bid,
		AuctionStart: s.Server.AuctionServer.Auction.timeofStart.String(),
	}

	var wg = sync.WaitGroup{}

	//TODO Send result to client

	//TODO Send result to backup servers

	for port, target := range s.BackupList {
		wg.Add(1)
		ack, err := target.Sync(context.Background(), clientBid)
		if err != nil {
			log.Printf("Could not sync with backup server: %v", err)
			continue
		}

		// create switch case for ack.Status
		switch ack.Status {
		case "0":
			//Success
			log.Printf("Primary Server %v synced with backup server: %v", s.Server.Port, port)
			wg.Done()
		case "1":
			//Fail
			//TODO Handle fail - Is the server still alive?
			log.Printf("Primary Server %v Could not sync with backup server: %v", s.Server.Port, target)
			return err
		case "2":
			// Exception
			// TODO Was would happen in this case?
			log.Printf("Primary Server %v: Exception when syncing with backup server: %v", s.Server.Port, target)
			return err
		}
	}

	wg.Wait()
	log.Printf("All backup servers was synced successfully")
	return nil
}

func (s *ConsensusServer) Sync(_ context.Context, clientBid *Consensus.ClientBid) (ack *Consensus.Ack, err error) {

	if clientBid.Bid > s.Server.AuctionServer.Auction.HighestBid {
		//Take the timestamp from the bid, and set the auction time to that in case Primary server is down
		s.Server.AuctionServer.Auction.timeofStart, err = time.Parse(time.RFC3339, clientBid.AuctionStart)
		s.Server.AuctionServer.Auction.HighestBid = clientBid.Bid
		s.Server.AuctionServer.Auction.HighestBidder = clientBid.Id
		log.Printf("Server %s was successfully synced with primary server", s.Server.Port)
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
func (s *ConsensusServer) Ping(_ context.Context, ack *Consensus.Ack) (*Consensus.PingResponse, error) {
	//Respond to ping
	return &Consensus.PingResponse{
		Status:    "0",
		IsPrimary: s.Server.isPrimaryServer,
	}, nil
}

func (s *ConsensusServer) PingServer(target Consensus.ConsensusClient) {

	ack := &Consensus.Ack{
		Status: "4",
	}

	for {
		time.Sleep(1 * time.Second) // Every 1 seconds the server should ping the other servers

		failure := ""

		clientDeadline := time.Now().Add(5 * time.Second)
		ctx, cancel := context.WithDeadline(context.Background(), clientDeadline)

		_, err := target.Ping(ctx, ack)
		if err != nil || errors.Is(ctx.Err(), context.Canceled) {
			// Handle the error if needed
			failure = "failure"
		}
		cancel()

		//Get Port of client from backuplist
		targetPort := ""
		for port, server := range s.BackupList {
			if server == target {
				targetPort = port
			}
		}
		// Wait for either the goroutine to complete or the timeout

		switch {
		case failure == "failure":
			//remove server from backuplist
			// Handle the error if needed
			if err != nil {
				//Update list and slice
				delete(s.BackupList, targetPort)

				for i, port := range s.PortList {
					if port == targetPort {
						s.PortList = append(s.PortList[:i], s.PortList[i+1:]...)
					}
				}
				if targetPort == s.PortOfPrimary {
					log.Printf("Server: %v has detected that the primary server is down, initiating election", s.Server.Port)
					s.initiateElection()

					break

				} else if targetPort != s.PortOfPrimary {
					//We've pinged a replica, and detected that it is down
					break

				}
			}
		case failure == "":
			continue
		}
	}
	//If the code reaches this point, it's time to find a new successor
	port, err := s.FindSuccessor()
	if err == nil {
		client := s.BackupList[port]
		s.PingServer(client)
	}
	log.Printf("Could not find successor: %v", err)

}

// ======================================= ELECTION  ===================================================================

// election lets all the other servers know that a primary is down, implicitly that makes the sender the primary server
func (s *ConsensusServer) initiateElection() {

	//If no other nodes are in the ring, this node can proclaim itself as primary
	if len(s.BackupList) == 0 {
		log.Printf("Server: %v is the new primary server, no other backup server was found", s.Server.Port)
		s.Server.isPrimaryServer = true
		s.PortOfPrimary = s.Server.Port
		return
	}

	//Initialize election list with own port, and send out to next node in the ring
	election := &Consensus.Command{}
	election.Ports = append(election.Ports, s.Server.Port)

	//Send election to successor:
	successorPort, err := s.FindSuccessor()
	if err != nil {
		log.Printf("Node %v Could not find successor: %v", s.Server.Port, err)
	}

	_, err = s.BackupList[successorPort].Election(context.Background(), election)
	if err != nil {
		log.Printf("Server: %v could not send election command to server: %v (Election)", s.Server.Port, successorPort)
	}

}

func (s *ConsensusServer) Election(_ context.Context, incomming *Consensus.Command) (*Consensus.Empty, error) {
	//Check if own port is in the list, if so then the election needs to be completed
	for _, port := range incomming.Ports {
		log.Printf("Server: %v received election command from server: %v (Election)", s.Server.Port, port)
		if port == s.Server.Port {

			s.ElectAndCoordinate(incomming)
			return &Consensus.Empty{}, nil
		}
	}

	//Send command to successor
	command := &Consensus.Command{}
	//If own port is not in the list, append it and send to successor
	command.Ports = append(command.Ports, s.Server.Port)
	successorPort, err := s.FindSuccessor()
	if err != nil {
		log.Printf("Node %v Could not find successor: %v(Election)", s.Server.Port, err)
	}
	s.BackupList[successorPort].Election(context.Background(), command)

	return &Consensus.Empty{}, nil
}

func (s *ConsensusServer) ElectAndCoordinate(command *Consensus.Command) {
	//Sort so that the highest port is the last in the received list
	sort.Slice(command.Ports, func(i, j int) bool {
		return command.Ports[i] < command.Ports[j]
	})
	//Coordinate new leader with all nodes in the ring
	result := &Consensus.Coordinator{
		Port: command.Ports[len(command.Ports)-1],
	}
	log.Printf("Server: %v is the new primary server", result.Port)

	for _, node := range s.BackupList {
		_, err := node.Leader(context.Background(), result)
		if err != nil {
			log.Printf("Server: %v could not send leader command to server: %v (ElectAndCoordinate)", s.Server.Port, node)
		}
		log.Printf("Server: %v sent leader command to server: %v (ElectAndCoordinate)", s.Server.Port, node)
	}

}
func (s *ConsensusServer) Leader(_ context.Context, coordinator *Consensus.Coordinator) (*Consensus.Empty, error) {

	//TODO Consider dropping the hasHadElection bool

	if coordinator.Port == s.Server.Port {
		log.Printf("Server: %v is the new primary server", s.Server.Port)
		s.Server.isPrimaryServer = true
		s.PortOfPrimary = s.Server.Port
		return &Consensus.Empty{}, nil
	}
	//Update own state
	s.Server.isPrimaryServer = false
	s.PortOfPrimary = coordinator.Port
	return &Consensus.Empty{}, nil
}

func (s *ConsensusServer) FindSuccessor() (port string, err error) {

	if len(s.PortList) == 0 {
		//Choose own port if no other nodes are present
		return s.Server.Port, nil
	} else if len(s.PortList) == 1 {
		//No sort on one element
		return s.PortList[0], nil
	}

	//Sort the port list
	sort.Slice(s.PortList, func(i, j int) bool {
		return s.PortList[i] < s.PortList[j]
	})

	for _, port := range s.PortList {
		if port > s.Server.Port {
			return port, nil
		}
	}

	//Return the first port, as there are no larger ports available in the ring
	return s.PortList[0], nil
}

// ======================================= AUCTION PART OF THE SERVER ==================================================

func (s *AuctionServer) Bid(ctx context.Context, bidRequest *Auction.BidRequest) (resp *Auction.BidResponse, err error) {
	// Ack : 0 = success, 1 = fail, 2 = exception
	// 0 : bid accepted and synced between servers
	// 1 : bid not accepted, either too low, could not sync (maybe), (auction is over?)
	// 2 : Some exception happened, not primary server
	if !s.Server.isPrimaryServer {
		resp = &Auction.BidResponse{
			Status: "exception",
		}
		return resp, nil
	}

	//Bid received at time:
	bidTime := time.Now()

	//Step 2: Is auction complete
	if s.Auction.isStarted && !s.Auction.isActive {

		resp = &Auction.BidResponse{
			Status: "fail",
		}
		return resp, nil
		//Evaluate if the 60 seconds has passed since the auction was started
	} else if s.Auction.isStarted && s.Auction.Duration < bidTime.Sub(s.Auction.timeofStart).Seconds() {

		s.Auction.isActive = false
		resp = &Auction.BidResponse{
			Status: "fail",
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
			Status: "fail",
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
		Status: "success",
	}

	return resp, nil
}

func (s *AuctionServer) StartAuction() {

	//TODOSkrive til de andre servere for at prøve at synkronisere auction clocks

	//Initiate auction
	s.Auction.isStarted = true
	s.Auction.isActive = true
	s.Auction.timeofStart = time.Now()
	//TODO Start timer here and continue - Maybe it should be its own method call

}

func (s *AuctionServer) Result(ctx context.Context, resultRequest *Auction.ResultRequest) (resp *Auction.ResultResponse, err error) {
	//States for an auction: result || the highest bid so far
	resp = &Auction.ResultResponse{
		Id:  s.Auction.HighestBidder,
		Bid: s.Auction.HighestBid,
	}
	//Time of request
	requestTime := time.Now()

	if !s.Auction.isActive {
		resp.Status = "EndResult"
		return resp, nil
	} else if !s.Auction.isStarted && !s.Auction.isActive {
		resp.Status = "NotStarted"
		return resp, nil
	} else if s.Auction.isStarted && s.Auction.Duration < requestTime.Sub(s.Auction.timeofStart).Seconds() {
		s.Auction.isActive = false
		resp.Status = "EndResult"
	} else {
		resp.Status = "Result"
		return resp, nil
	}
	return
}

func (c *AuctionServer) GetPrimaryServer(_ context.Context, _ *Auction.Empty) (*Auction.ServerResponse, error) {
	isPrimary := c.Server.isPrimaryServer
	resp := &Auction.ServerResponse{PrimaryStatus: isPrimary, PrimaryServerPort: c.Server.ConsensusServer.PortOfPrimary}
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
	s.Server.lamportClock = MaxLA(inComingClock, s.Server.lamportClock) + 1
}

// ======================================= Helper Methods ===============================================================

func MaxLA(a, b int32) int32 {
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
	//TODO Initiate the ping process

	// Enter || Ctrl + C to exit
	exit := ""
	fmt.Scanln(&exit)
}
