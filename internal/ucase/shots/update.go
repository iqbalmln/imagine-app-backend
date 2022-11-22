package ucase

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"gitlab.privy.id/go_graphql/internal/appctx"
	"gitlab.privy.id/go_graphql/internal/consts"
	"gitlab.privy.id/go_graphql/internal/entity"
	"gitlab.privy.id/go_graphql/internal/repositories"
	"gitlab.privy.id/go_graphql/internal/ucase/contract"
	"gitlab.privy.id/go_graphql/pkg/logger"
)

type updateShot struct {
	shotRepository repositories.Shot
}

func NewShotUpdate(shotRepository repositories.Shot) contract.UseCase {
	return &updateShot{shotRepository: shotRepository}
}

func (us updateShot) Serve(data *appctx.Data) appctx.Response {
	payload := entity.Shot{}
	err := data.Cast(&payload)
	if err != nil {
		fmt.Println(err)
		logger.Error(logger.MessageFormat("[Shot update] parsinng body request error: %v", err))
		return *appctx.NewResponse().WithCode(consts.CodeBadRequest).WithError(err.Error())
	}

	params := mux.Vars(data.Request)
	rawID := params["id"]

	parseId, err := strconv.ParseInt(rawID, 10, 8)

	if err != nil {
		fmt.Println(err)
		logger.Error(logger.MessageFormat("[Shot updated] %v", err))

		return *appctx.NewResponse().WithCode(consts.CodeInternalServerError)
	}
	id := int(parseId)
	shot, err := us.shotRepository.GetShotID(data.Request.Context(), id)
	if err != nil {
		fmt.Println(err)
		logger.Error(logger.MessageFormat("[Shot updated] %v", err))
		return *appctx.NewResponse().WithCode(consts.CodeInternalServerError).WithError(err.Error())
	}
	shot.Title = payload.Title
	shot.Description = payload.Description
	shot.Category = payload.Category
	shot.UpdatedAt = time.Now()

	err = us.shotRepository.UpdateShot(data.Request.Context(), shot)
	if err != nil {
		fmt.Println(err)
		logger.Error(logger.MessageFormat("[shot update] %v", err))
		return *appctx.NewResponse().WithCode(consts.CodeInternalServerError).WithError(err.Error())
	}
	return *appctx.NewResponse().WithCode(consts.CodeSuccess).WithData(shot)
}
