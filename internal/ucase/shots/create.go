package ucase

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"gitlab.privy.id/go_graphql/internal/appctx"
	"gitlab.privy.id/go_graphql/internal/entity"
	"gitlab.privy.id/go_graphql/internal/repositories"
)

type createShot struct {
	shotRepository repositories.Shot
}

func (cs createShot) Serve(data *appctx.Data) appctx.Response {
	var input entity.Shot

	err := json.NewDecoder(data.Request.Body).Decode(&input)
	if err != nil {
		err = fmt.Errorf("decoding request body : %w", err)

		return *appctx.NewResponse().
			WithCode(http.StatusInternalServerError).
			WithError(err)
	}

	now := time.Now()
	shot := entity.Shot{
		Title:       input.Title,
		IMG:         input.IMG,
		Description: input.Description,
		Category:    input.Category,
		CreatedAt:   now,
		UpdatedAt:   now,
		DeletedAt:   input.DeletedAt,
	}

	output, err := cs.shotRepository.CreateShot(data.Request.Context(), &shot)
	if err != nil {
		err = fmt.Errorf("upload shot: %w", err)
		fmt.Println(err)
		return *appctx.NewResponse().
			WithCode(http.StatusInternalServerError).
			WithError(err)

	}
	return *appctx.NewResponse().
		WithCode(http.StatusCreated).
		WithEntity("Uploaded Shot").
		WithState("Uploaded Shot Success").
		WithData(output)
}

func NewCreateShot(shotRepository repositories.Shot) *createShot {
	return &createShot{
		shotRepository: shotRepository,
	}
}
