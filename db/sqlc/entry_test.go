package db

import (
	"context"
	"testing"
	"time"

	"github.com/JMustang/OldBank/util"
	"github.com/stretchr/testify/require"
)

// TestEntry stores the test entry data and its associated account
type TestEntry struct {
	entry   Entry
	account Account
	params  CreateEntryParams
}

func createRandomEntry(t *testing.T, account Account) TestEntry {
	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomMoney(), // Pode ser positivo ou negativo
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, arg.AccountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)
	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

	return TestEntry{
		entry:   entry,
		account: account,
		params:  arg,
	}
}

func TestCreateEntry(t *testing.T) {
	// Primeiro criamos uma conta para associar com a entry
	account := createRandomAccount(t).account

	// Testamos criação com valor positivo
	t.Run("positive amount", func(t *testing.T) {
		arg := CreateEntryParams{
			AccountID: account.ID,
			Amount:    util.RandomPositiveMoney(),
		}

		entry, err := testQueries.CreateEntry(context.Background(), arg)
		require.NoError(t, err)
		require.NotEmpty(t, entry)
		require.Equal(t, arg.AccountID, entry.AccountID)
		require.Equal(t, arg.Amount, entry.Amount)
		require.NotZero(t, entry.ID)
		require.NotZero(t, entry.CreatedAt)
	})

	// Testamos criação com valor negativo
	t.Run("negative amount", func(t *testing.T) {
		arg := CreateEntryParams{
			AccountID: account.ID,
			Amount:    util.RandomNegativeMoney(),
		}

		entry, err := testQueries.CreateEntry(context.Background(), arg)
		require.NoError(t, err)
		require.NotEmpty(t, entry)
		require.Equal(t, arg.AccountID, entry.AccountID)
		require.Equal(t, arg.Amount, entry.Amount)
		require.NotZero(t, entry.ID)
		require.NotZero(t, entry.CreatedAt)
	})
}

func TestGetEntry(t *testing.T) {
	account := createRandomAccount(t).account
	entry1 := createRandomEntry(t, account)

	entry2, err := testQueries.GetEntry(context.Background(), entry1.entry.ID)
	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	compareEntries(t, entry1.entry, entry2)
}

func TestListEntries(t *testing.T) {
	account := createRandomAccount(t).account

	// Criar 10 entries para a mesma conta
	var createdEntries []Entry
	for i := 0; i < 10; i++ {
		entry := createRandomEntry(t, account)
		createdEntries = append(createdEntries, entry.entry)
	}

	// Testar diferentes combinações de limit e offset
	testCases := []struct {
		name     string
		arg      ListEntriesParams
		expected int
	}{
		{
			name: "First 5 entries",
			arg: ListEntriesParams{
				AccountID: account.ID,
				Limit:     5,
				Offset:    0,
			},
			expected: 5,
		},
		{
			name: "Last 5 entries",
			arg: ListEntriesParams{
				AccountID: account.ID,
				Limit:     5,
				Offset:    5,
			},
			expected: 5,
		},
		{
			name: "All entries",
			arg: ListEntriesParams{
				AccountID: account.ID,
				Limit:     10,
				Offset:    0,
			},
			expected: 10,
		},
		{
			name: "No entries (offset too high)",
			arg: ListEntriesParams{
				AccountID: account.ID,
				Limit:     5,
				Offset:    15,
			},
			expected: 0,
		},
		{
			name: "Invalid account ID",
			arg: ListEntriesParams{
				AccountID: account.ID + 1000, // ID que não existe
				Limit:     5,
				Offset:    0,
			},
			expected: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			entries, err := testQueries.ListEntries(context.Background(), tc.arg)
			require.NoError(t, err)
			require.Len(t, entries, tc.expected)

			if tc.expected > 0 {
				for _, entry := range entries {
					require.NotEmpty(t, entry)
					require.Equal(t, account.ID, entry.AccountID)
				}
			}
		})
	}
}

// Helper function to compare entries
func compareEntries(t *testing.T, entry1, entry2 Entry) {
	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, time.Second)
}

// TestEntryTxs testa casos específicos de transações
func TestEntryTxs(t *testing.T) {
	account := createRandomAccount(t).account

	// Teste de criação de múltiplas entries em sequência
	t.Run("sequential entries", func(t *testing.T) {
		n := 5
		errs := make(chan error)
		results := make(chan Entry)

		// Criar entries em goroutines paralelas
		for i := 0; i < n; i++ {
			go func() {
				entry, err := testQueries.CreateEntry(context.Background(), CreateEntryParams{
					AccountID: account.ID,
					Amount:    util.RandomMoney(),
				})

				errs <- err
				results <- entry
			}()
		}

		// Verificar resultados
		for i := 0; i < n; i++ {
			err := <-errs
			require.NoError(t, err)

			entry := <-results
			require.NotEmpty(t, entry)
			require.Equal(t, account.ID, entry.AccountID)
		}
	})
}

// package db

// import (
// 	"context"
// 	"testing"
// 	"time"

// 	"github.com/JMustang/OldBank/util"
// 	"github.com/stretchr/testify/require"
// )

// func createRandomEntry(t *testing.T, account Account) Entry {
// 	arg := CreateEntryParams{
// 		AccountID: account.ID,
// 		Amount:    util.RandomMoney(),
// 	}

// 	entry, err := testQueries.CreateEntry(context.Background(), arg)
// 	require.NoError(t, err)
// 	require.NotEmpty(t, entry)

// 	require.Equal(t, arg.AccountID, entry.AccountID)
// 	require.Equal(t, arg.Amount, entry.Amount)

// 	require.NotZero(t, entry.ID)
// 	require.NotZero(t, entry.CreatedAt)

// 	return entry
// }

// func TestCreateEntry(t *testing.T) {
// 	account := createRandomAccount(t)
// 	createRandomEntry(t, account)
// }

// func TestGetEntry(t *testing.T) {
// 	account := createRandomAccount(t)
// 	entry1 := createRandomEntry(t, account)
// 	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)
// 	require.NoError(t, err)
// 	require.NotEmpty(t, entry2)

// 	require.Equal(t, entry1.ID, entry2.ID)
// 	require.Equal(t, entry1.AccountID, entry2.AccountID)
// 	require.Equal(t, entry1.Amount, entry2.Amount)
// 	require.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, time.Second)
// }

// func TestListEntries(t *testing.T) {
// 	account := createRandomAccount(t)
// 	for i := 0; i < 10; i++ {
// 		createRandomEntry(t, account)
// 	}

// 	arg := ListEntriesParams{
// 		AccountID: account.ID,
// 		Limit:     5,
// 		Offset:    5,
// 	}

// 	entries, err := testQueries.ListEntries(context.Background(), arg)
// 	require.NoError(t, err)
// 	require.Len(t, entries, 5)

// 	for _, entry := range entries {
// 		require.NotEmpty(t, entry)
// 		require.Equal(t, arg.AccountID, entry.AccountID)
// 	}
// }
