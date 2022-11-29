package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTranferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)

	fmt.Println(">> Before : ", account1.Balance, account2.Balance)

	// num of threads
	n := 5
	amount := int64(10)

	// chanels to communicate errors and resutls between go sub routines and main go routine
	errs := make(chan error)
	reslts := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			ctx := context.Background()
			results, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountId: account1.ID,
				ToAccountId:   account2.ID,
				Amount:        amount,
			})
			errs <- err
			reslts <- results
		}()
	}

	// check resutls
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
		result := <-reslts
		require.NotEmpty(t, result)

		// check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check from entry
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, fromEntry.AccountID, account1.ID)
		require.Equal(t, fromEntry.Amount, -amount)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		// check to entry
		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, toEntry.AccountID, account2.ID)
		require.Equal(t, toEntry.Amount, amount)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		// check account
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, fromAccount.ID, account1.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, toAccount.ID, account2.ID)

		fmt.Println(">> Tx : ", fromAccount.Balance, toAccount.Balance)

		// check ballance
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
	}

	// check acccount updated balance
	updatedFromAccount, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedFromAccount)

	updatedToAccount, err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedToAccount)

	fmt.Println(">> after : ", updatedFromAccount.Balance, updatedToAccount.Balance)

	require.Equal(t, account1.Balance-int64(n)*amount, updatedFromAccount.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updatedToAccount.Balance)
}

func TestTrasferTxDeadlock(t *testing.T) {

	store := NewStore(testDB)

	account1 := CreateRandomAccount(t)
	account2 := CreateRandomAccount(t)

	n := 10
	amount := int64(10)

	errs := make(chan error)
	for i := 0; i < n; i++ {

		fromAccountID := account1.ID
		toAccountID := account2.ID

		if i%2 == 1 {
			fromAccountID = account2.ID
			toAccountID = account1.ID
		}

		go func() {
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountId: fromAccountID,
				ToAccountId:   toAccountID,
				Amount:        amount,
			})
			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
	}

	updatedAccount1, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount1)

	updatedAccount2, err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updatedAccount2)

	// check if balance stayed the same !
	require.Equal(t, account1.Balance, updatedAccount1.Balance)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)
}
