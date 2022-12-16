package token

import (
	"testing"
	"time"

	"github.com/brkss/simplebank/utils"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/require"
)



func TestJWTMaked(t *testing.T){

	maker, err := NewJWTMaker(utils.RandomString(34));
	require.NoError(t, err)

	username := utils.RandomOwner()
	duration := time.Minute

	issueAt := time.Now()
	expireAt := time.Now().Add(duration)

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)

	require.NotZero(t, payload.ID)
	require.Equal(t, payload.Username, username)
	require.WithinDuration(t, payload.IssuedAt, issueAt, time.Second)
	require.WithinDuration(t, payload.ExpiredAt, expireAt, time.Second) 
}

func TestExpiredJWT(t *testing.T){
	maker, err := NewJWTMaker(utils.RandomString(34))
	require.NoError(t, err)

	username := utils.RandomOwner() 
	duration := time.Minute // token lifespan

	token, err := maker.CreateToken(username, -duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}

func TestInvalidJWTSignature(t *testing.T){
	payload, err := NewPayload(utils.RandomOwner(), time.Minute)
	require.NoError(t, err)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	token, err := jwtToken.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)

	maker, err := NewJWTMaker(utils.RandomString(32))
	require.NoError(t, err)
	
	payload, err = maker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Empty(t, payload)
}

func TestInvalidToken(t *testing.T){
	maker1, err := NewJWTMaker(utils.RandomString(34))
	maker2, err := NewJWTMaker(utils.RandomString(34))
	require.NoError(t, err)
	
	token, err := maker1.CreateToken(utils.RandomOwner(), time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker2.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)
}
