package handler

import (
	"context"
	"errors"
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
	return nil, errors.New("unimplemented")
}

func (r *RaftServer) AppendEntries(ctx context.Context, req *pb.AppendEntryRequest) (*pb.AppendEntryResponse, error) {

}
