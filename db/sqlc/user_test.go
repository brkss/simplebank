package db

import (
	"context"
  "time"
	"testing"

	utils "github.com/brkss/simplebank/utils"
	"github.com/stretchr/testify/require"
)


func createRandomUser(t *testing.T) User {

  arg := CreateUserParams{
    Username: utils.RandomOwner(),
    FullName: utils.RandomOwner(),
    HashedPassword: "secret",
    Email: utils.RandomEmail(),
  }

  user, err := testQueries.CreateUser(context.Background(), arg)
  require.NoError(t, err)
  require.NotEmpty(t, user)


  require.Equal(t, user.Username, arg.Username)
  require.Equal(t, user.FullName, arg.FullName)
  require.Equal(t, user.HashedPassword, arg.HashedPassword)
  require.Equal(t, user.Email, arg.Email)

  require.NotZero(t, user.CreatedAt)
  require.True(t, user.PasswordChanged.IsZero())

  return user
}

func TestCreateUser(t *testing.T){
  createRandomUser(t)
}

func TestGetUser(t *testing.T){

  user1 := createRandomUser(t);

  user2, err := testQueries.GetUser(context.Background(), user1.Username)
  require.NoError(t, err)
  require.NotEmpty(t, user2)

  require.Equal(t, user1.Username, user2.Username);
  require.Equal(t, user1.FullName, user2.FullName);
  require.Equal(t, user1.HashedPassword, user2.HashedPassword);
  require.Equal(t, user1.Email, user2.Email);
  require.Equal(t, user1.PasswordChanged, user2.PasswordChanged);
  require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second) 
}

