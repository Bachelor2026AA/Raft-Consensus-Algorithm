package handler

import (
	"context"
	"raft-consensus/node"
	pb "raft-consensus/proto"
	"raft-consensus/repository"
)

type RaftServer struct {
	pb.UnimplementedRaftServer
	repo   repository.LogEntryRepository
	kvrepo repository.KeyValueRepository
	node   node.Node
}

func (r *RaftServer) RequestVote(ctx context.Context, req *pb.RequestVoteRequest) (*pb.RequestVoteResponse, error) {
	cTerm := int(req.GetTerm())
	cLogTerm := int(req.GetLogTerm())
	cId := int(req.GetCandidateId())
	cLogLength := int(req.GetLogLength())

	if cTerm > r.node.CurrentTerm {
		r.node.CurrentTerm = cTerm
		r.node.CurrentRole = node.Follower
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

func (r *RaftServer) AppendEntries(ctx context.Context, req *pb.AppendEntryRequest) (*pb.AppendEntryResponse, error) {

}
