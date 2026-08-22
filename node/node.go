package node

type ServerState int

const (
	Follower ServerState = iota
	Candidate
	Leader
)

type Node struct {
	ID          int
	ServerState ServerState
	Term        int
}
