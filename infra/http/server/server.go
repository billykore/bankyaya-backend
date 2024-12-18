package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.bankyaya.org/app/backend/infra/http/handler"
	"go.bankyaya.org/app/backend/pkg/config"
	"go.bankyaya.org/app/backend/pkg/logger"
)

// Server to run.
type Server struct {
	router *Router
}

// New creates new Server.
func New(router *Router) *Server {
	return &Server{
		router: router,
	}
}

// Serve start the Server.
func (s *Server) Serve() {
	s.router.Run()
}

// Router get all request to handlers and returns the response produce by handlers.
type Router struct {
	cfg             *config.Config
	log             *logger.Logger
	router          *echo.Echo
	transferHandler *handler.TransferHandler
}

// NewRouter returns new Router.
func NewRouter(
	cfg *config.Config,
	log *logger.Logger,
	router *echo.Echo,
	transferHandler *handler.TransferHandler,
) *Router {
	return &Router{
		cfg:             cfg,
		log:             log,
		router:          router,
		transferHandler: transferHandler,
	}
}

func (r *Router) useMiddlewares() {
	r.router.Use(middleware.Logger())
	r.router.Use(middleware.Recover())
}

func (r *Router) swagger() {
	r.router.GET("/swagger/*", echoSwagger.WrapHandler)
}

func (r *Router) run() {
	port := r.cfg.HTTP.Port
	r.log.Infof("running on port ::[:%v]", port)
	if err := r.router.Start(":" + port); err != nil {
		r.log.Fatalf("failed to run on port [::%v]", port)
	}
}

// Run runs the server.
func (r *Router) Run() {
	r.useMiddlewares()
	r.swagger()
	r.setTransferRoutes()
	r.run()
}

func (r *Router) setTransferRoutes() {
	r.router.POST("/transfer/inquiry", r.transferHandler.Inquiry)
	r.router.POST("/transfer/payment", r.transferHandler.Payment)
}
