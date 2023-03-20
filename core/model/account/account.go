package account

import (
	"context"
	"soporte-go/core/model/user"
	"time"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthenticationResponse struct {
	Token string        `json:"access_token"`
	User  user.UserAuth `json:"user"`
}

type RegisterAuthResponse struct {
	Access string               `json:"access_token"`
	User   user.ClienteResponse `json:"user"`
}

type RegisterForm struct {
	Username   *string `json:"username,omitempty"`
	Password   *string `json:"password"`
	EmpresaId  int     `json:"empresa_id,omitempty"`
	SuperiorId string  `json:"superior_id"`
	Email      *string `json:"email"`
	// IsAdmin    bool    `json:"is_admiin"`
	Rol int `json:"rol,omitempty"`
}

type User struct {
	Username  *string    `json:"username,omitempty"`
	Password  *string    `json:"password"`
	UserId    string     `json:"user_id,omitempty"`
	LastLogin *time.Time `json:"last_Login,omitempty"`
	CreatedOn *time.Time `json:"created_on,omitempty"`
	Email     *string    `json:"email"`
	Estado    int        `json:"estado"`
	EmpresaId int        `json:"empresa_id"`
	Rol       *int       `json:"rol,omitempty"`
}

type AccountUseCase interface {
	// Fetch(ctx context.Context,  num int64) ([]User, string, error)
	Login(ctx context.Context, loginRequest *LoginRequest) (user.UserAuth, error)
	// Update(ctx context.Context, ar *User) error
	RegisterCliente(context.Context, *RegisterForm) (user.ClienteResponse, error)
	RegisterFuncionario(context.Context, *RegisterForm) (user.UserAuth, error)
	DeleteUser(context.Context, string) (err error)
}

// ArticleRepository represent the article's repository contract
type AccountRepository interface {
	// Fetch(ctx context.Context, num int64) (res []User, err error)
	Login(ctx context.Context, loginRequest *LoginRequest) (res user.UserAuth, err error)
	// Update(ctx context.Context, ar *User) error
	RegisterCliente(context.Context, *RegisterForm) (user.ClienteResponse, error)
	RegisterFuncionario(context.Context, *RegisterForm) (user.UserAuth, error)
	DeleteUser(context.Context, string) (err error)
	// ValidateInvitation(ctx context.Context,mail *string,rol *int)(error)
}
