package raft

import (
	pb "raft-consensus/proto"
	"raft-consensus/repository"
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
	repo   repository.LogEntryRepository
	kvrepo repository.KeyValueRepository
}

func NewNode(ID int) *Node {
	return &Node{
		ID:           ID,
		Term:         0,
		VotedFor:     -1,
		CommitLength: 0,

		Role:          Follower,
		Leader:        -1,
		VotesReceived: map[int]struct{}{},
		SentLength:    map[int]int{},
		AckedLength:   map[int]int{},
	}
}
