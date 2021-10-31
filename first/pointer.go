package main

import (
	"errors"
	"fmt"
)

var ErrInsufficientFunds = errors.New("balance is not enough")

type Wallet struct {
	balance Bitcoin
}

type Bitcoin int

type Stringer interface {
    String() string
}

func (b Bitcoin) String() string {
    return fmt.Sprintf("%d BTC", b)
}

func (w *Wallet) Balance() Bitcoin {
	return w.balance
}

func (w *Wallet) Deposit(amount Bitcoin) {
	w.balance += amount
}

func (w *Wallet) Withdraw(amount Bitcoin) error {
	if w.balance >= amount {
		w.balance -= amount
		return nil
	} else {
		return ErrInsufficientFunds
	}
}
