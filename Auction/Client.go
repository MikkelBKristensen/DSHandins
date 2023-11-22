package main

import (
	"fmt"
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

	// succes : bid accepted and synced between servers
	// fail : bid not accepted, either too low, could not sync (maybe), (auction is over?)
	// exception : Some exception happened, maybe timeout?

	//TODO Add log
	switch resp.Status {
	case "success":
		fmt.Println("Bid was accepted")
		return //We don't need to do anything at this point

	case "fail":
		fmt.Println("Bid was not accepted")

	case "exception":
		fmt.Println("Bid was not accepted, exception: %v", resp.Status)

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
	switch resp.Status {
	case "EndResult":
		fmt.Printf("Auction is over, highest bidder is: %d with: %d DKK", resp.Id, resp.Bid)

	case "NotStarted":
		fmt.Println("You must bid to start the auction, no bid has been placed yet.")
	case "Result":
		fmt.Printf("Auction is still running, highest bidder is: %d with: %d DKK", resp.Id, resp.Bid)
	}
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
