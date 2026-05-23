package main

import (
	"strconv"

	"github.com/Bruwald/CodeSignalPractice/assessments/industry_coding/bankingsystem"
	"github.com/shopspring/decimal"
)

func main() {
	bankingSystem := bankingsystem.NewBankingSystem()
	queries := [][]string{
		{"CREATE_ACCOUNT", "1", "account1"},
		{"CREATE_ACCOUNT", "2", "account2"},
		{"CREATE_ACCOUNT", "2", "account2"}, // Duplicate account creation
		{"CREATE_ACCOUNT", "3", "account3"},
		{"DEPOSIT", "4", "account1", "1000.00"},
		{"DEPOSIT", "5", "account2", "500.00"},
		{"DEPOSIT", "6", "account3", "200.00"},
		{"TRANSFER", "7", "account1", "account2", "200.00"},
		{"TRANSFER", "8", "account1", "account3", "100.00"},
		{"TRANSFER", "9", "account2", "account3", "200.00"},
		{"PAY", "10", "account1", "100.00"},
		{"PAY", "11", "account2", "200.00"},
		{"PAY", "12", "account3", "300.00"},
		{"DEPOSIT", "100000000000000", "account1", "100.00"}, // For cashback processing
		{"DEPOSIT", "100000000000001", "account2", "200.00"}, // Transactions after cashback
		{"PAY", "100000000000002", "account3", "300.00"},
		{"PAY", "100000000000003", "account1", "100.00"},
		{"PAY", "100000000000004", "account2", "200.00"},
		{"PAY", "100000000000005", "account2", "300.00"},
		{"MERGE", "100000000000006", "account1", "account3"},
	}

	processQueries(queries, bankingSystem)
	bankingSystem.Print()
}

func processQueries(queries [][]string, bankingSystem bankingsystem.BankingSystem) {
	for _, query := range queries {
		cmd := query[0]

		switch cmd {
		case "CREATE_ACCOUNT":
			timestamp, _ := strconv.Atoi(query[1])
			accountID := query[2]
			bankingSystem.CreateAccount(int64(timestamp), bankingsystem.AccountID(accountID))
			_ = bankingSystem.CreateAccount(int64(timestamp), bankingsystem.AccountID(accountID))

		case "DEPOSIT":
			timestamp, _ := strconv.Atoi(query[1])
			accountID := query[2]
			amount, _ := decimal.NewFromString(query[3])
			bankingSystem.Deposit(int64(timestamp), bankingsystem.AccountID(accountID), amount)

		case "TRANSFER":
			timestamp, _ := strconv.Atoi(query[1])
			sourceAccountID := query[2]
			targetAccountID := query[3]
			amount, _ := decimal.NewFromString(query[4])
			bankingSystem.Transfer(int64(timestamp), bankingsystem.AccountID(sourceAccountID), bankingsystem.AccountID(targetAccountID), amount)

		case "PAY":
			timestamp, _ := strconv.Atoi(query[1])
			accountID := query[2]
			amount, _ := decimal.NewFromString(query[3])
			bankingSystem.Pay(int64(timestamp), bankingsystem.AccountID(accountID), amount)

		case "MERGE":
			timestamp, _ := strconv.Atoi(query[1])
			sourceAccountID := query[2]
			targetAccountID := query[3]
			bankingSystem.MergeAccounts(int64(timestamp), bankingsystem.AccountID(sourceAccountID), bankingsystem.AccountID(targetAccountID))
		}
	}
}
