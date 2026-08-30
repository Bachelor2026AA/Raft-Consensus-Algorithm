package raft

import pb "raft-consensus/proto"

func Deliver(log Log) {
	panic("Not implemented")
}

func CommitLogEntries() {
	panic("Not implemented")
}
func ReplicateLog(id int, follower int) {
	panic("Not implemented")
}

func (n *Node) handleAppendEntriesReqeust(follower int, resp *pb.AppendEntriesResponse) {

	if int(resp.GetTerm()) > n.Term {
		n.Term = int(resp.GetTerm())
		n.Role = Follower
		n.VotedFor = -1
		return
	}

	if int(resp.GetTerm()) == n.Term && n.Role == Leader {
		if resp.GetSuccess() == true && resp.GetAck() >= int32(n.AckedLength[follower]) {
			n.SentLength[follower] = int(resp.GetAck())
			n.AckedLength[follower] = int(resp.GetAck())
			CommitLogEntries()
		} else if n.SentLength[follower] > 0 {
			n.SentLength[follower] = n.SentLength[follower] - 1
			ReplicateLog(n.ID, follower)
		}
	}
}
