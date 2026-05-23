package bankingsystem

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/Bruwald/CodeSignalPractice/utils"
	"github.com/shopspring/decimal"
)

type (
	TransactionType string
	PaymentStatus   string
	AccountID       string
	PaymentID       string

	BankingSystem interface {
		// Creates a new account with the given identifier if it doesn’t already exist.
		// Returns true if the account was successfully created or false if an account
		// with account_id already exists.
		CreateAccount(timestamp int64, accountID AccountID) bool
		// Deposits the given amount of money to the specified account account_id.
		// Returns the balance of the account after the operation has been processed.
		// If the specified account doesn’t exist, should return nil.
		Deposit(timestamp int64, accountID AccountID, amount decimal.Decimal) *decimal.Decimal
		// Transfer the given amount of money from account source_account_id to account target_account_id.
		// Returns the balance of source_account_id if the transfer was successful or nil otherwise.
		// Returns nil if source_account_id or target_account_id doesn’t exist.
		// Returns nil if source_account_id and target_account_id are the same.
		// Returns nil if account source_account_id has insufficient funds to perform the transfer.
		Transfer(timestamp int64, sourceAccountID AccountID, targetAccountID AccountID, amount decimal.Decimal) *decimal.Decimal
		// Return the identifiers of the top n accounts with the highest outgoing transactions
		// - the total amount of money either transferred out of or paid/withdrawn (the pay operation will be introduced in level 3)
		// - sorted in descending order, or in case of a tie, sorted alphabetically by account_id in ascending order.
		// The result should be a list of strings in the following format: ["<account_id_1>(<total_outgoing_1>)", "<account_id_2>(<total_outgoing_2>)", ..., "<account_id_n>(<total_outgoing_n>)"].
		TopSpenders(timestamp int64, n int) []string
		// should withdraw the given amount of money from the specified account.
		// All withdraw transactions provide a 2% cashback – 2% of the withdrawn amount (rounded down to the nearest integer) will be refunded
		// to the account 24 hours after the withdrawal.
		// If the withdrawal is successful, returns a string with a unique identifier for the payment transaction
		// in this format: "payment[ordinal number of withdraws from all accounts]" — e.g., "payment1", "payment2", etc.
		Pay(timestamp int64, accountID AccountID, amount decimal.Decimal) *string
		// Return the status of the payment transaction for the given payment.
		GetPaymentStatus(timestamp int64, accountID AccountID, paymentID PaymentID) *string
		// Merge sourceAccountID into the targetAccountID.
		// Returns True if accounts were successfully merged, or False otherwise.
		MergeAccounts(timestamp int64, sourceAccountID AccountID, targetAccountID AccountID) bool
		// return the total amount of money in the account account_id at the given timestamp time_at.
		// If the specified account did not exist at a given time time_at, returns nil.
		GetBalance(timestamp int64, accountID AccountID) *decimal.Decimal
		// Prints the current state of the banking system, including all accounts and their balances.
		Print()
	}

	SpenderRank struct {
		AccountID  AccountID
		TotalSpent decimal.Decimal
	}

	Transaction struct {
		Timestamp int64
		Type      TransactionType
		Amount    decimal.Decimal
	}

	Payment struct {
		Timestamp      int64
		PaymentID      PaymentID
		AccountID      AccountID
		Amount         decimal.Decimal
		CashbackAmount decimal.Decimal
		Status         PaymentStatus
		DueTime        int64
	}

	Account struct {
		AccountID    AccountID
		Timestamp    int64
		Balance      decimal.Decimal
		Transactions []Transaction
		Payments     map[PaymentID]*Payment
	}

	PendingCashback struct {
		Timestamp int64
		AccountID AccountID
		PaymentID PaymentID
		Amount    decimal.Decimal
	}

	BankingSystemImpl struct {
		Accounts         map[AccountID]*Account
		PaymentCount     int64
		PendingCashbacks []PendingCashback
	}
)

const (
	TransactionTypeOutgoingTransfer TransactionType = "outgoing_transfer"
	TransactionTypeIncomingTransfer TransactionType = "incoming_transfer"
	TransactionTypeIncomingPayment  TransactionType = "incoming_payment"
	TransactionTypeOutgoingPayment  TransactionType = "outgoing_payment"
	TransactionTypeCashback         TransactionType = "cashback"
	TransactionTypeDeposit          TransactionType = "deposit"

	PaymentStatusInProgress       PaymentStatus = "in_progress"
	PaymentStatusCashbackReceived PaymentStatus = "cashback_received"

	PaymentIDPrefix       string = "payment"
	CashbackWaitingPeriod int64  = 24 * 60 * 60 * 1000 // 24 hours in milliseconds
)

func NewBankingSystem() BankingSystem {
	return &BankingSystemImpl{
		Accounts:         map[AccountID]*Account{},
		PaymentCount:     int64(0),
		PendingCashbacks: []PendingCashback{},
	}
}

