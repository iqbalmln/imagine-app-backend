// Package router
package router

import (
	"context"
	"encoding/json"
	"net/http"
	"runtime/debug"

	"gitlab.privy.id/go_graphql/internal/appctx"
	"gitlab.privy.id/go_graphql/internal/bootstrap"
	"gitlab.privy.id/go_graphql/internal/consts"
	"gitlab.privy.id/go_graphql/internal/handler"
	"gitlab.privy.id/go_graphql/internal/middleware"
	"gitlab.privy.id/go_graphql/internal/repositories"
	"gitlab.privy.id/go_graphql/internal/ucase/login"
	ucase "gitlab.privy.id/go_graphql/internal/ucase/shots"
	"gitlab.privy.id/go_graphql/internal/ucase/users"
	"gitlab.privy.id/go_graphql/pkg/logger"
	"gitlab.privy.id/go_graphql/pkg/msgx"
	"gitlab.privy.id/go_graphql/pkg/routerkit"

	//"gitlab.privy.id/go_graphql/pkg/mariadb"
	//"gitlab.privy.id/go_graphql/internal/repositories"
	//"gitlab.privy.id/go_graphql/internal/ucase/example"

	ucaseContract "gitlab.privy.id/go_graphql/internal/ucase/contract"
)

type router struct {
	config *appctx.Config
	router *routerkit.Router
}

// NewRouter initialize new router wil return Router Interface
func NewRouter(cfg *appctx.Config) Router {
	bootstrap.RegistryMessage()
	bootstrap.RegistryLogger(cfg)

	return &router{
		config: cfg,
		router: routerkit.NewRouter(routerkit.WithServiceName(cfg.App.AppName)),
	}
}

func (rtr *router) handle(hfn httpHandlerFunc, svc ucaseContract.UseCase, mdws ...middleware.MiddlewareFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lang := r.Header.Get(consts.HeaderLanguageKey)
		if !msgx.HaveLang(consts.RespOK, lang) {
			lang = rtr.config.App.DefaultLang
			r.Header.Set(consts.HeaderLanguageKey, lang)
		}

		defer func() {
			err := recover()
			if err != nil {
				w.Header().Set(consts.HeaderContentTypeKey, consts.HeaderContentTypeJSON)
				w.WriteHeader(http.StatusInternalServerError)
				res := appctx.Response{
					Code: consts.CodeInternalServerError,
				}

				res.WithLang(lang)
				logger.Error(logger.MessageFormat("error %v", string(debug.Stack())))
				json.NewEncoder(w).Encode(res.Byte())

				return
			}
		}()

		ctx := context.WithValue(r.Context(), "access", map[string]interface{}{
			"path":      r.URL.Path,
			"remote_ip": r.RemoteAddr,
			"method":    r.Method,
		})

		req := r.WithContext(ctx)

		// validate middleware
		if !middleware.FilterFunc(w, req, rtr.config, mdws) {
			return
		}

		resp := hfn(req, svc, rtr.config)
		resp.WithLang(lang)
		rtr.response(w, resp)
	}
}

// response prints as a json and formatted string for DGP legacy
func (rtr *router) response(w http.ResponseWriter, resp appctx.Response) {
	w.Header().Set(consts.HeaderContentTypeKey, consts.HeaderContentTypeJSON)
	resp.Generate()
	w.WriteHeader(resp.Code)
	w.Write(resp.Byte())
	return
}

// postgresql
const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "intern-privy-simple-orders"
)

