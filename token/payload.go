package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// expired token error
var ErrExpiredToken = errors.New("token has expired") 
var ErrInvalidToken = errors.New("invalid token !");

type Payload struct {
    ID          uuid.UUID	`json:"id"`
	Username	string		`json:"username"` 
	IssuedAt	time.Time	`json:"issued_at"`
	ExpiredAt	time.Time	`json:"expired_at"`
}

func NewPayload(username string, duration time.Duration) (*Payload, error) {
	tokenId, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID: tokenId,
		Username: username,
		IssuedAt: time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	return payload, nil
}


// Valid check if the token is valid or not ! 
func (p *Payload)Valid() error {
	if time.Now().After(p.ExpiredAt){
		return ErrExpiredToken;
	}
	return nil;
}
