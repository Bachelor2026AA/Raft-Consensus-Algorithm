package handler

import (
	"context"
	"errors"
	pb "raft-consensus/proto"
)

type RaftServer struct {
	pb.UnimplementedRaftServer
}

func (r *RaftServer) RequestVote(ctx context.Context, req *pb.RequestVoteRequest) (*pb.RequestVoteResponse, error) {
	return nil, errors.New("unimplemented")
}

func (r *RaftServer) AppendEntries(ctx context.Context, req *pb.AppendEntryRequest) (*pb.AppendEntryResponse, error) {
	return nil, errors.New("unimplemented")
}
