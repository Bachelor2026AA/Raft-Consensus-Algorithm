package raft

import (
	"context"
	"log"
	pb "raft-consensus/proto"
	"time"
)

func (n *Node) StartElection() {
	log.Printf("node %d has started election", n.ID)
	n.Term++
	n.Role = Candidate
	n.VotedFor = n.ID
	clear(n.VotesReceived)
	n.VotesReceived[n.ID] = struct{}{}

	lastTerm := 0
	if len(n.Logs) > 0 {
		lastTerm = n.Logs[len(n.Logs)-1].Term
	}
	voteReq := &pb.RequestVoteRequest{
		Term:        int32(n.Term),
		CandidateId: int32(n.ID),
		LogTerm:     int32(lastTerm),
		LogLength:   int32(len(n.Logs)),
	}
	for _, client := range n.peerClients {
		voteResp, err := client.RequestVote(context.Background(), voteReq)
		if err != nil {
			log.Printf("error: %v", err)
			return
		}
		if n.Role == Candidate && int(voteResp.Term) == n.Term && voteResp.VoteGranted {
			n.VotesReceived[int(voteResp.NodeId)] = struct{}{}
			majority := (len(n.peerClients)+1)/2 + 1
			if len(n.VotesReceived) >= majority {
				log.Printf("node %d is the new leader for term %d", n.ID, n.Term)
				n.Role = Leader
				n.Leader = n.ID
				n.timer.Stop()

				ticker := time.NewTicker(2 * time.Second)
				defer ticker.Stop()

				for {
					select {
					case <-ticker.C:
						for id, _ := range n.peerClients {
							log.Print("leader sent heartbeat to followers")
							n.SentLength[id] = len(n.Logs)
							n.AckedLength[id] = 0
							n.replicateLog(context.Background(), n.Leader, id)
						}
						//case keyVal := <-n.channel:
						//	for id, _ := range n.peerClients {
						//		n.replicateLog(context.Background(), n.Leader, id)
						//	}

					}
				}
			}
		} else if int(voteResp.Term) > n.Term {
			n.Term = int(voteResp.Term)
			n.Role = Follower
			n.VotedFor = -1
			n.resetElectionTimer()
		}
	}
}

func (n *Node) replicateLog(ctx context.Context, leaderId, followerId int) {
	prefixLen := n.SentLength[followerId]
	var suffix []*pb.Log
	for i := prefixLen; i < len(n.Logs); i++ {
		suffix = append(suffix, &pb.Log{
			Term: int32(n.Term),
			Data: n.Logs[i].Command,
		})
	}
	prefixTerm := 0
	if prefixLen > 0 {
		prefixTerm = n.Logs[prefixLen-1].Term
	}
	req := &pb.AppendEntriesRequest{
		LeaderId:     int32(leaderId),
		Term:         int32(n.Term),
		PrefixLen:    int32(prefixLen),
		PrefixTerm:   int32(prefixTerm),
		LeaderCommit: int32(n.CommitLength),
		Suffix:       suffix,
	}
	resp, err := n.peerClients[followerId].AppendEntries(ctx, req)
	if err != nil {
		log.Printf("error: %v", err)
		return
	}
	if int(resp.Term) > n.Term {
		n.Role = Follower
		n.Term = int(resp.Term)
		n.VotedFor = -1
		return
	}
	if resp.Success {

	}
}
