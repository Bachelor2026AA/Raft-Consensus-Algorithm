package domain

import "errors"

type LogEntry struct {
	Index   int
	Term    int
	Command string
	Key     string
	Value   int
}

func (rl *LogEntry) Validate() error {
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
