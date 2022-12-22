package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	mockdb "github.com/brkss/simplebank/db/mock"
	db "github.com/brkss/simplebank/db/sqlc"
	"github.com/brkss/simplebank/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

type eqCreateUserMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (e eqCreateUserMatcher) Matches(x interface{}) bool {

	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	err := utils.VerifyPassword(e.password, arg.HashedPassword)
	if err != nil {
		return false
	}

	e.arg.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(e.arg, x)
}

func (e eqCreateUserMatcher) String() string {
	return fmt.Sprintf("matches arg %v, password arg %v", e.arg, e.password)
}

func EqCreateUserParam(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserMatcher{arg, password}
}

func TestCreateUserAPI(t *testing.T) {
	user, password := randomUser(t)

	testCases := []struct {
		name          string
		body          gin.H
		buildStabs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"username":  user.Username,
				"email":     user.Email,
				"password":  password,
				"full_name": user.FullName,
			},
			buildStabs: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Username: user.Username,
					FullName: user.FullName,
					Email:    user.Email,
				}
				store.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParam(arg, password)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatch(t, recorder.Body, user)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStabs(store)

			server := NewTestServer(t, store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/users"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}

}

func randomUser(t *testing.T) (db.User, string) {

	password := utils.RandomString(10)
	hashedPassword, err := utils.HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)
	user := db.User{
		Username:        utils.RandomString(6),
		HashedPassword:  hashedPassword,
		FullName:        utils.RandomOwner(),
		Email:           fmt.Sprintf("%s@test.com", utils.RandomString(10)),
		PasswordChanged: time.Now(),
		CreatedAt:       time.Now(),
	}

	return user, password
}

func requireBodyMatch(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var user_res db.User
	err = json.Unmarshal(data, &user_res)
	require.NoError(t, err)

	fmt.Printf("email: %s, username: %s\n", user.Email, user.Username)

	require.Equal(t, user.FullName, user_res.FullName)
	require.Equal(t, user.Username, user_res.Username)
	require.Equal(t, user.Email, user_res.Email)
}
