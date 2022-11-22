// Package server
package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	gohandlers "github.com/gorilla/handlers"

	"gitlab.privy.id/go_graphql/internal/appctx"
	"gitlab.privy.id/go_graphql/internal/consts"
	"gitlab.privy.id/go_graphql/internal/router"
	"gitlab.privy.id/go_graphql/pkg/logger"
)

// NewHTTPServer creates http server instance
// returns: Server instance
func NewHTTPServer() Server {
	cfg := appctx.NewConfig()
	return &httpServer{
		config: cfg,
		router: router.NewRouter(cfg),
	}
}

// httpServer as HTTP server implementation
type httpServer struct {
	config *appctx.Config
	router router.Router
}

// Run runs the http server gracefully
// returns:
//
//	err: error operation
func (h *httpServer) Run(ctx context.Context) error {
	var err error

	ch := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"*"}))
	server := http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%d", h.config.App.Port),
		Handler:      ch(h.router.Route()),
		ReadTimeout:  time.Duration(h.config.App.ReadTimeoutSecond) * time.Second,
		WriteTimeout: time.Duration(h.config.App.WriteTimeoutSecond) * time.Second,
	}

	go func() {
		err = server.ListenAndServe()
		if err != http.ErrServerClosed {
			logger.Error(logger.MessageFormat("http server got %v", err), logger.EventName(consts.LogEventNameServiceStarting))
		}
	}()

	<-ctx.Done()

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer func() {
		cancel()
	}()

	if err = server.Shutdown(ctxShutDown); err != nil {
		logger.Fatal(logger.MessageFormat("server Shutdown Failed:%v", err), logger.EventName(consts.LogEventNameServiceTerminated))
	}

	logger.Info("server exited properly", logger.EventName(consts.LogEventNameServiceTerminated))

	if err == http.ErrServerClosed {
		err = nil
	}

	return err
}

// Done runs event wheen service stopped
func (h *httpServer) Done() {
	logger.Info("service http stopped", logger.EventName(consts.LogEventNameServiceTerminated))
}

// Config  func to handle get config will return Config object
func (h *httpServer) Config() *appctx.Config {
	return h.config
}
