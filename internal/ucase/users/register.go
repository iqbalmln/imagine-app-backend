package users

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"gitlab.privy.id/go_graphql/internal/appctx"
	"gitlab.privy.id/go_graphql/internal/entity"
	"gitlab.privy.id/go_graphql/internal/repositories"
)

type registerUser struct {
	userRepository repositories.User
}

func (ru registerUser) Serve(data *appctx.Data) appctx.Response {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(data.Request.Body).Decode(&input)
	if err != nil {
		err = fmt.Errorf("decoding request body : %w", err)
		fmt.Println(err)
		return *appctx.NewResponse().
			WithCode(http.StatusInternalServerError).
			WithError(err)
	}
	now := time.Now()
	user := entity.User{
		Name:      input.Name,
		Email:     input.Email,
		Password:  input.Password,
		CreatedAt: now,
		DeletedAt: now,
	}
	err = ru.userRepository.RegisterUser(data.Request.Context(), &user)
	fmt.Println(user)
	if err != nil {
		err = fmt.Errorf("register user: %w", err)
		fmt.Println(err)
		return *appctx.NewResponse().
			WithCode(http.StatusInternalServerError).
			WithError(err)
	}
	return *appctx.NewResponse().
		WithCode(http.StatusCreated).
		WithEntity("Registered User").
		WithState("Register User Success").
		WithData(user)
}

func NewRegisterUser(userRepository repositories.User) *registerUser {
	return &registerUser{
		userRepository: userRepository,
	}
}
