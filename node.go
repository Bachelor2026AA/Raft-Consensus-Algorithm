package raft

import (
	"math/rand/v2"
	pb "raft-consensus/proto"
	"raft-consensus/repository"
	"time"
)

type Role int

const (
	Follower Role = iota
	Candidate
	Leader
)

type Node struct {
	// persist if crashed

	ID           int
	Term         int
	VotedFor     int
	Logs         []Log
	CommitLength int

	// reset on crash

	Role          Role
	Leader        int
	VotesReceived map[int]struct{}
	SentLength    map[int]int
	AckedLength   map[int]int

	pb.UnimplementedRaftServer
	logEntryRepo repository.LogEntryRepository
	keyValueRepo repository.KeyValueRepository

	peerClients map[int]pb.RaftClient
	peerIds     []int

	timer *time.Timer
}

func NewNode(
	ID int,
	logEntryRepo repository.LogEntryRepository,
	keyValueRepo repository.KeyValueRepository,
	peerClients map[int]pb.RaftClient,
) *Node {
	n := &Node{
		ID:           ID,
		Term:         0,
		VotedFor:     -1,
		CommitLength: 0,

		Role:          Follower,
		Leader:        -1,
		VotesReceived: map[int]struct{}{},
		SentLength:    map[int]int{},
		AckedLength:   map[int]int{},

		logEntryRepo: logEntryRepo,
		keyValueRepo: keyValueRepo,

		peerClients: peerClients,
	}
	minDur := 8 * time.Second
	maxDur := 10 * time.Second
	n.timer = time.AfterFunc(minDur+rand.N(maxDur-minDur), n.StartElection)
	return n
}

func (n *Node) resetElectionTimer() {
	minDur := 8 * time.Second
	maxDur := 10 * time.Second
	n.timer.Stop()
	n.timer.Reset(minDur + rand.N(maxDur-minDur))
}
