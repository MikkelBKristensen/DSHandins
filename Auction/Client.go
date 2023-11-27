package main

import (
	"errors"
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

	//src: https://grpc.io/blog/deadlines/
	clientDeadline := time.Now().Add(5 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), clientDeadline)
	resp, err := c.auctioneer.Result(ctx, &result)

	//Check if the context deadline has exceeded
	if err != nil || errors.Is(ctx.Err(), context.Canceled) {
		c.switchServer()
		c.requestResult()
		log.Fatalf("could request result: %v", err)
	}
	cancel()

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

		//Continue to loop around until we find the primary server
		if i == 5009 {
			i = 5001
		}
		conn, err := grpc.Dial("localhost:"+port, grpc.WithInsecure())
		if err != nil {
			log.Printf("could not connect to server on port: %s", port)
		}
		target := Auction.NewAuctionClient(conn)

		//isPrimary, _ := target.GetPrimaryServer(context.Background(), &Auction.Empty{})

		//if !isPrimary.PrimaryStatus {
		//	continue
		//}
		log.Printf("Client %d connected to server on port: %s", c.Id, port)

		// Set the auctioneer to the primary server
		c.auctioneer = target
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
	err := client.connectToServer()
	if err != nil {
		log.Printf("Client: %V, Could not connect to server", client.Id)
	}

	time.Sleep(5 * time.Second)
	client.sendBid(10)

}
