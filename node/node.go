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

func newNode(ID int) *Node {
	p := new(Node)
	p.ID = ID
	p.ServerState = Follower
	return p
}
