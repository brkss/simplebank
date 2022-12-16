package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const MIN_SECRET_KEY_SIZE = 32

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) (Maker, error) {

	if len(secretKey) < MIN_SECRET_KEY_SIZE {
		return nil, fmt.Errorf("invalid secret: must have atleast %d characters !", MIN_SECRET_KEY_SIZE);
	}
	
	return &JWTMaker{secretKey: secretKey}, nil;
}

func (j *JWTMaker)CreateToken(username string, duration time.Duration)(string, error){

	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload) 
	token, err := jwtToken.SignedString([]byte(j.secretKey));
	
	if err != nil {
		return "", err
	}
	return token, nil
}

func (j *JWTMaker)VerifyToken(token string)(*Payload, error){
	keyFunc := func(token *jwt.Token)(interface{}, error){
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(j.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken){
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}
	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}
