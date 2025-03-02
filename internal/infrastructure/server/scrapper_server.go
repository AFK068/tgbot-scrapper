package server

import (
	"github.com/AFK068/bot/config"
	api "github.com/AFK068/bot/internal/api/openapi/scrapper/v1"
	"github.com/AFK068/bot/internal/application/scrapper"
	"github.com/AFK068/bot/internal/domain"
	handler "github.com/AFK068/bot/internal/infrastructure/handler/scrapper"
	"github.com/AFK068/bot/internal/infrastructure/logger"
	"github.com/AFK068/bot/internal/middleware"

	"github.com/labstack/echo/v4"
)

type ScrapperServer struct {
	Config    *config.ScrapperConfig
	Handler   *handler.ScrapperHandler
	Scheduler *scrapper.Scrapper
	Echo      *echo.Echo
	Repo      domain.ChatLinkRepository
	Logger    *logger.Logger
}

func NewScrapperServer(
	cfg *config.ScrapperConfig,
	repo domain.ChatLinkRepository,
	hd *handler.ScrapperHandler,
	sd *scrapper.Scrapper,
	log *logger.Logger,
) *ScrapperServer {
	return &ScrapperServer{
		Echo:      echo.New(),
		Config:    cfg,
		Repo:      repo,
		Handler:   hd,
		Scheduler: sd,
		Logger:    log,
	}
}

func (s *ScrapperServer) Start() error {
	api.RegisterHandlers(s.Echo, s.Handler)

	// Middleware for checking the user authentication.
	s.Echo.Use(middleware.AuthLinkMiddleware(s.Repo, s.Logger))

	// Run the scrapper.
	s.Scheduler.Run(scrapper.DefaultJobDuration)

	return s.Echo.Start(s.Config.Host + ":" + s.Config.Port)
}
