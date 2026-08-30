package raft

import (
	"context"
	"fmt"
	pb "raft-consensus/proto"
	"time"
)

func (n *Node) StartElection() {
	fmt.Printf("node %d has started election", n.ID)
	n.Term++
	n.Role = Candidate
	n.VotedFor = n.ID
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
			// TODO log error
			return
		}
		if n.Role == Candidate && int(voteResp.Term) == n.Term && voteResp.VoteGranted {
			n.VotesReceived[n.ID] = struct{}{}
			if len(n.VotesReceived) >= (len(n.peerClients) + 1/2) {
				n.Role = Leader
				n.Leader = n.ID
				n.timer.Stop()
				for id, _ := range n.peerClients {
					n.SentLength[id] = len(n.Logs)
					n.AckedLength[id] = 0
					n.replicateLog(context.Background(), n.Leader, id)
				}
			}
		} else if int(voteResp.Term) > n.Term {
			n.Term = int(voteResp.Term)
			n.Role = Follower
			n.VotedFor = -1
			n.timer.Reset(2 * time.Second)
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
