package raft

import (
	"context"
	"log"
	customLog "raft-consensus/log"
	pb "raft-consensus/proto"
)

func (n *Node) RequestVote(ctx context.Context, req *pb.RequestVoteRequest) (*pb.RequestVoteResponse, error) {
	cTerm := int(req.GetTerm())
	cLogTerm := int(req.GetLogTerm())
	cId := int(req.GetCandidateId())
	cLogLength := int(req.GetLogLength())

	if cTerm > n.Term {
		n.Term = cTerm
		n.Role = Follower
		n.VotedFor = -1
	}

	voteGranted := vote(
		n.Term,
		n.VotedFor,
		n.Logs,
		cTerm,
		cLogTerm,
		cLogLength,
		cId,
	)
	if voteGranted {
		n.VotedFor = cId
		n.resetElectionTimer()
	}
	log.Printf("node %d was requested to vote for candidate %d, voted granted: %t", n.ID, cId, voteGranted)
	return &pb.RequestVoteResponse{
		Term:        int32(n.Term),
		VoteGranted: voteGranted,
		NodeId:      int32(n.ID),
	}, nil
}

func (n *Node) AppendEntries(ctx context.Context, req *pb.AppendEntriesRequest) (*pb.AppendEntriesResponse, error) {
	n.resetElectionTimer()
	if len(req.Suffix) > 0 && len(n.Logs) > int(req.PrefixLen) {
		index := min(len(n.Logs), int(req.PrefixLen)+len(req.Suffix)) - 1

		if n.Logs[index].Term != int(req.Suffix[index-int(req.PrefixLen)].Term) {
			n.Logs = n.Logs[:index]
		}
	}

	if len(req.Suffix)+int(req.PrefixLen) > len(n.Logs) {
		for index := 0; index < len(req.Suffix); index++ {
			n.Logs = append(n.Logs, customLog{
				Term: int(req.Suffix[index].Term),
			})
		}
	}

	if req.LeaderCommit > int32(n.CommitLength) {
		for index := n.CommitLength; index < int(req.LeaderCommit) && index < len(n.Logs); index++ {
			Deliver(n.Logs[index])
		}
		n.CommitLength = min(int(req.LeaderCommit), len(n.Logs))
	}
	if n.CommitLength == min(int(req.LeaderCommit), len(n.Logs)) {
		return &pb.AppendEntriesResponse{
			Term:    int32(n.Term),
			Ack:     int32(len(n.Logs)),
			Success: true,
		}, nil
	}
	return &pb.AppendEntriesResponse{
		Term:    int32(n.Term),
		Ack:     0,
		Success: false,
	}, nil

}

func vote(
	term, votedFor int, logs []customLog,
	cTerm, cLogTerm, cLogLength, cId int,
) bool {
	lastTerm := 0
	if len(logs) > 0 {
		lastTerm = logs[len(logs)-1].Term
	}
	logOk := cLogTerm > lastTerm || (cLogTerm == lastTerm && cLogLength >= len(logs))
	return cTerm == term && logOk && (votedFor == -1 || votedFor == cId)
}
