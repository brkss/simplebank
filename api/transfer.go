package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	db "github.com/brkss/simplebank/db/sqlc"
	"github.com/brkss/simplebank/token"
	"github.com/gin-gonic/gin"
)

type CreateTransferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=1"`
	Currency      string `json:"currency" binding:"required,oneof=USD EUR CAD"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var request CreateTransferRequest
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	fromAccount, valid := server.validAccount(ctx, request.FromAccountID, request.Currency)
	if !valid {
		return
	}
	payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if payload.Username != fromAccount.Owner {
		err = errors.New("account doesn't belong to the authicated user !")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	_, valid = server.validAccount(ctx, request.ToAccountID, request.Currency)
	if !valid {
		return
	}

	arg := db.TransferTxParams{
		FromAccountId: request.FromAccountID,
		ToAccountId:   request.ToAccountID,
		Amount:        request.Amount,
	}
	results, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, results)
}

func (server *Server) validAccount(ctx *gin.Context, accountID int64, currency string) (db.Account, bool) {

	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return account, false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return account, false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch %s / %s", accountID, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return account, false
	}

	return account, true
}
