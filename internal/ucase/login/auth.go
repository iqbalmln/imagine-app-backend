package login

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"gitlab.privy.id/go_graphql/internal/appctx"
	"gitlab.privy.id/go_graphql/internal/entity"
	"gitlab.privy.id/go_graphql/internal/repositories"
	"gitlab.privy.id/go_graphql/pkg/logger"
)

type auth struct {
	loginRepository repositories.Auth
}

func NewAuth(loginRepository repositories.Auth) *auth {
	return &auth{
		loginRepository: loginRepository,
	}
}

func (a auth) Serve(data *appctx.Data) appctx.Response {
	var payload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(data.Request.Body).Decode(&payload)
	if err != nil {
		logger.Error(err)
		return *appctx.NewResponse().WithError(err.Error()).
			WithStatus("Error").
			WithCode(http.StatusInternalServerError).WithMessage("Register Error").
			WithEntity("registerError").WithState("registerError")
	}

	payreturn, err := a.loginRepository.Login(data.Request.Context(), payload.Email)
	if err != nil {
		fmt.Println("heree cuyy")
	}

	fmt.Println("=================PAYLOD")
	fmt.Println(payload.Password)
	fmt.Println("=================PASSWORD")
	fmt.Println(payreturn.Password)

	match := a.loginRepository.CheckPasswordHash(payreturn.Password, payload.Password)

	// if err := bcrypt.CompareHashAndPassword([]byte(payreturn.Role), []byte(payload.Password)); err != nil {
	// 	fmt.Println("errorr compare")
	// }

	if !match {
		fmt.Println("nengkene")

		return *appctx.NewResponse().
			WithStatus("Error").
			WithCode(http.StatusInternalServerError).WithMessage("Compare Password Error").
			WithEntity("comparePasswordError").WithState("comparePasswordError")

	}
	var JWT_SIGNING_METHOD = jwt.SigningMethodHS256
	var JWT_SIGNATURE_KEY = []byte("fgusbhinjklergoiernglkengkjerbngkerugerb8367yt8734597012840uohrgkdngkerbng")
	claims := entity.MyClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "YoiQbal",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 60)),
		},
		ID:       payreturn.ID,
		Username: payreturn.Username,
		// Role:     payreturn.Role,
	}

	token := jwt.NewWithClaims(
		JWT_SIGNING_METHOD,
		claims,
	)

	signedToken, err := token.SignedString(JWT_SIGNATURE_KEY)
	if err != nil {
		fmt.Println("error disini")
	}

	return *appctx.NewResponse().WithData(signedToken).
		WithStatus("SUCCESS").
		WithCode(http.StatusOK).WithMessage("Register Success").
		WithEntity("register").WithState("registerSuccess")
}
