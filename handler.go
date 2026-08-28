package raft

import (
	"context"
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
	}
	return &pb.RequestVoteResponse{
		Term:        int32(n.Term),
		VoteGranted: voteGranted,
		NodeId:      int32(n.ID),
	}, nil
}

func (n *Node) AppendEntries(ctx context.Context, req *pb.AppendEntriesRequest) (*pb.AppendEntriesResponse, error) {
	//if len(req.Suffix) > 0 && len(log) > int(req.PrefixLen) {
	//	index := min(len(log), int(req.PrefixLen)+len(req.Suffix)) - 1
	//
	//	if log[index] != req.Suffix[index-int(req.PrefixLen)] {
	//		log = log[:index]
	//	}
	//}
	//
	//if len(req.Suffix)+int(req.PrefixLen) > len(log) {
	//	for index := 0; index < len(req.Suffix); index++ {
	//		log = append(log, req.Suffix[index])
	//	}
	//}
	//// resolve conflict
	//
	//if req.LeaderCommit > int32(commitLength) {
	//	for index := commitLength; index < int(req.LeaderCommit) && index < len(log); index++ {
	//		deliver(log[index])
	//	}
	//
	//	commitLength = min(int(req.LeaderCommit), len(log))
	//}
	//
	//return &pb.AppendEntriesResponse{
	//	Term:    currentTerm,
	//	Ack:     int32(len(log)),
	//	Success: true,
	//}, nil
	return nil, nil
}
