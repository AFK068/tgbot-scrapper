package server

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"

	"github.com/AFK068/bot/internal/application/scrapper"
	"github.com/AFK068/bot/internal/domain"
	"github.com/AFK068/bot/internal/infrastructure/httpapi/scrapperapi"
	"github.com/AFK068/bot/internal/infrastructure/logger"
	"github.com/AFK068/bot/internal/middleware"

	scrappertypes "github.com/AFK068/bot/internal/api/openapi/scrapper/v1"
)

type ScrapperServer struct {
	Config    *scrapper.Config
	Handler   *scrapperapi.ScrapperHandler
	Scheduler *scrapper.Scrapper
	Echo      *echo.Echo
	Repo      domain.ChatLinkRepository
	Logger    *logger.Logger
}

func NewScrapperServer(
	cfg *scrapper.Config,
	repo domain.ChatLinkRepository,
	hd *scrapperapi.ScrapperHandler,
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
	scrappertypes.RegisterHandlers(s.Echo, s.Handler)

	// Middleware for checking the user authentication.
	s.Echo.Use(middleware.AuthLinkMiddleware(s.Repo, s.Logger))

	// Run the scrapper.
	s.Scheduler.Run(scrapper.DefaultJobDuration)

	return s.Echo.Start(s.Config.Host + ":" + s.Config.Port)
}

func (s *ScrapperServer) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return s.Echo.Shutdown(ctx)
}

func (s *ScrapperServer) RegisterHooks(lc fx.Lifecycle, log *logger.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			log.Logger.Info("Starting scrapper server")

			go func() {
				if err := s.Start(); err != nil && err != http.ErrServerClosed {
					log.Error("Failed to start scrapper server", "error", err)
				}
			}()

			return nil
		},
		OnStop: func(context.Context) error {
			log.Logger.Info("Stopping scrapper server")

			if err := s.Stop(); err != nil {
				log.Error("Failed to stop scrapper server", "error", err)
			}

			return nil
		},
	})
}
