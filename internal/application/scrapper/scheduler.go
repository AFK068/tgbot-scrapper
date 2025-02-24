package scrapper

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"

	botapi "github.com/AFK068/bot/internal/api/openapi/bot/v1"
	"github.com/AFK068/bot/internal/domain"
	"github.com/AFK068/bot/internal/infrastructure/clients/bot"
	"github.com/AFK068/bot/pkg/client/github"
	"github.com/AFK068/bot/pkg/client/stackoverflow"
	"github.com/AFK068/bot/pkg/utils"
	"github.com/go-co-op/gocron/v2"
)

type Scheduler struct {
	scheduler           gocron.Scheduler
	repository          domain.ChatLinkRepository
	stackOverflowClient *stackoverflow.Client
	gitHubClient        *github.Client
	botClient           *bot.Client
}

func NewScrapperScheduler(
	repository domain.ChatLinkRepository,
	stackoverflowClient *stackoverflow.Client,
	githubClient *github.Client,
	botClient *bot.Client,
) (*Scheduler, error) {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		return nil, fmt.Errorf("failed to create scheduler: %w", err)
	}

	return &Scheduler{
		scheduler:           scheduler,
		repository:          repository,
		stackOverflowClient: stackoverflowClient,
		gitHubClient:        githubClient,
		botClient:           botClient,
	}, nil
}

func (s *Scheduler) Run(jobDuration time.Duration) {
	_, err := s.scheduler.NewJob(
		gocron.DurationJob(
			jobDuration,
		),
		gocron.NewTask(
			s.scrappeLinksTask,
		),
	)

	if err != nil {
		fmt.Println("failed to create job:", err)
	}

	s.scheduler.Start()
}

func (s *Scheduler) Stop() error {
	if err := s.scheduler.Shutdown(); err != nil {
		return fmt.Errorf("failed to shutdown scheduler: %w", err)
	}

	return nil
}

func (s *Scheduler) notifyBot(ctx context.Context, link *domain.Link) error {
	update := botapi.LinkUpdate{
		Url:       &link.URL,
		TgChatIds: utils.SliceInt64Ptr(s.repository.GetChatIDsByLink(link)),
	}

	if err := s.botClient.PostUpdates(ctx, update); err != nil {
		return fmt.Errorf("failed to post updates: %w", err)
	}

	return nil
}

func (s *Scheduler) checkLinkForUpdate(ctx context.Context, link *domain.Link) (bool, error) {
	switch link.Type {
	case domain.StackoverflowType:
		question, err := s.stackOverflowClient.GetQuestion(ctx, link.URL)
		if err != nil {
			return false, fmt.Errorf("failed to get question: %w", err)
		}

		return time.Unix(question.LastActivityDate, 0).After(link.LastCheck), nil
	case domain.GithubType:
		repo, err := s.gitHubClient.GetRepo(ctx, link.URL)
		if err != nil {
			return false, fmt.Errorf("failed to get repo: %w", err)
		}

		return repo.UpdatedAt.After(link.LastCheck), nil
	default:
		return false, fmt.Errorf("unknown link type: %s", link.Type)
	}
}

func (s *Scheduler) scrappeLinksTask() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	links := s.repository.GetAllLinks()
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
				return
			}

			if ctx.Err() != nil {
				return
			}

			// Check if the link needs to be updated.
			needUpdate, err := s.checkLinkForUpdate(ctx, l)
			if err != nil {
				return
			}

			// Skip if the link does not need to be updated.
			if !needUpdate {
				return
			}

			// Notify the bot about the update.
			if err := s.notifyBot(ctx, l); err != nil {
				return
			}

			// Update the last check time.
			err = s.repository.UpdateLastCheck(l)
			if err != nil {
				return
			}
		}(link)
	}

	wg.Wait()
}
