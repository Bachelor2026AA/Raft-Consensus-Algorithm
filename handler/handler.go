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
	
