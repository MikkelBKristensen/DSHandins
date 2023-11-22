package main

import (
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
	//TODO Send result to client

	//TODO Send result to backup servers
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
			//TODO Handle fail - Is the server still alive?
			log.Printf("Could not sync with backup server: %v", s.BackupList[target])
		case "2":
			// Exception
			log.Printf("Exception when syncing with backup server: %v", s.BackupList[target])
		}
	}
	return nil
}

func (s *ConsensusServer) Sync(_ context.Context, clientBid *Consensus.ClientBid) (ack *Consensus.Ack, err error) {

	s.Server.AuctionServer.Auction.HighestBid = clientBid.Bid
	s.Server.AuctionServer.Auction.HighestBidder = clientBid.Id

	return &Consensus.Ack{
		Status: "0",
	}, nil
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
		time.Sleep(5 * time.Second) // Every 5 seconds the server should ping the other servers

		go func() {
			_, err := targetServer.Ping(context.Background(), ack)
			errCh <- err // Send the error to the channel
		}()

		// Wait for either the goroutine to complete or the timeout

		select {
		case err := <-errCh:
			// Handle the error if needed
			if err != nil {
				if targetServer.Server.isPrimaryServer {
					//TODO CAll for election

				} else {
					//TODO remove server from backuplist

					for i := range s.BackupList {
						if {
							
						}
					}

					log.Printf("Could not ping server: %v", targetServer.Server.Port)
				}
			}
		case <-time.After(10 * time.Second):
			// Timeout after 10 seconds
			fmt.Println("Timeout reached, the goroutine did not complete in time.")
		}
	}

}

// ======================================= AUCTION PART OF THE SERVER ==================================================

func (s *AuctionServer) Bid(ctx context.Context, bidRequest *Auction.BidRequest) (resp *Auction.BidResponse, err error) {
	// Ack : 0 = success, 1 = fail, 2 = exception
	// 0 : bid accepted and synced between servers
	// 1 : bid not accepted, either too low, could not sync (maybe), (auction is over?)
	// 2 : Some exception happened, maybe timeout?

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
