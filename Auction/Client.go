package main

import (
	Auction "github.com/MikkelBKristensen/DSHandins/Auction/Proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"strconv"
	"time"
)

type Client struct {
	Id         int32
	auctioneer Auction.AuctionClient
}

func CreateClient() *Client {
	id := createRandId()
	return &Client{
		Id: id,
	}
}

func (c *Client) sendBid(amount int32) {
	log.Printf("Client %d wants to send bid of: %d", c.Id, amount)
	bid := Auction.BidRequest{
		Id:  c.Id,
		Bid: amount,
	}
	resp, err := c.auctioneer.Bid(context.Background(), &bid)
	if err != nil {
		c.switchServer()
		c.sendBid(amount)
		log.Fatalf("could not place bid: %v", err)
	}

	// success : bid accepted and synced between servers
	// fail : bid not accepted, either too low, could not sync (maybe), (auction is over?)
	// exception : Some exception happened, maybe timeout?

	//TODO Add log
	switch resp.Status {
	case "success":
		log.Println("Bid was accepted")
		return //We don't need to do anything at this point

	case "fail":
		log.Println("Bid was not accepted")

	case "exception":
		log.Printf("Bid was not accepted, exception: %v", resp.Status)

		//Maybe incorporate a goroutine here, to find the lead server faster
		c.switchServer()
		c.sendBid(amount)

	default:
		c.switchServer()
		c.sendBid(amount)
		log.Fatalf("Unknown status: %v", resp.Status)
	}

}

func (c *Client) requestResult() {
	log.Printf("Client %d wants to request result", c.Id)
	result := Auction.ResultRequest{
		Id: c.Id,
	}
	resp, err := c.auctioneer.Result(context.Background(), &result)
	if err != nil {
		c.switchServer()
		c.requestResult()
		log.Fatalf("could request result: %v", err)
	}

	switch resp.Status {
	case "EndResult":
		log.Printf("Auction is over, highest bidder is: %d with: %d DKK", resp.Id, resp.Bid)
	case "NotStarted":
		log.Println("You must bid to start the auction, no bid has been placed yet.")
	case "Result":
		log.Printf("Auction is still running, highest bidder is: %d with: %d DKK", resp.Id, resp.Bid)
	}
	//States for an auction: result || the highest bid

}

func (c *Client) switchServer() {
	err := c.connectToServer()
	if err != nil {
		log.Fatalf("could not switch to server: %v", err)
	}
}
func (c *Client) connectToServer() error {
	// Find primary server and create connection
	for i := 5001; i < 5010; i++ {
		// Start by connecting to port 5001 and call GetPrimaryServer to find the primary server from that server
		port := strconv.Itoa(i)
		conn, err := grpc.Dial("localhost:"+port, grpc.WithInsecure())
		if err != nil {
			log.Printf("could not connect to server on port: %s", port)
		}
		primaryServerPort, _ := Auction.NewAuctionClient(conn).GetPrimaryServer(context.Background(), &Auction.Empty{})

		// Connect to primary server
		primaryConn, err := grpc.Dial("localhost:"+primaryServerPort.PrimaryPort, grpc.WithInsecure())
		if err != nil {
			log.Printf("could not connect to server on port: %s", primaryConn)
		}

		// Set the auctioneer to the primary server
		c.auctioneer = Auction.NewAuctionClient(primaryConn)
		return nil
	}
	return nil
}

func MaxL(a int32, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

func createRandId() int32 {
	minId := 1
	maxId := 10000
	rand.Seed(time.Now().UnixNano())
	id := rand.Int31n(int32(rand.Intn(maxId-minId) + minId))
	return id
}

// main method
func main() {
	client := CreateClient()

}
