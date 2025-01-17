package db

import (
	"context"
	"testing"

	"github.com/JMustang/OldBank/util"
	"github.com/stretchr/testify/require"
)

// TestTransfer stores the test transfer data and associated accounts
type TestTransfer struct {
	transfer    Transfer
	fromAccount Account
	toAccount   Account
	params      CreateTransferParams
}

func createRandomTransfer(t *testing.T, fromAccount, toAccount Account) TestTransfer {
	arg := CreateTransferParams{
		FromAccountID: fromAccount.ID,
		ToAccountID:   toAccount.ID,
		Amount:        util.RandomPositiveMoney(), // Transferências devem ser sempre positivas
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)
	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return TestTransfer{
		transfer:    transfer,
		fromAccount: fromAccount,
		toAccount:   toAccount,
		params:      arg,
	}
}

func TestCreateTransfer(t *testing.T) {
	fromAccount := createRandomAccount(t).account
	toAccount := createRandomAccount(t).account

	// Teste básico de criação de transferência
	t.Run("valid transfer", func(t *testing.T) {
		testTransfer := createRandomTransfer(t, fromAccount, toAccount)
		require.NotEmpty(t, testTransfer.transfer)
		compareTransfers(t, testTransfer.params, testTransfer.transfer)
	})

	// Teste de transferência com valor zero
	t.Run("zero amount", func(t *testing.T) {
		arg := CreateTransferParams{
			FromAccountID: fromAccount.ID,
			ToAccountID:   toAccount.ID,
			Amount:        0,
		}

		transfer, err := testQueries.CreateTransfer(context.Background(), arg)
		require.NoError(t, err)
		require.NotEmpty(t, transfer)
		compareTransfers(t, arg, transfer)
	})

	// Teste de transferência para a mesma conta
	t.Run("same account", func(t *testing.T) {
		arg := CreateTransferParams{
			FromAccountID: fromAccount.ID,
			ToAccountID:   fromAccount.ID,
			Amount:        util.RandomPositiveMoney(),
		}

		transfer, err := testQueries.CreateTransfer(context.Background(), arg)
		require.NoError(t, err)
		require.NotEmpty(t, transfer)
		compareTransfers(t, arg, transfer)
	})
}

func TestGetTransfer(t *testing.T) {
	fromAccount := createRandomAccount(t).account
	toAccount := createRandomAccount(t).account
	transfer1 := createRandomTransfer(t, fromAccount, toAccount)

	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.transfer.ID)
	require.NoError(t, err)
	require.NotEmpty(t, transfer2)

	compareTransfers(t, transfer1.params, transfer2)
}

func TestListTransfers(t *testing.T) {
	fromAccount := createRandomAccount(t).account
	toAccount := createRandomAccount(t).account

	// Criar 10 transferências
	for i := 0; i < 5; i++ {
		createRandomTransfer(t, fromAccount, toAccount)
		createRandomTransfer(t, toAccount, fromAccount) // Criar transferências nos dois sentidos
	}

	testCases := []struct {
		name     string
		arg      ListTransfersParams
		expected int
		checkFn  func(t *testing.T, transfers []Transfer)
	}{
		{
			name: "List outgoing transfers",
			arg: ListTransfersParams{
				FromAccountID: fromAccount.ID,
				ToAccountID:   fromAccount.ID,
				Limit:         5,
				Offset:        0,
			},
			expected: 5,
			checkFn: func(t *testing.T, transfers []Transfer) {
				for _, transfer := range transfers {
					require.True(t, transfer.FromAccountID == fromAccount.ID || transfer.ToAccountID == fromAccount.ID)
				}
			},
		},
		{
			name: "List with offset",
			arg: ListTransfersParams{
				FromAccountID: fromAccount.ID,
				ToAccountID:   fromAccount.ID,
				Limit:         5,
				Offset:        5,
			},
			expected: 5,
			checkFn: func(t *testing.T, transfers []Transfer) {
				for _, transfer := range transfers {
					require.True(t, transfer.FromAccountID == fromAccount.ID || transfer.ToAccountID == fromAccount.ID)
				}
			},
		},
		{
			name: "List all transfers",
			arg: ListTransfersParams{
				FromAccountID: fromAccount.ID,
				ToAccountID:   fromAccount.ID,
				Limit:         10,
				Offset:        0,
			},
			expected: 10,
			checkFn: func(t *testing.T, transfers []Transfer) {
				for _, transfer := range transfers {
					require.True(t, transfer.FromAccountID == fromAccount.ID || transfer.ToAccountID == fromAccount.ID)
				}
			},
		},
		{
			name: "No transfers (invalid account)",
			arg: ListTransfersParams{
				FromAccountID: fromAccount.ID + 1000, // ID inválido
				ToAccountID:   toAccount.ID + 1000,   // ID inválido
				Limit:         5,
				Offset:        0,
			},
			expected: 0,
			checkFn: func(t *testing.T, transfers []Transfer) {
				require.Empty(t, transfers)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			transfers, err := testQueries.ListTransfers(context.Background(), tc.arg)
			require.NoError(t, err)
			require.Len(t, transfers, tc.expected)

			if tc.expected > 0 {
				for _, transfer := range transfers {
					require.NotEmpty(t, transfer)
				}
				tc.checkFn(t, transfers)
			}
		})
	}
}

// TestTransferTxs testa casos específicos de transações
func TestTransferTxs(t *testing.T) {
	account1 := createRandomAccount(t).account
	account2 := createRandomAccount(t).account

	// Teste de criação de múltiplas transferências em sequência
	t.Run("sequential transfers", func(t *testing.T) {
		n := 5
		errs := make(chan error)
		results := make(chan Transfer)

		// Criar transferências em goroutines paralelas
		for i := 0; i < n; i++ {
			go func() {
				transfer, err := testQueries.CreateTransfer(context.Background(), CreateTransferParams{
					FromAccountID: account1.ID,
					ToAccountID:   account2.ID,
					Amount:        util.RandomPositiveMoney(),
				})

				errs <- err
				results <- transfer
			}()
		}

		// Verificar resultados
		for i := 0; i < n; i++ {
			err := <-errs
			require.NoError(t, err)

			transfer := <-results
			require.NotEmpty(t, transfer)
			require.Equal(t, account1.ID, transfer.FromAccountID)
			require.Equal(t, account2.ID, transfer.ToAccountID)
		}
	})
}

// Helper function to compare transfers
func compareTransfers(t *testing.T, params CreateTransferParams, transfer Transfer) {
	require.Equal(t, params.FromAccountID, transfer.FromAccountID)
	require.Equal(t, params.ToAccountID, transfer.ToAccountID)
	require.Equal(t, params.Amount, transfer.Amount)
	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)
}
