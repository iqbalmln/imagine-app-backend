package login

import (
	"fmt"
	"net/http"

	"gitlab.privy.id/go_graphql/internal/appctx"
	"gitlab.privy.id/go_graphql/internal/repositories"
)

type getUsers struct {
	loginRepository repositories.Auth
}

func (gu getUsers) Serve(data *appctx.Data) appctx.Response {
	users, err := gu.loginRepository.Get(data.Request.Context())
	if err != nil {
		fmt.Println("getting users")
		return *appctx.NewResponse().WithCode(http.StatusInternalServerError).WithError(err)

	}
	return *appctx.NewResponse().WithCode(http.StatusOK).WithData(users)

}

func NewGetUsers(loginRepository repositories.Auth) *getUsers {
	return &getUsers{
		loginRepository: loginRepository,
	}
}
