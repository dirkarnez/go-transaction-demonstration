package main

import (
	"fmt"
	"sync"
	"time"
)

// BankAccount represents a bank account with a balance
type BankAccount struct {
	balance int
	mu      sync.Mutex
}

// Deposit adds money to the account
func (b *BankAccount) Deposit(amount int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Deposit tying to get lock")
	b.mu.Lock()
	fmt.Println("Deposit get the lock")
	defer b.mu.Unlock()
	b.balance += amount
	fmt.Printf("Deposited: %d, New Balance: %d\n", amount, b.balance)
}

// Withdraw subtracts money from the account
func (b *BankAccount) Withdraw(amount int, wg *sync.WaitGroup) {
	defer wg.Done()
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.balance >= amount {
		b.balance -= amount
		time.Sleep(5 * time.Second)
		fmt.Printf("Withdrew: %d, New Balance: %d\n", amount, b.balance)
	} else {
		fmt.Printf("Failed to withdraw: %d, Insufficient funds. Current Balance: %d\n", amount, b.balance)
	}
}

func main() {
	var wg sync.WaitGroup
	account := BankAccount{balance: 1800}

	// 1800 + 500 - 200 - 800 + 300 - 600 = 1000
	// Simulate concurrent deposits and withdrawals
	transactions := []struct {
		action string
		amount int
	}{
		{"deposit", 500},
		{"withdraw", 200},
		{"withdraw", 800},
		{"deposit", 300},
		{"withdraw", 600},
	}

	for _, transaction := range transactions {
		wg.Add(1)
		if transaction.action == "deposit" {
			go account.Deposit(transaction.amount, &wg)
		} else if transaction.action == "withdraw" {
			go account.Withdraw(transaction.amount, &wg)
		}
	}

	// Wait for all goroutines to finish
	wg.Wait()
	fmt.Printf("Final Balance: %d\n", account.balance)
}
