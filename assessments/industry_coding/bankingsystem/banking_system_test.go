package bankingsystem

import (
	"testing"

	"github.com/Bruwald/CodeSignalPractice/utils"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestBankingSystem_CreateAccount(t *testing.T) {
	bankingSystem := NewBankingSystem()

	testCases := []struct {
		name      string
		timestamp int64
		accountID AccountID
		expected  bool
	}{
		{
			name:      "Should create a new account successfully",
			timestamp: 1,
			accountID: "account1",
			expected:  true,
		},
		{
			name:      "Should not create a duplicate account",
			timestamp: 2,
			accountID: "account1",
			expected:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := bankingSystem.CreateAccount(tc.timestamp, tc.accountID)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestBankingSystem_Deposit(t *testing.T) {
	bankingSystem := NewBankingSystem()

	testCases := []struct {
		name          string
		bankingSystem BankingSystem
		timestamp     int64
		accountID     AccountID
		amount        decimal.Decimal
		setup         func()
		cleanup       func()
		expected      *decimal.Decimal
	}{
		{
			name:      "Should deposit to existing account",
			timestamp: 1,
			accountID: "account1",
			amount:    decimal.NewFromFloat(100.0),
			setup: func() {
				_ = bankingSystem.CreateAccount(1, "account1")
			},
			cleanup: func() {
				bankingSystem = NewBankingSystem()
			},
			expected: utils.GetPointer(decimal.NewFromFloat(100.0)),
		},
		{
			name:      "Should not deposit to non-existing account",
			timestamp: 2,
			accountID: "account2",
			amount:    decimal.NewFromFloat(50.0),
			setup:     func() {},
			cleanup:   func() {},
			expected:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()
			result := bankingSystem.Deposit(tc.timestamp, tc.accountID, tc.amount)
			if tc.expected == nil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
				assert.True(t, result.Equal(*tc.expected))
			}
			tc.cleanup()
		})
	}
}

func TestBankingSystem_Transfer(t *testing.T) {
	bankingSystem := NewBankingSystem()

	testCases := []struct {
		name            string
		timestamp       int64
		sourceAccountID AccountID
		targetAccountID AccountID
		amount          decimal.Decimal
		setup           func()
		cleanup         func()
		expected        *decimal.Decimal
	}{
		{
			name:            "Should transfer between accounts successfully",
			timestamp:       1,
			sourceAccountID: "account1",
			targetAccountID: "account2",
			amount:          decimal.NewFromFloat(50.0),
			setup: func() {
				_ = bankingSystem.CreateAccount(1, "account1")
				_ = bankingSystem.CreateAccount(2, "account2")
				_ = bankingSystem.Deposit(1, "account1", decimal.NewFromFloat(100.0))
			},
			cleanup: func() {
				bankingSystem = NewBankingSystem()
			},
			expected: utils.GetPointer(decimal.NewFromFloat(50.0)),
		},
		{
			name:            "Should not transfer if source account has insufficient funds",
			timestamp:       2,
			sourceAccountID: "account1",
			targetAccountID: "account2",
			amount:          decimal.NewFromFloat(150.0),
			setup: func() {
				_ = bankingSystem.CreateAccount(1, "account1")
				_ = bankingSystem.CreateAccount(2, "account2")
				_ = bankingSystem.Deposit(1, "account1", decimal.NewFromFloat(100.0))
			},
			cleanup: func() {
				bankingSystem = NewBankingSystem()
			},
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()
			result := bankingSystem.Transfer(tc.timestamp, tc.sourceAccountID, tc.targetAccountID, tc.amount)
			if tc.expected == nil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
				assert.True(t, result.Equal(*tc.expected))
			}
			tc.cleanup()
		})
	}
}

