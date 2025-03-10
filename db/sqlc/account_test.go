package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/JMustang/OldBank/util"
	"github.com/stretchr/testify/require"
)

// TestAccount stores the test account data
type TestAccount struct {
	account Account
	params  CreateAccountParams
}

func createRandomAccount(t *testing.T) TestAccount {
	user := createRandomUser(t)

	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return TestAccount{
		account: account,
		params:  arg,
	}
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	testAccount := createRandomAccount(t)
	account2, err := testQueries.GetAccount(context.Background(), testAccount.account.ID)

	require.NoError(t, err)
	require.NotEmpty(t, account2)

	compareAccounts(t, testAccount.account, account2)
}

func TestUpdateAccount(t *testing.T) {
	testAccount := createRandomAccount(t)

	arg := UpdateAccountParams{
		ID:      testAccount.account.ID,
		Balance: util.RandomMoney(),
	}

	updatedAccount, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount)

	// Verificar se apenas o saldo foi alterado
	require.Equal(t, testAccount.account.ID, updatedAccount.ID)
	require.Equal(t, testAccount.account.Owner, updatedAccount.Owner)
	require.Equal(t, arg.Balance, updatedAccount.Balance)
	require.Equal(t, testAccount.account.Currency, updatedAccount.Currency)
	require.WithinDuration(t, testAccount.account.CreatedAt, updatedAccount.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	testAccount := createRandomAccount(t)

	// Testar a deleção
	err := testQueries.DeleteAccount(context.Background(), testAccount.account.ID)
	require.NoError(t, err)

	// Verificar se a conta foi realmente deletada
	account2, err := testQueries.GetAccount(context.Background(), testAccount.account.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account2)
}

func TestListAccounts(t *testing.T) {
	user := createRandomUser(t)

	var createdAccounts []Account
	targetOwner := user.Username

	// Lista de moedas únicas para garantir que não haja duplicatas
	currencies := []string{
		"USD", "EUR", "GBP", "JPY", "CAD",
		"AUD", "CHF", "CNY", "SEK", "NZD", // 10 moedas diferentes
	}

	// Criar 10 contas com o mesmo owner e moedas únicas
	for i := 0; i < 10; i++ {
		arg := CreateAccountParams{
			Owner:    targetOwner,
			Balance:  util.RandomMoney(),
			Currency: currencies[i], // Usa a moeda correspondente ao índice
		}

		account, err := testQueries.CreateAccount(context.Background(), arg)
		require.NoError(t, err)
		createdAccounts = append(createdAccounts, account)
	}

	// ... restante do teste permanece igual ...

	// Testar diferentes combinações de limit e offset
	testCases := []struct {
		name     string
		arg      ListAccountsParams
		expected int
	}{
		{
			name: "First 5 accounts",
			arg: ListAccountsParams{
				Owner:  targetOwner,
				Limit:  5,
				Offset: 0,
			},
			expected: 5,
		},
		{
			name: "Last 5 accounts",
			arg: ListAccountsParams{
				Owner:  targetOwner,
				Limit:  5,
				Offset: 5,
			},
			expected: 5,
		},
		{
			name: "All accounts",
			arg: ListAccountsParams{
				Owner:  targetOwner,
				Limit:  10,
				Offset: 0,
			},
			expected: 10,
		},
		{
			name: "No accounts (offset too high)",
			arg: ListAccountsParams{
				Owner:  targetOwner,
				Limit:  5,
				Offset: 15,
			},
			expected: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			accounts, err := testQueries.ListAccounts(context.Background(), tc.arg)
			require.NoError(t, err)
			require.Len(t, accounts, tc.expected)

			for _, account := range accounts {
				require.NotEmpty(t, account)
				require.Equal(t, targetOwner, account.Owner)
			}
		})
	}
}

// Helper function to compare accounts
func compareAccounts(t *testing.T, account1, account2 Account) {
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

// package db

// import (
// 	"context"
// 	"database/sql"
// 	"testing"
// 	"time"

// 	"github.com/JMustang/OldBank/util"
// 	"github.com/stretchr/testify/require"
// )

// func createRandomAccount(t *testing.T) Account {
// 	arg := CreateAccountParams{
// 		Owner:    util.RandomOwner(),
// 		Balance:  util.RandomMoney(),
// 		Currency: util.RandomCurrency(),
// 	}

// 	account, err := testQueries.CreateAccount(context.Background(), arg)
// 	require.NoError(t, err)
// 	require.NotEmpty(t, account)

// 	require.Equal(t, arg.Owner, account.Owner)
// 	require.Equal(t, arg.Balance, account.Balance)
// 	require.Equal(t, arg.Currency, account.Currency)

// 	require.NotZero(t, account.ID)
// 	require.NotZero(t, account.CreatedAt)

// 	return account
// }

// func TestCreateAccount(t *testing.T) {
// 	createRandomAccount(t)
// }

// func TestGetAccount(t *testing.T) {
// 	account1 := createRandomAccount(t)
// 	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
// 	require.NoError(t, err)
// 	require.NotEmpty(t, account2)

// 	require.Equal(t, account1.ID, account2.ID)
// 	require.Equal(t, account1.Owner, account2.Owner)
// 	require.Equal(t, account1.Balance, account2.Balance)
// 	require.Equal(t, account1.Currency, account2.Currency)
// 	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)

// }

// func TestUpdateAccount(t *testing.T) {
// 	account1 := createRandomAccount(t)

// 	arg := UpdateAccountParams{
// 		ID:      account1.ID,
// 		Balance: util.RandomMoney(),
// 	}

// 	account2, err := testQueries.UpdateAccount(context.Background(), arg)
// 	require.NoError(t, err)
// 	require.NotEmpty(t, account2)

// 	require.Equal(t, account1.ID, account2.ID)
// 	require.Equal(t, account1.Owner, account2.Owner)
// 	require.Equal(t, arg.Balance, account2.Balance)
// 	require.Equal(t, account1.Currency, account2.Currency)
// 	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
// }

// func TestDeleteAccount(t *testing.T) {
// 	account1 := createRandomAccount(t)
// 	err := testQueries.DeleteAccount(context.Background(), account1.ID)
// 	require.NoError(t, err)

// 	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
// 	require.Error(t, err)
// 	require.EqualError(t, err, sql.ErrNoRows.Error())
// 	require.Empty(t, account2)
// }

// func TestListAccounts(t *testing.T) {
// 	for i := 0; i < 10; i++ {
// 		createRandomAccount(t)
// 	}

// 	arg := ListAccountsParams{
// 		Limit:  5,
// 		Offset: 5,
// 	}

// 	accounts, err := testQueries.ListAccounts(context.Background(), arg)
// 	require.NoError(t, err)
// 	require.Len(t, accounts, 5)

//		for _, account := range accounts {
//			require.NotEmpty(t, account)
//		}
//	}
