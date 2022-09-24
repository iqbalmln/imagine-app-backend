// Package router
package router

import (
	"net/http"

	"gitlab.com/go_graphql/internal/appctx"
	"gitlab.com/go_graphql/internal/ucase/contract"
	"gitlab.com/go_graphql/pkg/routerkit"
)

// httpHandlerFunc is a contract http handler for router
type httpHandlerFunc func(request *http.Request, svc contract.UseCase, conf *appctx.Config) appctx.Response

// Router is a contract router and must implement this interface
type Router interface {
	Route() *routerkit.Router
}