func TestBankingSystem_TopSpenders(t *testing.T) {
	bankingSystem := NewBankingSystem()

	testCases := []struct {
		name          string
		bankingSystem BankingSystem
		timestamp     int64
		n             int
		setup         func()
		cleanup       func()
		expected      []string
	}{
		{
			name:      "Should return top spenders",
			timestamp: 1,
			n:         10,
			setup: func() {
				_ = bankingSystem.CreateAccount(1, "account1")
				_ = bankingSystem.CreateAccount(2, "account2")
				_ = bankingSystem.CreateAccount(3, "account3")
				_ = bankingSystem.Deposit(1, "account1", decimal.NewFromFloat(100.0))
				_ = bankingSystem.Deposit(2, "account2", decimal.NewFromFloat(100.0))
				_ = bankingSystem.Transfer(1, "account1", "account2", decimal.NewFromFloat(50.0))
				_ = bankingSystem.Transfer(1, "account2", "account1", decimal.NewFromFloat(50.0))
			},
			cleanup: func() {
				bankingSystem = NewBankingSystem()
			},
			expected: []string{"account1(50.00)", "account2(50.00)", "account3(0.00)"},
		},
		{
			name:      "Should return top 2 spenders even when there are 3 spender accounts",
			timestamp: 1,
			n:         2,
			setup: func() {
				_ = bankingSystem.CreateAccount(1, "account1")
				_ = bankingSystem.CreateAccount(2, "account2")
				_ = bankingSystem.CreateAccount(3, "account3")
				_ = bankingSystem.Deposit(3, "account1", decimal.NewFromFloat(100.0))
				_ = bankingSystem.Deposit(4, "account2", decimal.NewFromFloat(100.0))
				_ = bankingSystem.Deposit(4, "account3", decimal.NewFromFloat(20.0))
				_ = bankingSystem.Transfer(5, "account1", "account2", decimal.NewFromFloat(100.0))
				_ = bankingSystem.Transfer(6, "account2", "account3", decimal.NewFromFloat(50.0))

			},
			cleanup: func() {
				bankingSystem = NewBankingSystem()
			},
			expected: []string{"account1(100.00)", "account2(50.00)"},
		},
		{
			name:      "Should return empty list if no spenders",
			timestamp: 2,
			n:         2,
			setup:     func() {},
			cleanup:   func() {},
			expected:  []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()
			result := bankingSystem.TopSpenders(tc.timestamp, tc.n)
			assert.Equal(t, tc.expected, result)
			tc.cleanup()
		})
	}
}

func TestBankingSystem_Pay(t *testing.T) {
	bankingSystem := NewBankingSystem()

	testCases := []struct {
		name      string
		timestamp int64
		accountID AccountID
		amount    decimal.Decimal
		setup     func()
		cleanup   func()
		expected  *string
	}{
		{
			name:      "Should withdraw and return payment ID",
			timestamp: 1,
			accountID: "account1",
			amount:    decimal.NewFromFloat(50.0),
			setup: func() {
				_ = bankingSystem.CreateAccount(1, "account1")
				_ = bankingSystem.Deposit(1, "account1", decimal.NewFromFloat(100.0))
			},
			cleanup: func() {
				bankingSystem = NewBankingSystem()
			},
			expected: utils.GetPointer("payment1"),
		},
		{
			name:      "Should not withdraw from non-existing account",
			timestamp: 2,
			accountID: "account2",
			amount:    decimal.NewFromFloat(50.0),
			setup:     func() {},
			cleanup:   func() {},
			expected:  nil,
		},
		{
			name:      "Should not withdraw if account has insufficient funds",
			timestamp: 3,
			accountID: "account1",
			amount:    decimal.NewFromFloat(150.0),
			setup: func() {
				_ = bankingSystem.CreateAccount(1, "account1")
				_ = bankingSystem.Deposit(1, "account1", decimal.NewFromFloat(100.0))
			},
			cleanup: func() {
				bankingSystem = NewBankingSystem()
			},
			expected: nil,
		},
		{
			name:      "Should increment payment ID for multiple payments",
			timestamp: 4,
			accountID: "account1",
			amount:    decimal.NewFromFloat(25.0),
			setup: func() {
				_ = bankingSystem.CreateAccount(1, "account1")
				_ = bankingSystem.Deposit(1, "account1", decimal.NewFromFloat(200.0))
				_ = bankingSystem.Pay(1, "account1", decimal.NewFromFloat(50.0))
			},
			cleanup: func() {
				bankingSystem = NewBankingSystem()
			},
			expected: utils.GetPointer("payment2"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()
			result := bankingSystem.Pay(tc.timestamp, tc.accountID, tc.amount)
			if tc.expected == nil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
				assert.Equal(t, *tc.expected, *result)
			}
			tc.cleanup()
		})
	}
}

