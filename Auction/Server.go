package main

import (
	"context"
	"errors"
	Consensus "github.com/MikkelBKristensen/DSHandins/Auction/Consensus"
	Auction "github.com/MikkelBKristensen/DSHandins/Auction/Proto"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc"
	"log"
	"net"
	"os"
	"slices"
	"strconv"
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
	BackupList      []Consensus.ConsensusClient
	Server          *Server
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
	Bidders       []int32
	HighestBid    int32
	HighestBidder int32
	time.Duration
	isActive bool
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
	return nil
}

func (s *Server) StartServer() {

	//TODO Find port to create server on
	for i := 1; i < 3; i++ {
		port := "500" + strconv.Itoa(i)
		//Try to create server on port
		err := s.CreateServer(port)

		//If fail, connect to that port
		if err != nil {
			connectionError := s.ConsensusConnect(port)
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

func (s *Server) ConsensusConnect(port string) error {
	conn, err := grpc.Dial(port, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	client := Consensus.NewConsensusClient(conn)
	s.ConsensusServer.BackupList = append(s.ConsensusServer.BackupList, client)
	return nil
}

func (s *ConsensusServer) sendSync(bidReq *Auction.BidRequest) error {
	clientBid := &Consensus.ClientBid{
		Id:        bidReq.Id,
		Bid:       bidReq.Bid,
		Timestamp: bidReq.Timestamp,
	}
	for target := range s.BackupList {
		ack, err := s.BackupList[target].Sync(context.Background(), clientBid)
		if err != nil {
			log.Printf("Could not sync with backup server: %v", err)
			return err
		}
		// create switch case for ack.Status
		switch ack.Status {
		case "0":
			//Success
			log.Printf("Synced with backup server: %v", s.BackupList[target])
		case "1":
			//Fail
			log.Printf("Could not sync with backup server: %v", s.BackupList[target])
		case "2":
			// Exception
			log.Printf("Exception when syncing with backup server: %v", s.BackupList[target])
		}
	}
	return nil
}

func (s *ConsensusServer) Sync(_ context.Context, clientBid *Consensus.ClientBid) (ack *Consensus.Ack, err error) {

	//TODO Check if the bidder is registered

	//TODO Update the highest bid and bidder, sync, send success response
	//TODO Send Sync with replica servers
	//TODO Send ack
	//TODO Send result to client
	//TODO Send result to backup servers
	return &Consensus.Ack{
		Status: "0",
	}, nil
}

// ======================================= AUCTION PART OF THE SERVER ==================================================

func (s *AuctionServer) Bid(ctx context.Context, bidRequest *Auction.BidRequest) (resp *Auction.BidResponse, err error) {
	// Ack : 0 = success, 1 = fail, 2 = exception
	// 0 : bid accepted and synced between servers
	// 1 : bid not accepted, either too low, could not sync (maybe), (auction is over?)
	// 2 : Some exception happened, maybe timeout?

	//Step 1: Sync clock, Update and increment
	s.UpdateAndIncrementClock(bidRequest.Timestamp) //Sync clock and Increment

	//Step 2: Is auction still active
	if !s.Auction.isActive {

		resp = &Auction.BidResponse{
			Status:    "1",
			Timestamp: s.Server.lamportClock,
		}
		return resp, nil
	}

	//Step 3: Check if the bidder is registered
	//If there are no other bidders, start Auction timer
	if len(s.Auction.Bidders) == 0 {

		// Initiate auction
		s.Auction.isActive = true
		//TODO Start timer here and continue - Maybe it should be its own method call

		//Register bidder
		s.Auction.Bidders = append(s.Auction.Bidders, s.Auction.Bidders[bidRequest.Id])

	} else if !slices.Contains(s.Auction.Bidders, bidRequest.Id) {
		// If the bidder is not contained in the slice, add the bidder
		s.Auction.Bidders = append(s.Auction.Bidders, s.Auction.Bidders[bidRequest.Id])
	}

	//Step 4: Validate bid
	if bidRequest.Bid < s.Auction.HighestBid {
		resp = &Auction.BidResponse{
			Status:    "1",
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
		Status:    "0",
		Timestamp: s.Server.lamportClock,
	}

	return resp, nil
}

func (s *AuctionServer) Result(ctx context.Context, resultRequest *Auction.ResultRequest) (resp *Auction.ResultResponse, err error) {
	//States for an auction: result || the highest bid so far

	status := ""
	if !s.Auction.isActive {
		status = "result"
		return resp, nil
	} else {
		status = "highest bid so far"
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

}
