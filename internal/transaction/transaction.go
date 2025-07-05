package transaction

import (
	"lesta-battleship/server-core/internal/game"
	// "lesta-battleship/server-core/internal/match"
)

type Command interface {
	Apply(states *game.States) error
	Undo(states *game.States)
}

type Transaction struct {
	commands []Command
}

func NewTransaction() *Transaction {
	return &Transaction{}
}

func (tx *Transaction) Add(cmd Command) {
	tx.commands = append(tx.commands, cmd)
}

func (tx *Transaction) Execute(states *game.States) error {
	for i, cmd := range tx.commands {
		if err := cmd.Apply(states); err != nil {
			for j := i - 1; j >= 0; j-- {
				tx.commands[j].Undo(states)
			}
			return err
		}
	}
	return nil
}
