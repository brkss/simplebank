package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/brkss/simplebank/db/mock"
	db "github.com/brkss/simplebank/db/sqlc"
	"github.com/brkss/simplebank/utils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetAccount(t *testing.T) {
	account := randomAccount()
	testCases := []struct {
		name          string
		accountID     int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name:      "NotFound",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "InternalError",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "BadRequest",
			accountID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			// build stbus
			tc.buildStubs(store)

			// start test server and send requests
			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/account/%d", tc.accountID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}

}

func randomAccount() db.Account {
	return db.Account{
		ID:       utils.RandomInt(0, 10000),
		Owner:    utils.RandomOwner(),
		Currency: utils.RandomCurrency(),
		Balance:  utils.RandomMoney(),
	}
}

// requireMatchAccount check the account we got from the request
func requireMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {

	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)

	require.Equal(t, account, gotAccount)
}

func TestListAccounts(t *testing.T) {
	accounts := createRandomAccounts(10)
	testCases := []struct {
		name          string
		accounts      []db.Account
		limit         int64
		offset        int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:     "OK",
			accounts: accounts,
			limit:    5,
			offset:   5,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Eq(db.ListAccountsParams{Limit: 5, Offset: 5})).
					Times(1).
					Return(accounts[0:5], nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				checkListAccountResponse(t, recorder.Body, accounts[0:5])
			},
		},
		{
			name:     "BadRequest",
			accounts: accounts,
			limit:    -1,
			offset:   -1,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
				//checkListAccountResponse(t, recorder.Body, accounts[0:5])
			},
		},

		{
			name:     "InternalServerError",
			accounts: accounts,
			limit:    5,
			offset:   5,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Eq(db.ListAccountsParams{Limit: 5, Offset: 5})).
					Times(1).
					Return([]db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
				//checkListAccountResponse(t, recorder.Body, accounts[0:5])
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			// build stubs
			tc.buildStubs(store)

			recorder := httptest.NewRecorder()
			server := NewServer(store)

			url := fmt.Sprintf("/accounts/%d/%d", tc.limit, tc.offset)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)

			// check response
		})
	}
}

func createRandomAccounts(n int) []db.Account {

	var accounts []db.Account

	for i := 0; i < n; i++ {
		accounts = append(accounts, db.Account{
			ID:       utils.RandomInt(0, 1000),
			Owner:    utils.RandomOwner(),
			Currency: utils.RandomCurrency(),
			Balance:  utils.RandomMoney(),
		})
	}

	return accounts
}

func checkListAccountResponse(t *testing.T, body *bytes.Buffer, accounts []db.Account) {

	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotAccounts []db.Account
	err = json.Unmarshal(data, &gotAccounts)
	require.NoError(t, err)

	for i := 0; i < len(accounts); i++ {
		require.Equal(t, accounts[i], gotAccounts[i])
	}
}

/*
func TestCreateAccount(t *testing.T) {

		arg := CreateAccountRequest{
			Owner:    utils.RandomOwner(),
			Currency: utils.RandomCurrency(),
		}
		testCases := []struct {
			name          string
			arg           CreateAccountRequest
			buildStabs    func(store *mockdb.MockStore)
			checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
		}{
			{
				name: "OK",
				arg:  arg,
				buildStabs: func(store *mockdb.MockStore) {
					store.EXPECT().
						CreateAccount(gomock.Any(), gomock.Eq(arg)).
						Times(1).
						Return(db.Account{
							ID:       utils.RandomInt(0, 1000),
							Owner:    arg.Owner,
							Currency: arg.Currency,
							Balance:  0,
						}, nil)

				},
				checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
					require.Equal(t, recorder.Code, http.StatusOK)
					checkCreateAccountResponse(t, recorder.Body, arg)
				},
			},
		}


			for i := range testCases {
				tc := testCases[i]
				t.Run(tc.name, func(t *testing.T) {

					ctrl := gomock.NewController(t)

					store := mockdb.NewMockStore(ctrl)
					tc.buildStabs(store)

					recorder := httptest.NewRecorder()
					server := NewServer(store)

					url := "/account"
					body := bytes.NewReader(arg)
					request, err := http.NewRequest(http.MethodPost, url, body)
					require.NoError(t, err)
					server.router.ServeHTTP(recorder, request)
				})
			}
	}
*/
func checkCreateAccountResponse(t *testing.T, body *bytes.Buffer, arg CreateAccountRequest) {

	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var account db.Account
	err = json.Unmarshal(data, &account)
	require.NoError(t, err)

	require.Equal(t, account.Balance, 0)
	require.Equal(t, account.Owner, arg.Owner)
	require.Equal(t, account.Currency, arg.Currency)

}
