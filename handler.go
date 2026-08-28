package raft

import (
	"context"
	pb "raft-consensus/proto"
	"raft-consensus/repository"
)

type RaftServer struct {
	pb.UnimplementedRaftServer
	repo   repository.LogEntryRepository
	kvrepo repository.KeyValueRepository
	node   Node
}

func (r *RaftServer) RequestVote(ctx context.Context, req *pb.RequestVoteRequest) (*pb.RequestVoteResponse, error) {
	cTerm := int(req.GetTerm())
	cLogTerm := int(req.GetLogTerm())
	cId := int(req.GetCandidateId())
	cLogLength := int(req.GetLogLength())

	if cTerm > r.node.CurrentTerm {
		r.node.CurrentTerm = cTerm
		r.node.CurrentRole = Follower
		r.node.VotedFor = -1
	}
	lastTerm := 0
	if len(r.node.Logs) > 0 {
		lastTerm = r.node.Logs[len(r.node.Logs)-1].Term
	}
	logOk := cLogTerm > lastTerm ||
		(cLogTerm == lastTerm && cLogLength >= len(r.node.Logs))
	voteGranted := false
	if cTerm == r.node.CurrentTerm && logOk &&
		(r.node.VotedFor == -1 || r.node.VotedFor == cId) {
		r.node.VotedFor = cId
		voteGranted = true
	}
	return &pb.RequestVoteResponse{
		Term:        int32(r.node.CurrentTerm),
		VoteGranted: voteGranted,
		NodeId:      int32(r.node.ID),
	}, nil
}

func (r *RaftServer) AppendEntries(ctx context.Context, req *pb.AppendEntriesRequest) (*pb.AppendEntriesResponse, error) {
	if len(req.Suffix) > 0 && len(log) > int(req.PrefixLen) {
		index := min(len(log), int(req.PrefixLen)+len(req.Suffix)) - 1

		if log[index] != req.Suffix[index-int(req.PrefixLen)] {
			log = log[:index]
		}
	}

	if len(req.Suffix)+int(req.PrefixLen) > len(log) {
		for index := 0; index < len(req.Suffix); index++ {
			log = append(log, req.Suffix[index])
		}
	}
	// resolve conflict

	if req.LeaderCommit > int32(commitLength) {
		for index := commitLength; index < int(req.LeaderCommit) && index < len(log); index++ {
			deliver(log[index])
		}

		commitLength = min(int(req.LeaderCommit), len(log))
	}

	return &pb.AppendEntriesResponse{
		Term:    currentTerm,
		Ack:     int32(len(log)),
		Success: true,
	}, nil
}
