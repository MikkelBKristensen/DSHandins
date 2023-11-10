package main

import (
	"context"
	"fmt"
	MeService "github.com/MikkelBKristensen/DSHandins/MutualExclusion/MeService"
	_ "github.com/MikkelBKristensen/DSHandins/MutualExclusion/MeService"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"
	"time"
)

type mutualExclusionServer struct {
	mu                sync.Mutex
	lamportClock      int64
	inCriticalSection bool
	pendingRequests   map[string]int64
}

func (s *mutualExclusionServer) RequestEntry(ctx context.Context, req *MeService.Request) (MeService.Response, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Update Lamport clock
	s.lamportClock++
	localTimestamp := s.lamportClock

	// If the local node is not in the critical section or has a lower timestamp, grant permission
	if !s.inCriticalSection || req.Timestamp < localTimestamp {
		// Grant permission
		response := &MeService.Response{
			Permission: true,
			Timestamp:  localTimestamp,
			NodeId:     "local", // You can replace this with the actual node ID
		}
		return response, nil
	}

	// Otherwise, queue the request
	s.pendingRequests[req.NodeId] = req.Timestamp

	// Deny permission for now
	response := &MeService.Response{
		Permission: false,
		Timestamp:  localTimestamp,
		NodeId:     "local", // You can replace this with the actual node ID
	}
	return response, nil
}

func (s *mutualExclusionServer) ReleaseEntry(ctx context.Context, req *MeService.Release) (*MeService.Empty, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Update Lamport clock
	s.lamportClock++

	// Remove the released node from the pending requests queue
	delete(s.pendingRequests, req.NodeId)

	// Check if the pending requests can now be granted
	for node, timestamp := range s.pendingRequests {
		if timestamp <= s.lamportClock {
			// Grant permission to the next node in the queue
			response := &MeService.Response{
				Permission: true,
				Timestamp:  s.lamportClock,
				NodeId:     "local", // You can replace this with the actual node ID
			}
			go s.sendPermission(node, response)
			break
		}
	}

	return &MeService.Empty{}, nil
}

func (s *mutualExclusionServer) sendPermission(nodeID string, response *MeService.Response) {
	// Simulate network delay
	time.Sleep(100 * time.Millisecond)

	// Connect to the node and send the permission response
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", "localhost", 50000), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()

	client := MeService.NewMutualExclusionServiceClient(conn)
	client.RequestEntry(context.Background(), &MeService.Request{
		Timestamp: response.Timestamp,
		NodeId:    response.NodeId,
	})
}

func main() {
	listen, err := net.Listen("tcp", ":50000")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	server := grpc.NewServer()

	mutualExclusionService := &mutualExclusionServer{
		pendingRequests: make(map[string]int64),
	}

	pb.RegisterMutualExclusionServiceServer(server, mutualExclusionService)

	log.Println("Server started on :50000")
	if err := server.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
