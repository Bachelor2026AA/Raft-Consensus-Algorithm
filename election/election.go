package election

import (
	"fmt"
	"math/rand"
	"raft-consensus/node"
)

func PrepareNodes(servers []node.Node) {
	var nobody = 0
	for i := 0; i <= len(servers)-1; i++ {
		servers[i].VotedFor = nobody
	}

}

func SelectCandidate(servers []node.Node) *node.Node {
	state := &servers[rand.Intn(4)]
	if state.ServerState == node.Follower {
		state.ServerState = node.Candidate
		state.VotedFor = state.ID
		state.Term++
	}
	var candidate = state.ID
	fmt.Printf("Suitble candidate has been found the candidate is node nr %d\n", candidate)

	return state
}

func StartElection(servers []node.Node) {
	var candidate = SelectCandidate(servers)
	var votes = CollectVotes(servers, candidate)
	var answers = MajorityVote(len(servers), votes)
	if answers {
		candidate.ServerState = node.Leader
	} else {
		fmt.Println("Election Failed")
	}
}

func RequestVote(candidate *node.Node, server *node.Node) int {
	if server.ServerState == node.Disconnected {
		fmt.Println("I cant communicate with this Node")
		return -1
	}
	var nobody int = 0
	var notvoted int = 0
	var notgranted int = 0
	var granted int = 1
	var alreadyvoted int = 2
	fmt.Printf("Hello i am candidate %d\n", candidate.ID)
	fmt.Printf("I am currently on term %d\n", candidate.Term)
	fmt.Println("I would like to request for your vote")

	// Now node being asked needs to compare its term with the term given
	if server.Term < candidate.Term {
		server.Term = candidate.Term
		fmt.Println("Since a new term has started im joining that term")
		server.VotedFor = nobody
	}

	if server.Term > candidate.Term {
		fmt.Println("Since the node is term behind it isnt op to date and this cant vote for it")
		return notgranted
	}

	if server.Term == candidate.Term {
		// have i voted?
		if server.VotedFor != notvoted {
			if server.VotedFor == candidate.ID {
				fmt.Println("i have already voted for you.")
				return alreadyvoted
			} else {
				fmt.Println("I have already voted for another candidate this term")
				return notgranted
			}
		}
	}
	server.VotedFor = candidate.ID
	return granted
}

func CollectVotes(servers []node.Node, candidate *node.Node) int {
	var granted = 1
	var count int = 1
	for index, votes := range servers {
		if candidate.ID == votes.ID {
			continue
		}
		var result = RequestVote(candidate, &servers[index])
		if result == granted {
			count++
		}
	}
	return count
}

func MajorityVote(numofservers int, votes int) bool {
	fmt.Println("Calculating if majority has been achieved")
	if votes >= numofservers/2+1 {
		fmt.Println("majority was achieved")
		return true
	}
	fmt.Println("majority was not achieved")
	return false
}
