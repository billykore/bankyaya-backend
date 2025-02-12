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
	qrisHandler     *handler.QRISHandler
	userHandler     *handler.UserHandler
}

// NewRouter returns new Router.
func NewRouter(
	cfg *config.Config,
	log *logger.Logger,
	router *echo.Echo,
	transferHandler *handler.TransferHandler,
	qrisHandler *handler.QRISHandler,
	userHandler *handler.UserHandler,
) *Router {
	return &Router{
		cfg:             cfg,
		log:             log,
		router:          router,
		transferHandler: transferHandler,
		qrisHandler:     qrisHandler,
		userHandler:     userHandler,
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
	r.log.Usecase("run").Infof("running on port ::[:%v]", port)
	if err := r.router.Start(":" + port); err != nil {
		r.log.Usecase("run").Fatalf("failed to run on port [::%v]", port)
	}
}

// Run runs the server.
func (r *Router) Run() {
	r.useMiddlewares()
	r.swagger()
	r.setTransferRoutes()
	r.setQRISRoutes()
	r.setUserRoutes()
	r.run()
}

func (r *Router) setTransferRoutes() {
	tr := r.router.Group("/transfer")
	tr.Use(authMiddleware())

	tr.POST("/inquiry", r.transferHandler.Inquiry)
	tr.POST("/payment", r.transferHandler.Payment)
}

func (r *Router) setQRISRoutes() {
	qr := r.router.Group("/qris")
	qr.Use(authMiddleware())

	qr.POST("/inquiry", r.qrisHandler.Inquiry)
	qr.POST("/payment", r.qrisHandler.Payment)
}

func (r *Router) setUserRoutes() {
	r.router.POST("/user/login", r.userHandler.Login)
}
