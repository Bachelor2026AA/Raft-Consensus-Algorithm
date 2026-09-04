package raft

import (
	"context"
	"log"
	raftLog "raft-consensus/log"
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
	term := int(req.GetTerm())
	prefixLen := int(req.GetTerm())
	prefixTerm := int(req.GetPrefixTerm())
	suffix := req.GetSuffix()
	leaderCommit := int(req.GetLeaderCommit())

	if term > n.Term {
		n.Term = int(req.Term)
		n.VotedFor = -1
	}
	if term == n.Term {
		n.Role = Follower
		n.Leader = int(req.LeaderId)
		n.resetElectionTimer()
	}
	logOk := (len(n.Logs) >= prefixLen) &&
		(prefixLen == 0 || n.Logs[prefixLen-1].Term == prefixTerm)
	if term == n.Term && logOk {
		n.appendEntriesInner(prefixLen, leaderCommit, suffix)
		ack := prefixLen + len(suffix)
		return &pb.AppendEntriesResponse{
			Term:    int32(n.Term),
			Ack:     int32(ack),
			Success: true,
		}, nil
	}
	return &pb.AppendEntriesResponse{
		Term:    int32(n.Term),
		Ack:     0,
		Success: false,
	}, nil
}

func (n *Node) appendEntriesInner(prefixLen, leaderCommit int, suffix []*pb.Log) {
	if len(suffix) > 0 && len(n.Logs) > int(prefixLen) {
		index := min(len(n.Logs), int(prefixLen)+len(suffix)) - 1

		if n.Logs[index].Term != int(suffix[index-int(prefixLen)].Term) {
			n.Logs = n.Logs[:index]
		}
	}

	if len(suffix)+int(prefixLen) > len(n.Logs) {
		for i := 0; i < len(suffix); i++ {
			n.Logs = append(n.Logs, raftLog.Log{
				Term: int(suffix[i].Term),
			})
		}
	}

	if leaderCommit > n.CommitLength {
		for i := n.CommitLength; i < leaderCommit && i < len(n.Logs); i++ {
			Deliver()
		}
		n.CommitLength = leaderCommit
	}
}

func vote(
	term, votedFor int, logs []raftLog.Log,
	cTerm, cLogTerm, cLogLength, cId int,
) bool {
	lastTerm := 0
	if len(logs) > 0 {
		lastTerm = logs[len(logs)-1].Term
	}
	logOk := cLogTerm > lastTerm || (cLogTerm == lastTerm && cLogLength >= len(logs))
	return cTerm == term && logOk && (votedFor == -1 || votedFor == cId)
}
