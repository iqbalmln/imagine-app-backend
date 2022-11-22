package ucase

import (
	"net/http"

	"gitlab.privy.id/go_graphql/internal/appctx"
	"gitlab.privy.id/go_graphql/internal/repositories"
)

type getShot struct {
	shotRepository repositories.Shot
}

func (gs getShot) Serve(data *appctx.Data) appctx.Response {
	shots, err := gs.shotRepository.Get(data.Request.Context())
	if err != nil {
		// Do Error Handling
	}
	// fmt.Println(shots)
	return *appctx.NewResponse().
		WithCode(http.StatusOK).
		WithData(shots)
}

func NewGetShot(shotRepository repositories.Shot) *getShot {
	return &getShot{
		shotRepository: shotRepository,
	}
}
