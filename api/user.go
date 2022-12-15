package api

import (
	"net/http"

  "time"

  pq "github.com/lib/pq"
	db "github.com/brkss/simplebank/db/sqlc"
	"github.com/brkss/simplebank/utils"
	"github.com/gin-gonic/gin"
)

type CreateUserRequest struct {
  Username  string `json:"username" binding:"required"`
  Email     string `json:"email" binding:"required"`
  FullName  string `json:"full_name" binding:"required"`
  Password  string `json:"password" binding:"required"`
}

type CreateUserResponse struct {
  Username          string    `json:"username"`
  FullName          string    `json:"full_name"`
  Email             string    `json:"email"`
  PasswordChanged   time.Time `json:"password_changed_at"`
  CreatedAt         time.Time `json:""created_at`
}

func (server *Server)createUser(ctx *gin.Context){
 
  var request CreateUserRequest
  err := ctx.ShouldBindJSON(&request)

  if err != nil {
    ctx.JSON(http.StatusBadRequest, errorResponse(err))
    return
  }

  hashedPassword, err := utils.HashPassword(request.Password)
  if err != nil {
    ctx.JSON(http.StatusInternalServerError, errorResponse(err))
  }

  arg := db.CreateUserParams{
    Username: request.Username,
    FullName: request.FullName,
    Email: request.Email,
    HashedPassword: hashedPassword,
  }


  user, err := server.store.CreateUser(ctx, arg)
  if err != nil {
    pqErr, ok := err.(*pq.Error)
    if ok {
      switch pqErr.Code.Name(){
        case "unique_violation":
          ctx.JSON(http.StatusForbidden, errorResponse(err))
          return
      }
    }
    ctx.JSON(http.StatusInternalServerError, errorResponse(err))
    return
  }
  
  resp := CreateUserResponse{
    Username: user.Username,
    Email: user.Email,
    FullName: user.FullName,
    CreatedAt: user.CreatedAt,
    PasswordChanged: user.PasswordChanged,

  }

  ctx.JSON(http.StatusOK, resp)
  return
}