func (b *BankingSystemImpl) CreateAccount(timestamp int64, accountID AccountID) bool {
	if _, exists := b.Accounts[accountID]; exists {
		return false
	}
	b.Accounts[accountID] = &Account{
		AccountID:    accountID,
		Timestamp:    timestamp,
		Balance:      decimal.Zero,
		Transactions: []Transaction{},
		Payments:     map[PaymentID]*Payment{},
	}
	return true
}

func (b *BankingSystemImpl) Deposit(timestamp int64, accountID AccountID, amount decimal.Decimal) *decimal.Decimal {
	b.processCashbacks(timestamp)

	account, exists := b.Accounts[accountID]
	if !exists {
		return nil
	}
	account.Balance = account.Balance.Add(amount)
	account.Transactions = append(account.Transactions, Transaction{
		Timestamp: timestamp,
		Type:      TransactionTypeDeposit,
		Amount:    amount,
	})
	return &account.Balance
}

func (b *BankingSystemImpl) Transfer(timestamp int64, sourceAccountID AccountID, targetAccountID AccountID, amount decimal.Decimal) *decimal.Decimal {
	b.processCashbacks(timestamp)

	sourceAccount, sourceAcountExists := b.Accounts[sourceAccountID]
	destinationAccount, destinationAccountExists := b.Accounts[targetAccountID]
	if !sourceAcountExists ||
		!destinationAccountExists ||
		sourceAccount == destinationAccount ||
		sourceAccount.Balance.Sub(amount).LessThan(decimal.Zero) {
		return nil
	}

	sourceAccount.Balance = sourceAccount.Balance.Sub(amount)
	destinationAccount.Balance = destinationAccount.Balance.Add(amount)
	sourceAccount.Transactions = append(sourceAccount.Transactions, Transaction{
		Timestamp: timestamp,
		Type:      TransactionTypeOutgoingTransfer,
		Amount:    amount,
	})
	destinationAccount.Transactions = append(destinationAccount.Transactions, Transaction{
		Timestamp: timestamp,
		Type:      TransactionTypeIncomingTransfer,
		Amount:    amount,
	})

	return &sourceAccount.Balance
}

func (b *BankingSystemImpl) TopSpenders(timestamp int64, n int) []string {
	b.processCashbacks(timestamp)

	if n <= 0 {
		return []string{}
	}

	spenderRanks := []SpenderRank{}

	for _, account := range b.Accounts {
		outgoingTransactionSum := decimal.Zero
		for _, transaction := range account.Transactions {
			if transaction.Type == TransactionTypeOutgoingTransfer ||
				transaction.Type == TransactionTypeOutgoingPayment {
				outgoingTransactionSum = outgoingTransactionSum.Add(transaction.Amount)
			}
		}
		spenderRanks = append(spenderRanks, SpenderRank{
			AccountID:  account.AccountID,
			TotalSpent: outgoingTransactionSum,
		})
	}

	sort.Slice(spenderRanks, func(i, j int) bool {
		if spenderRanks[i].TotalSpent != spenderRanks[j].TotalSpent {
			return spenderRanks[i].TotalSpent.GreaterThan(spenderRanks[j].TotalSpent)
		}
		return spenderRanks[i].AccountID < spenderRanks[j].AccountID
	})

	if len(spenderRanks) > n {
		spenderRanks = spenderRanks[:n]
	}

	topSpenders := make([]string, 0, n)
	for _, spenderRank := range spenderRanks {
		topSpenders = append(topSpenders,
			string(spenderRank.AccountID)+"("+spenderRank.TotalSpent.StringFixed(2)+")")
	}

	return topSpenders
}

func (b *BankingSystemImpl) Pay(timestamp int64, accountID AccountID, amount decimal.Decimal) *string {
	b.processCashbacks(timestamp)

	account, exists := b.Accounts[accountID]
	if !exists ||
		account.Balance.Sub(amount).LessThan(decimal.Zero) {
		return nil
	}

	b.PaymentCount++
	account.Balance = account.Balance.Sub(amount)

	paymentID := PaymentIDPrefix + strconv.FormatInt(b.PaymentCount, 10)
	cashbackAmount := amount.Mul(decimal.NewFromFloat32(0.02)).Round(2)
	account.Payments[PaymentID(paymentID)] = &Payment{
		Timestamp:      timestamp,
		PaymentID:      PaymentID(paymentID),
		AccountID:      accountID,
		Amount:         amount,
		CashbackAmount: cashbackAmount,
		Status:         PaymentStatusInProgress,
		DueTime:        timestamp + CashbackWaitingPeriod,
	}
	account.Transactions = append(account.Transactions, Transaction{
		Timestamp: timestamp,
		Type:      TransactionTypeOutgoingPayment,
		Amount:    amount,
	})

	b.PendingCashbacks = append(b.PendingCashbacks, PendingCashback{
		Timestamp: timestamp + CashbackWaitingPeriod,
		AccountID: accountID,
		PaymentID: PaymentID(paymentID),
		Amount:    cashbackAmount,
	})

	return &paymentID
}

