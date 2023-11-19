package main

import (
	"context"
	Consensus "github.com/MikkelBKristensen/DSHandins/Auction/Consensus"
	Auction "github.com/MikkelBKristensen/DSHandins/Auction/Proto"
	_ "google.golang.org/grpc"
	"log"
	"net"
	"time"
)

type Server struct {
	ConsensusServer Consensus.ConsensusServer
	ConsensusClient Consensus.ConsensusClient
	AuctionClients  []Auction.AuctionClient
	BackupList      []Consensus.ConsensusServer
	lamportClock    int32
}

// AuctionServer This struct is used in the Client/Server interaction happening in the Auction service
type AuctionServer struct {
	AuctionServer Auction.AuctionServer
	AuctionClients
	Auction         *ActiveAuction
	Server          *Server
	Id              int32
	isPrimaryServer bool
}

// ActiveAuction This struct is used as a container to collect the relevant auction data
type ActiveAuction struct {
	Bidders       map[int32]bool
	HighestBid    int32
	HighestBidder int32
	time.Duration
	isActive bool
}

// =======================================CONSENSUS PART OF THE SERVER==========================================================
func CreateServer() {
	listener, err := net.Listen("tcp", ":"+p.port)
	if err != nil {
		log.Fatalf("Could not listen to port: %s: %v", p.port, err)
	}
}

// =======================================AUCTION PART OF THE SERVER===========================================================
func (s *AuctionServer) Bid(ctx context.Context, bidRequest *Auction.BidRequest) (resp *Auction.BidResponse, err error) {
	// Ack : 0 = success, 1 = fail, 2 = exception
	// 0 : bid accepted and synced between servers
	// 1 : bid not accepted, either too low, could not sync (maybe), (auction is over?)

	//Step 1: Sync clock, Update and increment
	s.UpdateAndIncrementClock(bidRequest.Timestamp) //Sync clock and Increment

	//Step 2: Is auction still active
	if !s.auction.isActive {

		resp = &Auction.BidResponse{
			status:    1,
			timestamp: s.server.lamportClock,
		}
		return resp, nil
	}

	//Step 3: Check if the bidder is registered
	//If there are no other bidders, start Auction timer
	if len(s.auction.Bidders) == 0 {
		//TODO Start timer here and continue

		// Initate auction
		s.auction.isActive = true

		//Register bidder
		s.auction.Bidders[bidRequest.Id] = true

	} else if s.auction.Bidders[bidRequest.Id] { // Will evaluate to false if person is not in the map
		//If not registered, then register bidder
		s.auction.Bidders[bidRequest.Id] = true
	}

	//Step 4: Validate bid
	if bidRequest.Bid < s.auction.HighestBid {
		resp = &Auction.BidResponse{
			status:    1,
			timestamp: s.server.lamportClock,
		}
		return resp, nil
	}

	//Step 5: Update the highest bid and bidder, sync, send success response
	s.auction.HighestBid = bidRequest.Bid
	s.auction.HighestBidder = bidRequest.Id

	//TODO Send Sync with replica servers

	resp = &Auction.BidResponse{
		status:    0,
		timestamp: s.server.lamportClock,
	}

	return resp, nil
}

func (s *AuctionServer) Result(ctx context.Context, resultRequest *Auction.ResultRequest) (resp *Auction.ResultResponse, err error) {
	//States for an auction: result || highest bid so far

	status := ""
	if !s.auction.isActive {
		status = "result"
		return resp, nil
	} else {
		status = "highest bid so far"
	}

	resp = &Auction.ResultResponse{
		id:        s.auction.HighestBidder,
		bid:       s.auction.HighestBid,
		status:    status,
		timestamp: s.server.lamportClock,
	}

	return resp, nil
}

// =======================================Lamport clock=========================================================
func MaxL(a int32, b int32) int32 {
	if a > b {
		return a
	}
	return b
}
func (s *AuctionServer) UpdateAndIncrementClock(inComingClock int32) {
	s.server.lamportClock = MaxL(inComingClock, s.server.lamportClock) + 1
}