func TestBankingSystem_GetPaymentStatus(t *testing.T) {
	bankingSystem := NewBankingSystem()

	testCases := []struct {
		name      string
		timestamp int64
		accountID AccountID
		paymentID PaymentID
		setup     func()
		cleanup   func()
		expected  *string
	}{
		{
			name:      "Should return payment status in progress",
			timestamp: 1,
			accountID: "account1",
			paymentID: "payment1",
			setup: func() {
				_ = bankingSystem.CreateAccount(1, "account1")
				_ = bankingSystem.Deposit(1, "account1", decimal.NewFromFloat(100.0))
				_ = bankingSystem.Pay(1, "account1", decimal.NewFromFloat(50.0))
			},
			cleanup: func() {
				bankingSystem = NewBankingSystem()
			},
			expected: utils.GetPointer(string(PaymentStatusInProgress)),
		},
		{
			name:      "Should not get status for non-existing account",
			timestamp: 2,
			accountID: "account2",
			paymentID: "payment1",
			setup:     func() {},
			cleanup:   func() {},
			expected:  nil,
		},
		{
			name:      "Should not get status for non-existing payment",
			timestamp: 3,
			accountID: "account1",
			paymentID: "payment999",
			setup: func() {
				_ = bankingSystem.CreateAccount(1, "account1")
				_ = bankingSystem.Deposit(1, "account1", decimal.NewFromFloat(100.0))
				_ = bankingSystem.Pay(1, "account1", decimal.NewFromFloat(50.0))
			},
			cleanup: func() {
				bankingSystem = NewBankingSystem()
			},
			expected: nil,
		},
		{
			name:      "Should return cashback received status after waiting period",
			timestamp: 1000000000000,
			accountID: "account1",
			paymentID: "payment1",
			setup: func() {
				_ = bankingSystem.CreateAccount(1, "account1")
				_ = bankingSystem.Deposit(1, "account1", decimal.NewFromFloat(100.0))
				_ = bankingSystem.Pay(1, "account1", decimal.NewFromFloat(50.0))
			},
			cleanup: func() {
				bankingSystem = NewBankingSystem()
			},
			expected: utils.GetPointer(string(PaymentStatusCashbackReceived)),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()
			result := bankingSystem.GetPaymentStatus(tc.timestamp, tc.accountID, tc.paymentID)
			if tc.expected == nil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
				assert.Equal(t, *tc.expected, *result)
			}
			tc.cleanup()
		})
	}
}

