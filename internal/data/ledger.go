package data

import "fmt"

type Ledger struct {
	transactions []Transaction
}

func NewLedger() *Ledger {
	return &Ledger{
		transactions: make([]Transaction, 0),
	}
}

func (l *Ledger) AddTransaction(t Transaction) {
	l.transactions = append(l.transactions, t)
}

func (l *Ledger) AddTransactions(tSlice []Transaction) {
	l.transactions = append(l.transactions, tSlice...)
}

func (l *Ledger) GetTransaction(id int) (Transaction, error) {
	for _, t := range l.transactions {
		if t.ID == id {
			return t, nil
		}
	}
	return Transaction{}, fmt.Errorf("no transaction with ID %d was found", id)
}

func (l *Ledger) GetTransactions() []Transaction {
	return l.transactions
}