// Route preparing http router and will return mux router object
func (rtr *router) Route() *routerkit.Router {

	root := rtr.router.PathPrefix("/").Subrouter()
	inV1 := root.PathPrefix("/api/v1").Subrouter()
	//in := root.PathPrefix("/in/").Subrouter()
	// liveness := root.PathPrefix("/").Subrouter()
	//inV1 := in.PathPrefix("/v1/").Subrouter()

	// open tracer setup
	bootstrap.RegistryOpenTracing(rtr.config)

	//db := bootstrap.RegistryMariaMasterSlave(rtr.config.WriteDB, rtr.config.ReadDB, rtr.config.App.Timezone))
	//db := bootstrap.RegistryMariaDB(rtr.config.WriteDB, rtr.config.App.Timezone)
	db := bootstrap.RegistryPostgreSQLMasterSlave(rtr.config.ReadDB, rtr.config.WriteDB, rtr.config.App.Timezone)

	shotRepository := repositories.NewShotImplementation(db)
	userRepository := repositories.NewUserImplementation(db)
	loginRepository := repositories.NewLoginImplment(db)

	// Shot
	getShot := ucase.NewGetShot(shotRepository)
	getShotID := ucase.NewGetShotID(shotRepository)
	createShot := ucase.NewCreateShot(shotRepository)
	deleteShot := ucase.NewShotDelete(shotRepository)
	updateShot := ucase.NewShotUpdate(shotRepository)
	// root.Use(middleware.JwtAuthorization)

	//AUTH
	auth := rtr.router.PathPrefix("/auth").Subrouter()

	register := login.NewRegister(loginRepository)
	auth.HandleFunc("/register", rtr.handle(
		handler.HttpRequest,
		register,
	)).Methods(http.MethodPost)
	// User

	login := login.NewAuth(loginRepository)
	auth.HandleFunc("/login", rtr.handle(
		handler.HttpRequest,
		login,
	)).Methods(http.MethodPost)
	registerUser := users.NewRegisterUser(userRepository)

	inV1.HandleFunc("/users/register", rtr.handle(
		handler.HttpRequest,
		registerUser,
	)).Methods(http.MethodPost)

	inV1.HandleFunc("/shots", rtr.handle(
		handler.HttpRequest,
		getShot,
	)).Methods(http.MethodGet)

	inV1.HandleFunc("/shots/{id}", rtr.handle(
		handler.HttpRequest,
		getShotID,
	)).Methods(http.MethodGet)

	inV1.HandleFunc("/shots/create", rtr.handle(
		handler.HttpRequest,
		createShot,
	)).Methods(http.MethodPost)

	inV1.HandleFunc("/shots/delete/{id}", rtr.handle(
		handler.HttpRequest,
		deleteShot,
	)).Methods(http.MethodDelete)

	inV1.HandleFunc("/shots/update/{id}", rtr.handle(
		handler.HttpRequest,
		updateShot,
	)).Methods(http.MethodPut)

	// use case
	// healthy := ucase.NewHealthCheck()

	// healthy
	// liveness.HandleFunc("/liveness", rtr.handle(
	// 	handler.HttpRequest,
	// 	healthy,
	// )).Methods(http.MethodGet)

	// greeting := ucase.NewGreeting()

	// root.HandleFunc("/greeting/{name}", rtr.handle(
	// 	handler.HttpRequest,
	// 	greeting,
	// )).Methods(http.MethodPost)

	// db := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	// result :=

	// this is use case for example purpose, please delete
	//repoExample := repositories.NewExample(db)
	//el := example.NewExampleList(repoExample)
	//ec := example.NewPartnerCreate(repoExample
	//ed := example.NewExampleDelete(repoExample)

	// TODO: create your route here

	// this route for example rest, please delete
	// example list
	//inV1.HandleFunc("/example", rtr.handle(
	//    handler.HttpRequest,
	//    el,
	//)).Methods(http.MethodGet)

	//inV1.HandleFunc("/example", rtr.handle(
	//    handler.HttpRequest,
	//    ec,
	//)).Methods(http.MethodPost)

	//inV1.HandleFunc("/example/{id:[0-9]+}", rtr.handle(
	//    handler.HttpRequest,
	//    ed,
	//)).Methods(http.MethodDelete)

	return rtr.router

}