func TestBankingSystem_MergeAccounts(t *testing.T) {
	bankingSystem := NewBankingSystem()

	testCases := []struct {
		name            string
		timestamp       int64
		sourceAccountID AccountID
		targetAccountID AccountID
		setup           func()
		cleanup         func()
		expected        bool
		verifyState     func(t *testing.T)
	}{
		{
			name:            "Should merge accounts successfully",
			timestamp:       1,
			sourceAccountID: "account1",
			targetAccountID: "account2",
			setup: func() {
				_ = bankingSystem.CreateAccount(1, "account1")
				_ = bankingSystem.CreateAccount(2, "account2")
				_ = bankingSystem.Deposit(1, "account1", decimal.NewFromFloat(100.0))
				_ = bankingSystem.Deposit(2, "account2", decimal.NewFromFloat(50.0))
			},
			cleanup: func() {
				bankingSystem = NewBankingSystem()
			},
			expected: true,
			verifyState: func(t *testing.T) {
				impl := bankingSystem.(*BankingSystemImpl)
				assert.NotNil(t, impl.Accounts["account1"])
				assert.Nil(t, impl.Accounts["account2"])
				assert.True(t, impl.Accounts["account1"].Balance.Equal(decimal.NewFromFloat(150.0)))
			},
		},
		{
			name:            "Should not merge if source account does not exist",
			timestamp:       2,
			sourceAccountID: "account_missing",
			targetAccountID: "account1",
			setup: func() {
				_ = bankingSystem.CreateAccount(1, "account1")
				_ = bankingSystem.Deposit(1, "account1", decimal.NewFromFloat(100.0))
			},
			cleanup: func() {
				bankingSystem = NewBankingSystem()
			},
			expected: false,
		},
		{
			name:            "Should not merge if target account does not exist",
			timestamp:       3,
			sourceAccountID: "account1",
			targetAccountID: "account_missing",
			setup: func() {
				_ = bankingSystem.CreateAccount(1, "account1")
				_ = bankingSystem.Deposit(1, "account1", decimal.NewFromFloat(100.0))
			},
			cleanup: func() {
				bankingSystem = NewBankingSystem()
			},
			expected: false,
		},
		{
			name:            "Should merge and combine transactions",
			timestamp:       4,
			sourceAccountID: "account1",
			targetAccountID: "account2",
			setup: func() {
				_ = bankingSystem.CreateAccount(1, "account1")
				_ = bankingSystem.CreateAccount(2, "account2")
				_ = bankingSystem.Deposit(1, "account1", decimal.NewFromFloat(100.0))
				_ = bankingSystem.Deposit(2, "account2", decimal.NewFromFloat(50.0))
				_ = bankingSystem.Transfer(3, "account1", "account2", decimal.NewFromFloat(25.0))
			},
			cleanup: func() {
				bankingSystem = NewBankingSystem()
			},
			expected: true,
			verifyState: func(t *testing.T) {
				impl := bankingSystem.(*BankingSystemImpl)
				sourceAccount := impl.Accounts["account1"]
				assert.NotNil(t, sourceAccount)
				assert.Greater(t, len(sourceAccount.Transactions), 2)
				assert.Nil(t, impl.Accounts["account2"])
			},
		},
		{
			name:            "Should merge and transfer payments",
			timestamp:       5,
			sourceAccountID: "account1",
			targetAccountID: "account2",
			setup: func() {
				_ = bankingSystem.CreateAccount(1, "account1")
				_ = bankingSystem.CreateAccount(2, "account2")
				_ = bankingSystem.Deposit(1, "account1", decimal.NewFromFloat(100.0))
				_ = bankingSystem.Deposit(2, "account2", decimal.NewFromFloat(100.0))
				_ = bankingSystem.Pay(1, "account1", decimal.NewFromFloat(10.0))
				_ = bankingSystem.Pay(2, "account2", decimal.NewFromFloat(20.0))
			},
			cleanup: func() {
				bankingSystem = NewBankingSystem()
			},
			expected: true,
			verifyState: func(t *testing.T) {
				impl := bankingSystem.(*BankingSystemImpl)
				sourceAccount := impl.Accounts["account1"]
				assert.GreaterOrEqual(t, len(sourceAccount.Payments), 2)
				for _, payment := range sourceAccount.Payments {
					assert.Equal(t, "account1", string(payment.AccountID))
				}
				assert.Nil(t, impl.Accounts["account2"])
			},
		},
		{
			name:            "Should merge zero balance accounts",
			timestamp:       6,
			sourceAccountID: "account1",
			targetAccountID: "account2",
			setup: func() {
				_ = bankingSystem.CreateAccount(1, "account1")
				_ = bankingSystem.CreateAccount(2, "account2")
				_ = bankingSystem.Deposit(1, "account1", decimal.NewFromFloat(100.0))
			},
			cleanup: func() {
				bankingSystem = NewBankingSystem()
			},
			expected: true,
			verifyState: func(t *testing.T) {
				impl := bankingSystem.(*BankingSystemImpl)
				assert.True(t, impl.Accounts["account1"].Balance.Equal(decimal.NewFromFloat(100.0)))
				assert.Nil(t, impl.Accounts["account2"])
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()
			result := bankingSystem.MergeAccounts(tc.timestamp, tc.sourceAccountID, tc.targetAccountID)
			assert.Equal(t, tc.expected, result)
			if tc.verifyState != nil {
				tc.verifyState(t)
			}
			tc.cleanup()
		})
	}
}

