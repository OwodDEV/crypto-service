package http

import (
	"log/slog"

	"github.com/OwodDEV/crypto-service/internal/config"
	"github.com/OwodDEV/crypto-service/internal/metrics"
	"github.com/OwodDEV/crypto-service/internal/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/swagger"

	_ "github.com/OwodDEV/crypto-service/docs/swagger"
)

type Server struct {
	Service  *service.Service
	Metrics  *metrics.Metrics
	Config   *config.Config
	Validate *validator.Validate
	router   *fiber.App
}

func NewServer(srv *service.Service, mtr *metrics.Metrics, cfg *config.Config) (s *Server, err error) {
	s = &Server{
		Service:  srv,
		Metrics:  mtr,
		Config:   cfg,
		Validate: validator.New(),
	}
	return
}

// @title Auth Service API
// @version
// @description
// @BasePath /
func (s *Server) Run(errCh chan<- error) {
	s.router = fiber.New(fiber.Config{
		DisableStartupMessage:   true,
		CaseSensitive:           true,
		StrictRouting:           true,
		EnableTrustedProxyCheck: true,
	})

	// middleware
	s.router.Use("/api/*", s.TraceMiddleware())
	s.router.Use("/api/*", s.RequestLoggerMiddleware())

	// api routes
	s.router.Get("/api/wallet/:address", s.GetWalletHandler)
	s.router.Get("/api/transaction/:hash", s.GetTransactionHandler)

	// swagger
	s.router.Get("/swagger/*", swagger.HandlerDefault)

	// metrics
	s.router.Get("/metrics", adaptor.HTTPHandler(s.Metrics.PrometheusHandler()))

	// startup
	slog.Info(
		"starting the HTTP server",
		slog.String("port", s.Config.Transport.HTTP.Port),
	)
	err := s.router.Listen(s.Config.Transport.HTTP.Host + ":" + s.Config.Transport.HTTP.Port)
	if err != nil {
		slog.Error("HTTP router.Listen() failed", slog.Any("err", err))
		errCh <- err
	}
}

func (s *Server) Shutdown() {
	slog.Info("shutting down HTTP server...")
	s.router.Shutdown()
}
