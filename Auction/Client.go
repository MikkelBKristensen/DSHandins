package main

import (
	Auction "github.com/MikkelBKristensen/DSHandins/Auction/Proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
)

type Client struct {
	Id           int32
	lamportClock int32
	auctioneer   Auction.AuctionClient
	Servers      []string
}

func (c *Client) sendBid(amount int32) {
	c.lamportClock++

	bid := Auction.BidRequest{
		Id:        c.Id,
		Bid:       amount,
		Timestamp: c.lamportClock,
	}
	resp, err := c.auctioneer.Bid(context.Background(), &bid)
	if err != nil {
		log.Fatalf("could not place bid: %v", err)
	}
	//Sync clock according to response
	c.lamportClock = MaxL(c.lamportClock, resp.Timestamp)

	//@TODO Do something with resp here
	// Ack : 0 = success, 1 = fail, 2 = exception
	// 0 : bid accepted and synced between servers
	// 1 : bid not accepted, either too low, could not sync (maybe), (auction is over?)

	switch resp.Status {
	case "success":
		return //We don't need to do anything at this point
	case "fail":

	case "exception":

	default:
		log.Fatalf("Unknown status: %v", resp.Status)
	}

}
func (c *Client) requestResult() {
	c.lamportClock++

	result := Auction.ResultRequest{
		Id: c.Id,
		//Timestamp: c.lamportClock,
	}
	resp, err := c.auctioneer.Result(context.Background(), &result)
	if err != nil {
		log.Fatalf("could request result: %v", err)
	}
	//Sync clock according to response
	c.lamportClock = MaxL(c.lamportClock, resp.Timestamp)

	//@TODO Do something with resp her
	//States for an auction: result || highest bid

}

func (c *Client) joinAuction() error {

	// Connect to server and get client
	//TODO find a way to make client connect to primary server
	conn, err := grpc.Dial("5001", grpc.WithInsecure())
	if err != nil {
		return err
	}

	c.auctioneer = Auction.NewAuctionClient(conn)

	return nil
}

func MaxL(a int32, b int32) int32 {
	if a > b {
		return a
	}
	return b
}
