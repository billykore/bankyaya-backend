package server

import (
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	echoswagger "github.com/swaggo/echo-swagger"
	"go.bankyaya.org/app/backend/internal/adapter/http/handler"
	"go.bankyaya.org/app/backend/internal/adapter/http/middleware"
	"go.bankyaya.org/app/backend/internal/pkg/config"
	"go.bankyaya.org/app/backend/internal/pkg/logger"
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

// Router gets all requests to handlers and returns the response produce by handlers.
type Router struct {
	cfg              *config.Configs
	log              *logger.Logger
	router           *echo.Echo
	intrabankHandler *handler.Intrabank
	userHandler      *handler.UserHandler
	otpHandler       *handler.OTPHandler
}

// NewRouter returns new Router.
func NewRouter(
	cfg *config.Configs,
	log *logger.Logger,
	router *echo.Echo,
	transferHandler *handler.Intrabank,
	userHandler *handler.UserHandler,
	otpHandler *handler.OTPHandler,
) *Router {
	return &Router{
		cfg:              cfg,
		log:              log,
		router:           router,
		intrabankHandler: transferHandler,
		userHandler:      userHandler,
		otpHandler:       otpHandler,
	}
}

func (r *Router) useMiddlewares() {
	r.router.Use(echomiddleware.Logger())
	r.router.Use(echomiddleware.Recover())
}

func (r *Router) swagger() {
	r.router.GET("/swagger/*", echoswagger.WrapHandler)
}

func (r *Router) run() {
	port := r.cfg.App.Port
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
	r.setUserRoutes()
	r.run()
}

func (r *Router) setTransferRoutes() {
	tr := r.router.Group("/transfer/intrabank")
	tr.Use(middleware.AuthenticateUser())

	tr.POST("/inquiry", r.intrabankHandler.Inquiry)
	tr.POST("/payment", r.intrabankHandler.Payment)
}

func (r *Router) setUserRoutes() {
	r.router.POST("/user/login", r.userHandler.Login)
}

func (r *Router) setOTPRoutes() {
	or := r.router.Group("/otp")
	or.POST("/send", r.otpHandler.SendOTP)
	or.POST("/verify", r.otpHandler.VerifyOTP)
}
