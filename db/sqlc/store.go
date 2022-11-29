package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store struct {
	*Queries
	db *sql.DB
}

// TransferTxParams contains all the input prams needed to create transaction !
type TransferTxParams struct {
	FromAccountId int64 `json:"from_account_id"`
	ToAccountId   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult is the result of transfer transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

// execTx execute a function within the datbase transaction
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err %v, rb err %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

// TransferTx performs a moneyTransaction from one account to the other
// it create a transfer record, an account entries and update accounts balance within a single databse transaction
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountId,
			ToAccountID:   arg.ToAccountId,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountId,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountId,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}
		// TODO: update account balance

		if arg.FromAccountId < arg.ToAccountId {
			result.FromAccount, result.ToAccount, err = AddMoney(ctx, q, arg.FromAccountId, -arg.Amount, arg.ToAccountId, arg.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = AddMoney(ctx, q, arg.ToAccountId, arg.Amount, arg.FromAccountId, -arg.Amount)
		}

		return nil
	})

	return result, err
}

func AddMoney(
	ctx context.Context,
	q *Queries,
	account1ID int64,
	amount1 int64,
	account2ID int64,
	amount2 int64,
) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     account1ID,
		Amount: amount1,
	})
	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     account2ID,
		Amount: amount2,
	})
	if err != nil {
		return
	}
	return
}
