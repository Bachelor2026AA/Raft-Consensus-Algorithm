package raft

type Role int

const (
	Follower Role = iota
	Candidate
	Leader
)

type Node struct {
	// persist if crashed

	ID           int
	CurrentTerm  int
	VotedFor     int
	Logs         []Log
	CommitLength int

	// reset on crash

	CurrentRole   Role
	CurrentLeader int
	VotesReceived map[int]struct{}
	SentLength    map[int]int
	AckedLength   map[int]int
}

func New(ID int) *Node {
	return &Node{
		ID:           ID,
		CurrentTerm:  0,
		VotedFor:     -1,
		CommitLength: 0,

		CurrentRole:   Follower,
		CurrentLeader: -1,
		VotesReceived: map[int]struct{}{},
		SentLength:    map[int]int{},
		AckedLength:   map[int]int{},
	}
}
