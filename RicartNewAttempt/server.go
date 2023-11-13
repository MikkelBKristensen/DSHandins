// server.go
package main

import (
	"context"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"

	mutualexclusion "path/to/your/mutual_exclusion" // Update with your actual protobuf package path
)

type server struct {
	mu         sync.Mutex
	isLocked   bool
	highestTs  int64
	requesting map[string]int64
}

func (s *server) RequestEntry(ctx context.Context, req *mutualexclusion.Request) (*mutualexclusion.Response, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	nodeID := req.GetNodeId()
	timestamp := req.GetTimestamp()

	if s.isLocked {
		s.requesting[nodeID] = timestamp
		log.Printf("Received Request from Node %s at timestamp %d. Currently locked.", nodeID, timestamp)
		return &mutualexclusion.Response{Success: false}, nil
	}

	s.isLocked = true
	s.highestTs = timestamp

	log.Printf("Received Request from Node %s at timestamp %d. Granting access.", nodeID, timestamp)
	return &mutualexclusion.Response{Success: true}, nil
}

func (s *server) ReleaseEntry(ctx context.Context, req *mutualexclusion.Release) (*mutualexclusion.Response, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	nodeID := req.GetNodeId()
	timestamp := req.GetTimestamp()

	if !s.isLocked {
		log.Printf("Received Release from Node %s at timestamp %d. Not locked.", nodeID, timestamp)
		return &mutualexclusion.Response{Success: false}, nil
	}

	log.Printf("Received Release from Node %s at timestamp %d. Releasing lock.", nodeID, timestamp)
	s.isLocked = false

	// Process pending requests
	for id, ts := range s.requesting {
		if ts < s.highestTs {
			delete(s.requesting, id)
			continue
		}

		// Grant access to the pending request
		delete(s.requesting, id)
		s.isLocked = true
		s.highestTs = ts
		log.Printf("Granting access to pending request from Node %s at timestamp %d.", id, ts)
		break
	}

	return &mutualexclusion.Response{Success: true}, nil
}

func main() {
	port := ":50051" // Use your desired port

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}

	s := grpc.NewServer()
	mutualexclusion.RegisterMutualExclusionServer(s, &server{
		requesting: make(map[string]int64),
	})

	log.Printf("Server listening on port %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
