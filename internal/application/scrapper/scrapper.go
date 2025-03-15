package scrapper

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/go-co-op/gocron/v2"

	"github.com/AFK068/bot/internal/domain"
	"github.com/AFK068/bot/internal/infrastructure/clients/bot"
	"github.com/AFK068/bot/internal/infrastructure/logger"
	"github.com/AFK068/bot/pkg/client/github"
	"github.com/AFK068/bot/pkg/client/stackoverflow"
	"github.com/AFK068/bot/pkg/utils"

	botapi "github.com/AFK068/bot/internal/api/openapi/bot/v1"
)

const (
	DefaultJobDuration = 10 * time.Minute
)

type StackOverlowQuestionFetcher interface {
	GetQuestion(ctx context.Context, questionURL string) (*stackoverflow.Question, error)
}

type GitHubRepoFetcher interface {
	GetRepo(ctx context.Context, questionURL string) (*github.Repository, error)
}

type Scrapper struct {
	scheduler           gocron.Scheduler
	repository          domain.ChatLinkRepository
	stackOverflowClient StackOverlowQuestionFetcher
	gitHubClient        GitHubRepoFetcher
	botClient           bot.Service
	logger              *logger.Logger
}

func NewScrapperScheduler(
	repository domain.ChatLinkRepository,
	stackoverflowClient StackOverlowQuestionFetcher,
	githubClient GitHubRepoFetcher,
	botClient bot.Service,
	log *logger.Logger,
) (*Scrapper, error) {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		return nil, fmt.Errorf("failed to create scheduler: %w", err)
	}

	return &Scrapper{
		scheduler:           scheduler,
		repository:          repository,
		stackOverflowClient: stackoverflowClient,
		gitHubClient:        githubClient,
		botClient:           botClient,
		logger:              log,
	}, nil
}

func (s *Scrapper) Run(jobDuration time.Duration) {
	s.logger.Info("Starting scrapper", "jobDuration", jobDuration.String())
	_, err := s.scheduler.NewJob(
		gocron.DurationJob(
			jobDuration,
		),
		gocron.NewTask(
			s.scrappeLinksTask,
		),
	)

	if err != nil {
		s.logger.Error("Failed to create new job", "error", err)
		return
	}

	s.scheduler.Start()
	s.logger.Info("Scrapper started")
}

func (s *Scrapper) notifyBot(ctx context.Context, link *domain.Link) error {
	s.logger.Info("Notifying bot for link", "url", link.URL)

	chatIDs, err := s.repository.GetChatIDsByLink(ctx, link)
	if err != nil {
		s.logger.Error("Failed to get chat IDs by link", "error", err)
		return fmt.Errorf("failed to get chat IDs by link: %w", err)
	}

	update := botapi.LinkUpdate{
		Url:       &link.URL,
		TgChatIds: utils.SliceInt64Ptr(chatIDs),
	}

	if err := s.botClient.PostUpdates(ctx, update); err != nil {
		s.logger.Error("Failed to post updates", "error", err)
		return fmt.Errorf("failed to post updates: %w", err)
	}

	return nil
}

func (s *Scrapper) checkLinkForUpdate(ctx context.Context, link *domain.Link) (bool, error) {
	s.logger.Info("Checking link for update", "url", link.URL)

	switch link.Type {
	case domain.StackoverflowType:
		question, err := s.stackOverflowClient.GetQuestion(ctx, link.URL)
		if err != nil {
			s.logger.Error("Failed to get question", "error", err)
			return false, fmt.Errorf("failed to get question: %w", err)
		}

		return time.Unix(question.LastActivityDate, 0).After(link.LastCheck), nil
	case domain.GithubType:
		repo, err := s.gitHubClient.GetRepo(ctx, link.URL)
		if err != nil {
			s.logger.Error("Failed to get repo", "error", err)
			return false, fmt.Errorf("failed to get repo: %w", err)
		}

		return repo.UpdatedAt.After(link.LastCheck), nil
	default:
		s.logger.Error("Unknown link type", "type", link.Type)
		return false, fmt.Errorf("unknown link type: %s", link.Type)
	}
}

func (s *Scrapper) scrappeLinksTask() {
	s.logger.Info("Starting scrappeLinksTask")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	links, err := s.repository.GetAllLinks(ctx)
	if err != nil {
		s.logger.Error("Failed to get all links", "error", err)
		return
	}

	sema := make(chan struct{}, runtime.NumCPU()*4)

	var wg sync.WaitGroup

	for _, link := range links {
		wg.Add(1)

		go func(l *domain.Link) {
			defer wg.Done()

			select {
			case sema <- struct{}{}:
				defer func() { <-sema }()
			case <-ctx.Done():
				s.logger.Warn("Context done before processing link", "url", l.URL)
				return
			}

			if ctx.Err() != nil {
				s.logger.Warn("Context error", "error", ctx.Err())
				return
			}

			// Check if the link needs to be updated.
			needUpdate, err := s.checkLinkForUpdate(ctx, l)
			if err != nil {
				s.logger.Error("Error checking link for update", "error", err)
				return
			}

			// Skip if the link does not need to be updated.
			if !needUpdate {
				s.logger.Info("No update needed for link", "url", l.URL)
				return
			}

			// Notify the bot about the update.
			if err := s.notifyBot(ctx, l); err != nil {
				s.logger.Error("Error notifying bot: ", "error", err)
				return
			}

			// Update the last check time.
			err = s.repository.UpdateLastCheck(ctx, l)
			if err != nil {
				s.logger.Error("Error updating last check", "error", err)
				return
			}

			s.logger.Info("Successfully processed link", "url", l.URL)
		}(link)
	}

	wg.Wait()
	s.logger.Info("Finished scrappeLinksTask")
}