func (b *BankingSystemImpl) GetPaymentStatus(timestamp int64, accountID AccountID, paymentID PaymentID) *string {
	b.processCashbacks(timestamp)

	account, accountExists := b.Accounts[accountID]
	if !accountExists {
		return nil
	}

	payment, paymentExists := account.Payments[paymentID]
	if !paymentExists || payment.AccountID != accountID {
		return nil
	}

	return utils.GetPointer(string(payment.Status))
}

func (b *BankingSystemImpl) MergeAccounts(timestamp int64, sourceAccountID AccountID, targetAccountID AccountID) bool {
	b.processCashbacks(timestamp)

	sourceAccount, sourceExists := b.Accounts[sourceAccountID]
	targetAccount, targetExists := b.Accounts[targetAccountID]

	if !sourceExists || !targetExists {
		return false
	}

	sourceAccount.Balance = sourceAccount.Balance.Add(targetAccount.Balance)

	for _, transaction := range targetAccount.Transactions {
		sourceAccount.Transactions = append(sourceAccount.Transactions, transaction)
	}

	for paymentID, payment := range targetAccount.Payments {
		payment.AccountID = sourceAccountID
		sourceAccount.Payments[paymentID] = payment
	}

	for _, pendingCashback := range b.PendingCashbacks {
		if pendingCashback.AccountID == targetAccountID {
			pendingCashback.AccountID = sourceAccountID
		}
	}

	delete(b.Accounts, targetAccountID)

	return true
}

func (b *BankingSystemImpl) GetBalance(timestamp int64, accountID AccountID) *decimal.Decimal {
	b.processCashbacks(timestamp)

	account, exists := b.Accounts[accountID]
	if !exists {
		return nil
	}

	// Just to make sure it is sorted by time for the balance calculation loop.
	sort.Slice(account.Transactions, func(i, j int) bool {
		return account.Transactions[i].Timestamp < account.Transactions[j].Timestamp
	})

	balanceUpUntilTimestamp := decimal.Zero
	for _, transaction := range account.Transactions {
		if transaction.Timestamp < timestamp {
			switch transaction.Type {
			case TransactionTypeDeposit,
				TransactionTypeIncomingTransfer,
				TransactionTypeIncomingPayment,
				TransactionTypeCashback:
				balanceUpUntilTimestamp = balanceUpUntilTimestamp.Add(transaction.Amount)
			case TransactionTypeOutgoingTransfer,
				TransactionTypeOutgoingPayment:
				balanceUpUntilTimestamp = balanceUpUntilTimestamp.Sub(transaction.Amount)
			}
		} else {
			break
		}
	}

	return &balanceUpUntilTimestamp
}

func (b *BankingSystemImpl) Print() {
	fmt.Println("\nCurrent banking system:")
	if len(b.Accounts) == 0 {
		fmt.Println("  No accounts created yet")
		return
	}
	for _, account := range b.Accounts {
		fmt.Printf("  Account ID: %s, Balance: %s, Timestamp: %d\n",
			account.AccountID, account.Balance.StringFixed(2), account.Timestamp)

		if len(account.Transactions) > 0 {
			fmt.Println("    Transactions:")
			for _, txn := range account.Transactions {
				fmt.Printf("      Type: %s, Amount: %s, Timestamp: %d\n",
					txn.Type, txn.Amount.StringFixed(2), txn.Timestamp)
			}
		}

		if len(account.Payments) > 0 {
			fmt.Println("    Payments:")
			for paymentID, payment := range account.Payments {
				fmt.Printf("      Payment ID: %s, Status: %s, Amount: %s\n",
					paymentID, payment.Status, payment.Amount.StringFixed(2))
			}
		}
	}

	if len(b.PendingCashbacks) > 0 {
		fmt.Println("  Pending Cashbacks:")
		for _, cashback := range b.PendingCashbacks {
			fmt.Printf("    Account ID: %s, Payment ID: %s, Amount: %s, Timestamp: %d\n",
				cashback.AccountID, cashback.PaymentID, cashback.Amount.StringFixed(2), cashback.Timestamp)
		}
	}
}

func (b *BankingSystemImpl) processCashbacks(timestamp int64) {
	// Just to make sure it is sorted by time for the cashback processing loop,
	// in case there were some merges that changed account IDs for pending cashbacks.
	sort.Slice(b.PendingCashbacks, func(i, j int) bool {
		return b.PendingCashbacks[i].Timestamp < b.PendingCashbacks[j].Timestamp
	})

	i := 0
	for _, pendingCashback := range b.PendingCashbacks {
		if timestamp > pendingCashback.Timestamp {
			i++
			account, exists := b.Accounts[pendingCashback.AccountID]
			if !exists {
				continue
			}
			account.Balance = account.Balance.Add(pendingCashback.Amount)
			account.Transactions = append(account.Transactions, Transaction{
				Timestamp: timestamp,
				Type:      TransactionTypeCashback,
				Amount:    pendingCashback.Amount,
			})

			if payment, paymentExists := account.Payments[pendingCashback.PaymentID]; paymentExists {
				payment.Status = PaymentStatusCashbackReceived
			}
		}
	}
	if i > 0 {
		b.PendingCashbacks = b.PendingCashbacks[i:]
	}
}
