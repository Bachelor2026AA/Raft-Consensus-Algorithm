package raft

func vote(
	term, votedFor int, logs []Log,
	cTerm, cLogTerm, cLogLength, cId int,
) bool {
	lastTerm := 0
	if len(logs) > 0 {
		lastTerm = logs[len(logs)-1].Term
	}
	logOk := cLogTerm > lastTerm || (cLogTerm == lastTerm && cLogLength >= len(logs))
	return cTerm == term && logOk && (votedFor == -1 || votedFor == cId)
}
