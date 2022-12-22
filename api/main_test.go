package api

import (
	"os"
	"testing"
	"time"

	db "github.com/brkss/simplebank/db/sqlc"
	"github.com/brkss/simplebank/utils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func NewTestServer(t *testing.T, store db.Store) *Server {
	config := utils.Config{
		TokenSymetricKey: utils.RandomString(32),
		TokenDuration:    time.Minute,
	}
	server, err := NewServer(config, store)
	require.NoError(t, err)

	return server
}

func TestMain(m *testing.M) {

	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())

}
