package node

type ServerState int

const (
	Follower ServerState = iota
	Candidate
	Leader
	Disconnected
)

type Node struct {
	ID          int
	ServerState ServerState
	Term        int
	VotedFor    int
}

func New(ID int) *Node {
	return &Node{
	ID = ID,
	ServerState = Follower}
}
