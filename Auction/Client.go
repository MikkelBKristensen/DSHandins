package main

import (
	"errors"
	"fmt"
	Auction "github.com/MikkelBKristensen/DSHandins/Auction/Proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"os"
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
	clientDeadline := time.Now().Add(10 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), clientDeadline)
	resp, err := c.auctioneer.Bid(ctx, &bid)

	if err != nil || errors.Is(ctx.Err(), context.Canceled) {
		log.Printf("could not place bid, switching Server %v", err)
		c.switchServer()
		c.sendBid(amount)
		
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
	cancel()
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
	j := 0
	for i := 5001; i < 5010; i++ {

		// Start by connecting to port 5001 and call GetPrimaryServer to find the primary server from that server
		port := strconv.Itoa(i)

		//Continue to loop around until we find the primary server
		if i == 5009 {
			i = 5001
			j++
		}
		if j == 10 {
			log.Fatalf("Could not find primary server, shutting down client")
		}
		conn, err := grpc.Dial("localhost:"+port, grpc.WithInsecure())
		if err != nil {
			log.Printf("could not connect to server on port: %s", port)
			continue
		}
		target := Auction.NewAuctionClient(conn)

		clientDeadline := time.Now().Add(5 * time.Second)
		ctx, cancel := context.WithDeadline(context.Background(), clientDeadline)

		isPrimary, err := target.GetPrimaryServer(ctx, &Auction.Empty{})
		if err != nil {
			fmt.Println(err)
			continue
		}

		if !isPrimary.PrimaryStatus {
			continue
		}
		cancel()
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
func handleBid(c *Client) {
	fmt.Println("Type the amount you want to bid")
	var amount int32
	_, err := fmt.Scan(&amount)
	if err != nil {
		log.Printf("Error: %v", err)
	}
	c.sendBid(amount)
	fmt.Println("Enter new command")
}

func (c *Client) handleInput() {
	fmt.Println("Place a bid by typing: bid")
	fmt.Println("Request result by typing: result")
	fmt.Println("Exit by typing: exit")
	var command string

	for {
		_, err := fmt.Scan(&command)
		if err != nil {
			log.Printf("Error: %v", err)
		}
		if command == "bid" {
			handleBid(c)
		} else if command == "result" {
			c.requestResult()
		} else if command == "exit" {
			return
		} else {
			fmt.Println("Unknown command")
			c.handleInput()
		}
	}

}

// main method
func main() {

	// Set up the log
	f, err := os.OpenFile("logfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("[SERVER]: error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	//Create client
	client := CreateClient()
	err = client.connectToServer()
	if err != nil {
		log.Printf("Client: %V, Could not connect to server", client.Id)
	}

	client.handleInput()

}
