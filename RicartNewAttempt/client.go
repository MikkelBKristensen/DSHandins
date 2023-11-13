// client.go
package main

import (
	"context"
	"log"
	"math/rand"
	"sync"
	"time"

	"google.golang.org/grpc"

	mutualexclusion "github.com/MikkelBKristensen/DSHandins/NewAttempt/mutual_exclusion" // Update with your actual protobuf package path
)

func main() {
	rand.Seed(time.Now().UnixNano())

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		runPeer("Node1")
	}()

	go func() {
		defer wg.Done()
		runPeer("Node2")
	}()

	go func() {
		defer wg.Done()
		runPeer("Node3")
	}()

	wg.Wait()
}

func runPeer(nodeID string) {
	port := ":50051" // Use your desired port
	conn, err := grpc.Dial("localhost"+port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()

	client := mutualexclusion.NewMutualExclusionClient(conn)

	for i := 0; i < 5; i++ {
		time.Sleep(time.Duration(rand.Intn(3)) * time.Second)

		req := &mutualexclusion.Request{
			Timestamp: time.Now().UnixNano(),
			NodeId:    nodeID,
		}

		resp, err := client.RequestEntry(context.Background(), req)
		if err != nil {
			log.Printf("Error requesting entry: %v", err)
			continue
		}

		if resp.GetSuccess() {
			log.Printf("Node %s successfully entered the critical section.", nodeID)

			time.Sleep(time.Duration(rand.Intn(3)) * time.Second)

			releaseReq := &mutualexclusion.Release{
				Timestamp: time.Now().UnixNano(),
				NodeId:    nodeID,
			}

			releaseResp, err := client.ReleaseEntry(context.Background(), releaseReq)
			if err != nil {
				log.Printf("Error releasing entry: %v", err)
				continue
			}

			if releaseResp.GetSuccess() {
				log.Printf("Node %s successfully released the critical section.", nodeID)
			}
		} else {
			log.Printf("Node %s could not enter the critical section. Retrying...", nodeID)
		}
	}
}
