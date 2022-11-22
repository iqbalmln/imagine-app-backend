package ucase

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gitlab.privy.id/go_graphql/internal/appctx"
	"gitlab.privy.id/go_graphql/internal/consts"
	"gitlab.privy.id/go_graphql/internal/repositories"
	"gitlab.privy.id/go_graphql/internal/ucase/contract"
	"gitlab.privy.id/go_graphql/pkg/logger"
)

type getShotID struct {
	shotRepositoryID repositories.Shot
}

func NewGetShotID(shotRepositoryID repositories.Shot) contract.UseCase {
	return &getShotID{shotRepositoryID: shotRepositoryID}
}

func (gs getShotID) Serve(data *appctx.Data) appctx.Response {
	param := mux.Vars(data.Request)
	rawID := param["id"]

	// Convert param to int
	id, err := strconv.ParseInt(rawID, 10, 64)

	if err != nil {
		logger.Error(logger.MessageFormat("[shot-getid]%w", err))
		return *appctx.NewResponse().WithCode(consts.CodeBadRequest)
	}

	shot, err := gs.shotRepositoryID.GetShotID(data.Request.Context(), int(id))

	if err != nil {
		logger.Error(logger.MessageFormat("shot-getid %w", err))

		return *appctx.NewResponse().WithCode(consts.CodeInternalServerError)
	}

	return *appctx.NewResponse().
		WithData(shot).
		WithCode(http.StatusOK).
		WithMessage("Get Shot Detail Success")
}
