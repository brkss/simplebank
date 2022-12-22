package token

import (
	"testing"
	"time"

	"github.com/brkss/simplebank/utils"
	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T){
	pasetoMaker, err := NewPasetoMaker(utils.RandomString(32))
	require.NoError(t, err);

	username := utils.RandomOwner()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := time.Now().Add(duration)

	token, err := pasetoMaker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := pasetoMaker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, payload.Username, username)
	require.WithinDuration(t, payload.ExpiredAt, expiredAt, time.Second)
	require.WithinDuration(t, payload.IssuedAt, issuedAt, time.Second)
}


func TestExpiredPaseto(t *testing.T){
	maker, err := NewPasetoMaker(utils.RandomString(32))
	require.NoError(t, err)

	token, err := maker.CreateToken(utils.RandomOwner(), -time.Minute)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Empty(t, token)
}
