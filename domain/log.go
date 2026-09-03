package domain

import "errors"

type Log struct {
	Index   uint64
	Term    int
	Command string
}

func (rl *Log) Validate() error {
	if rl.Index <= 0 {
		return errors.New("Index cant be negative")
	}
	if rl.Term <= 0 {
		return errors.New("Term cant be negative")
	}
	if rl.Command == " " {
		return errors.New("Command cant be negative")
	}
	return nil
}
