package service

import (
	"github.com/Rhymond/go-money"
	"github.com/diezfx/split-app-backend/gen/ent"
	"github.com/google/uuid"
)

//go:generate go run github.com/dmarkham/enumer -type=TransactionType
type TransactionType int

const (
	Expense TransactionType = iota
	Transfer
)

type Transaction struct {
	ID              uuid.UUID
	Name            string
	TransactionType TransactionType
	Amount          *money.Money
	SourceID        string
	TargetIDs       []string
}

type Project struct {
	ID           uuid.UUID
	Name         string
	Transactions []Transaction
	Members      []string
}

func FromEntProject(project *ent.Project) Project {
	transactions := make([]Transaction, len(project.Edges.Transactions))
	for i, t := range project.Edges.Transactions {
		transactions[i] = FromEntTransaction(t)
	}

	return Project{
		ID:           project.ID,
		Name:         project.Name,
		Transactions: transactions,
		Members:      project.Members,
	}
}

func FromEntTransaction(trans *ent.Transaction) Transaction {
	return Transaction{
		ID:   trans.ID,
		Name: trans.Name, Amount: money.New(trans.Amount, money.EUR),
		SourceID:  trans.SourceID,
		TargetIDs: trans.TargetIds,
	}
}
