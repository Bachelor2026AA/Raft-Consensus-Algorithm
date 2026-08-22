package leaderchecker

import (
	"fmt"
	"raft-consensus/election"
	"raft-consensus/node"
)

func LeaderCheck(servers []node.Node) {
	fmt.Println("Checking if anyone is leader")
	var check = 0
	var leader int
	for _, s := range servers {
		if s.ServerState == node.Leader {
			check += 1
			leader = s.ID
		}
	}
	if check == 0 {
		fmt.Println("No current leader")
		election.StartElection(servers)
		return
	}
	fmt.Printf("Leader %d is still alive", leader)
}
