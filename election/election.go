package election

import (
	"fmt"
	"math/rand"
	"raft-consensus/node"
	"time"
)

var currentTerm int

func SelectCandidate(servers []node.Node) *node.Node {
	state := &servers[rand.Intn(4)]
	if state.ServerState == node.Follower {
		state.ServerState = node.Candidate
	}
	var candidate = state.ID
	fmt.Printf("Suitble candidate has been found the candidate is node nr %d\n", candidate)

	return state
}

func Voting(servers []node.Node, candidate *node.Node) int {
	var vote int
	var chosen = candidate
	vote = 1 // as the candidate votes for himself
	for _, votes := range servers {
		var check = 0
		if votes.ID == chosen.ID {
			continue
		} else {
			fmt.Printf("Did you want to vote for this candidate? nr %d\n", votes.ID)
			check = rand.Intn(2)
			if check == 1 {
				fmt.Println("yes")
				vote += 1
			} else {
				fmt.Println("No")
			}
		}
	}
	return vote
}

func StartElection(servers []node.Node) {

	var candidate = SelectCandidate(servers)
	currentTerm = currentTerm + 1
	//candidate.Term = candidate.Term + 1

	fmt.Printf("This is currently Term %d\n", currentTerm)

	fmt.Println("Voting for majority may begin")

	var answer = Voting(servers, candidate)

	if answer >= len(servers)/2+1 {
		fmt.Println("Leader found")
		candidate.ServerState = node.Leader

	} else {
		fmt.Println("No leader found because of lack of votes")
		fmt.Println("Starting new vote")
	}

	time.Sleep(time.Duration(rand.Intn(10)) * time.Second)

}
