package handler

import (
	"context"
	"errors"
	pb "raft-consensus/proto"
)

type raftServer struct {
	pb.UnimplementedRaftServer
}

func (r *raftServer) RequestVote(ctx context.Context, req *pb.RequestVoteRequest) (*pb.RequestVoteResponse, error) {
	return nil, errors.New("unimplemented")
}