func TestBankingSystem_GetBalance(t *testing.T) {
	bankingSystem := NewBankingSystem()

	testCases := []struct {
		name      string
		timestamp int64
		accountID AccountID
		setup     func()
		cleanup   func()
		expected  *decimal.Decimal
	}{
		{
			name:      "Should return nil for non-existing account",
			timestamp: 1,
			accountID: "account_missing",
			setup:     func() {},
			cleanup:   func() {},
			expected:  nil,
		},
		{
			name:      "Should get balance after deposit",
			timestamp: 2,
			accountID: "account1",
			setup: func() {
				_ = bankingSystem.CreateAccount(1, "account1")
				_ = bankingSystem.Deposit(1, "account1", decimal.NewFromFloat(100.0))
			},
			cleanup: func() {
				bankingSystem = NewBankingSystem()
			},
			expected: utils.GetPointer(decimal.NewFromFloat(100.0)),
		},
		{
			name:      "Should only count transactions before timestamp",
			timestamp: 2,
			accountID: "account1",
			setup: func() {
				_ = bankingSystem.CreateAccount(1, "account1")
				_ = bankingSystem.Deposit(1, "account1", decimal.NewFromFloat(100.0))
				_ = bankingSystem.Deposit(5, "account1", decimal.NewFromFloat(50.0))
			},
			cleanup: func() {
				bankingSystem = NewBankingSystem()
			},
			expected: utils.GetPointer(decimal.NewFromFloat(100.0)),
		},
		{
			name:      "Should get balance after transfer out",
			timestamp: 3,
			accountID: "account1",
			setup: func() {
				_ = bankingSystem.CreateAccount(1, "account1")
				_ = bankingSystem.CreateAccount(2, "account2")
				_ = bankingSystem.Deposit(1, "account1", decimal.NewFromFloat(100.0))
				_ = bankingSystem.Transfer(2, "account1", "account2", decimal.NewFromFloat(30.0))
			},
			cleanup: func() {
				bankingSystem = NewBankingSystem()
			},
			expected: utils.GetPointer(decimal.NewFromFloat(70.0)),
		},
		{
			name:      "Should get balance after transfer in",
			timestamp: 3,
			accountID: "account2",
			setup: func() {
				_ = bankingSystem.CreateAccount(1, "account1")
				_ = bankingSystem.CreateAccount(2, "account2")
				_ = bankingSystem.Deposit(1, "account1", decimal.NewFromFloat(100.0))
				_ = bankingSystem.Transfer(2, "account1", "account2", decimal.NewFromFloat(30.0))
			},
			cleanup: func() {
				bankingSystem = NewBankingSystem()
			},
			expected: utils.GetPointer(decimal.NewFromFloat(30.0)),
		},
		{
			name:      "Should get balance after payment (withdrawal)",
			timestamp: 3,
			accountID: "account1",
			setup: func() {
				_ = bankingSystem.CreateAccount(1, "account1")
				_ = bankingSystem.Deposit(1, "account1", decimal.NewFromFloat(100.0))
				_ = bankingSystem.Pay(2, "account1", decimal.NewFromFloat(40.0))
			},
			cleanup: func() {
				bankingSystem = NewBankingSystem()
			},
			expected: utils.GetPointer(decimal.NewFromFloat(60.0)),
		},
		{
			name:      "Should get balance after multiple deposits",
			timestamp: 6,
			accountID: "account1",
			setup: func() {
				_ = bankingSystem.CreateAccount(1, "account1")
				_ = bankingSystem.Deposit(1, "account1", decimal.NewFromFloat(100.0))
				_ = bankingSystem.Deposit(3, "account1", decimal.NewFromFloat(75.0))
				_ = bankingSystem.Deposit(5, "account1", decimal.NewFromFloat(25.0))
			},
			cleanup: func() {
				bankingSystem = NewBankingSystem()
			},
			expected: utils.GetPointer(decimal.NewFromFloat(200.0)),
		},
		{
			name:      "Should get balance with multiple transactions",
			timestamp: 10,
			accountID: "account1",
			setup: func() {
				_ = bankingSystem.CreateAccount(1, "account1")
				_ = bankingSystem.Deposit(1, "account1", decimal.NewFromFloat(100.0))
				_ = bankingSystem.CreateAccount(2, "account2")
				_ = bankingSystem.Transfer(3, "account1", "account2", decimal.NewFromFloat(25.0))
				_ = bankingSystem.Deposit(5, "account1", decimal.NewFromFloat(50.0))
				_ = bankingSystem.Pay(7, "account1", decimal.NewFromFloat(20.0))
			},
			cleanup: func() {
				bankingSystem = NewBankingSystem()
			},
			expected: utils.GetPointer(decimal.NewFromFloat(105.0)),
		},
		{
			name:      "Should get zero balance for newly created account",
			timestamp: 1,
			accountID: "account1",
			setup: func() {
				_ = bankingSystem.CreateAccount(1, "account1")
			},
			cleanup: func() {
				bankingSystem = NewBankingSystem()
			},
			expected: utils.GetPointer(decimal.Zero),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()
			result := bankingSystem.GetBalance(tc.timestamp, tc.accountID)
			if tc.expected == nil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
				assert.True(t, result.Equal(*tc.expected))
			}
			tc.cleanup()
		})
	}
}
