package db

import (
	"context"
	"testing"
	"time"

	"github.com/brkss/simplebank/utils"
	"github.com/stretchr/testify/require"
)

func CreateRandomAccount(t *testing.T, user User) Account {

	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  utils.RandomMoney(),
		Currency: utils.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
	return account
}

func TestCreateAccount(t *testing.T) {
	user := createRandomUser(t)
	CreateRandomAccount(t, user)
}

func TestGetAccount(t *testing.T) {
	user := createRandomUser(t)
	account_data := CreateRandomAccount(t, user)
	account, err := testQueries.GetAccount(context.Background(), account_data.ID)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, account_data.ID, account.ID)
	require.Equal(t, account_data.Owner, account.Owner)
	require.Equal(t, account_data.Balance, account.Balance)
	require.Equal(t, account_data.Currency, account.Currency)
	require.WithinDuration(t, account_data.CreatedAt, account.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	user := createRandomUser(t)
	account1 := CreateRandomAccount(t, user)
	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: utils.RandomMoney(),
	}
	account2, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account2.ID, account1.ID)
	require.Equal(t, account2.Owner, account1.Owner)
	require.Equal(t, arg.Balance, account2.Balance)
	require.Equal(t, account2.Currency, account1.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestListAccounts(t *testing.T) {
	user := createRandomUser(t)
	for i := 0; i < 10; i++ {
		CreateRandomAccount(t, user)
	}
	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
		Owner:  user.Username,
	}
	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, acc := range accounts {
		require.NotEmpty(t, acc)
	}

}
