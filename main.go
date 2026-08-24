package main

import (
	"raft-consensus/election"
	"raft-consensus/leaderchecker"
	"raft-consensus/node"
)

func main() {
	servers := []node.Node{
		{ID: 1, ServerState: node.Follower},
		{ID: 2, ServerState: node.Follower},
		{ID: 3, ServerState: node.Follower},
		{ID: 4, ServerState: node.Follower}}
	leaderchecker.LeaderCheck(servers)
	election.KillLeader(servers)
	leaderchecker.LeaderCheck(servers)

}
