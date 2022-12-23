package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/brkss/simplebank/token"
	"github.com/brkss/simplebank/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestAuthMiddleware(t *testing.T) {

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				token, err := tokenMaker.CreateToken(utils.RandomOwner(), time.Minute)
				require.NoError(t, err)
				authorization := fmt.Sprintf("%s %s", authorizationTypeBearer, token)
				request.Header.Set(authorizationHeaderKey, authorization)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusOK)
			},
		},
		{
			name: "ExpiredToken",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				token, err := tokenMaker.CreateToken(utils.RandomOwner(), -time.Minute)
				require.NoError(t, err)
				authorization := fmt.Sprintf("%s %s", authorizationTypeBearer, token)
				request.Header.Set(authorizationHeaderKey, authorization)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusUnauthorized)
			},
		},
		{
			name: "UnsetToken",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusUnauthorized)
			},
		},
		{
			name: "InvalidToken",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				authorization := fmt.Sprintf("%s %s", authorizationTypeBearer, "sdaskdasdsalkd")
				request.Header.Set(authorizationHeaderKey, authorization)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusUnauthorized)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := NewTestServer(t, nil)
			server.router.GET("/auth", authMiddleware(server.tokenMaker), func(ctx *gin.Context) {
				ctx.JSON(http.StatusOK, gin.H{})
			})

			url := "/auth"
			recorder := httptest.NewRecorder()
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)
			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}

}
