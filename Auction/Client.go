package main

import (
	Consensus "github.com/MikkelBKristensen/DSHandins/Handin5/Consensus"
	_ "google.golang.org/grpc"
)
struct Client{
	int64 lamportclock;
	var stream Consensus_SyncClient;
}

