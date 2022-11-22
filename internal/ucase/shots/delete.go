package ucase

import (
	"strconv"

	"github.com/gorilla/mux"
	"gitlab.privy.id/go_graphql/internal/appctx"
	"gitlab.privy.id/go_graphql/internal/consts"
	"gitlab.privy.id/go_graphql/internal/repositories"
	"gitlab.privy.id/go_graphql/internal/ucase/contract"
	"gitlab.privy.id/go_graphql/pkg/logger"
)

type deleteShot struct {
	shotRepository repositories.Shot
}

func NewShotDelete(shotRepository repositories.Shot) contract.UseCase {
	return &deleteShot{shotRepository: shotRepository}
}

func (u *deleteShot) Serve(data *appctx.Data) appctx.Response {
	params := mux.Vars(data.Request)
	rawID := params["id"]

	parseId, err := strconv.ParseInt(rawID, 10, 8)

	if err != nil {
		logger.Error(logger.MessageFormat("[shot delete] %v", err))

		return *appctx.NewResponse().WithCode(consts.CodeInternalServerError)
	}

	id := int(parseId)
	err = u.shotRepository.DeleteShot(data.Request.Context(), id)
	if err != nil {
		logger.Error(logger.MessageFormat("[shot deleted] %v", err))
		return *appctx.NewResponse().WithCode(consts.CodeInternalServerError).WithError(err.Error())
	}
	return *appctx.NewResponse().
		WithStatus("SUCCESS").
		WithEntity("deletedShot").
		WithState("deletedShotSuccess").
		WithMessage("Delete Shot Success")
}
