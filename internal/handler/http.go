// Package handler
package handler

import (
	"context"
	"net/http"

	"gitlab.privy.id/go_graphql/internal/appctx"
	"gitlab.privy.id/go_graphql/internal/consts"
	"gitlab.privy.id/go_graphql/internal/ucase/contract"
)

// HttpRequest handler func wrapper
func HttpRequest(request *http.Request, svc contract.UseCase, conf *appctx.Config) appctx.Response {
	ctx := context.WithValue(request.Context(), consts.CtxLang, request.Header.Get(consts.HeaderLanguageKey))

	req := request.WithContext(ctx)

	data := &appctx.Data{
		Request:     req,
		Config:      conf,
		ServiceType: consts.ServiceTypeHTTP,
	}

	return svc.Serve(data)
}
