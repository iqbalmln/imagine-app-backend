package login

import (
	"encoding/json"
	"net/http"
	"time"

	"gitlab.privy.id/go_graphql/internal/appctx"
	"gitlab.privy.id/go_graphql/internal/entity"
	"gitlab.privy.id/go_graphql/internal/repositories"
)

type register struct {
	loginRepository repositories.Auth
}

func NewRegister(loginRepository repositories.Auth) *register {
	return &register{
		loginRepository: loginRepository,
	}
}

func (r register) Serve(data *appctx.Data) appctx.Response {
	var payload entity.Login
	err := json.NewDecoder(data.Request.Body).Decode(&payload)
	if err != nil {
		return *appctx.NewResponse().WithError(err.Error()).
			WithStatus("Error").
			WithCode(http.StatusInternalServerError).WithMessage("Register Error").
			WithEntity("registerError").WithState("registerError")
	}

	hash, err := r.loginRepository.HashPassword(payload.Password)
	if err != nil {
		return *appctx.NewResponse().WithError(err.Error()).
			WithStatus("Error").
			WithCode(http.StatusInternalServerError).WithMessage("Hashing Error").
			WithEntity("errorHashing").WithState("errorHashing")
	}

	user := entity.Login{
		Username: payload.Username,
		Email:    payload.Email,
		Password: hash,
		// Role:      payload.Role,
		CreatedAt: time.Now(),
	}

	err = r.loginRepository.Register(data.Request.Context(), &user)

	if err != nil {
		return *appctx.NewResponse().WithError(err.Error()).
			WithStatus("Error").
			WithCode(http.StatusInternalServerError).WithMessage("Register Error").
			WithEntity("registerError").WithState("registerError")
	}
	return *appctx.NewResponse().WithData(user).
		WithStatus("SUCCESS").
		WithCode(http.StatusOK).WithMessage("Register Success").
		WithEntity("register").WithState("registerSuccess")
}
