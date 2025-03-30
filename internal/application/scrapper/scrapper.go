package scrapper

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/go-co-op/gocron/v2"

	"github.com/AFK068/bot/internal/domain"
	"github.com/AFK068/bot/internal/infrastructure/clients/bot"
	"github.com/AFK068/bot/internal/infrastructure/logger"
	"github.com/AFK068/bot/pkg/client/github"
	"github.com/AFK068/bot/pkg/client/stackoverflow"
	"github.com/AFK068/bot/pkg/utils"

	bottypes "github.com/AFK068/bot/internal/api/openapi/bot/v1"
)

const (
	DefaultJobDuration = 15 * time.Second

	PaginationLimit uint64 = 50
)

type StackOverlowQuestionFetcher interface {
	GetQuestion(ctx context.Context, questionURL string) (*stackoverflow.Question, error)
	GetActivity(ctx context.Context, question *stackoverflow.Question, lastCheckTime time.Time) ([]*stackoverflow.Activity, error)
}

type GitHubRepoFetcher interface {
	GetRepo(ctx context.Context, questionURL string) (*github.Repository, error)
	GetActivity(ctx context.Context, repository *github.Repository, lastCheckTime time.Time) ([]*github.Activity, error)
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

func (s *Scrapper) Stop() error {
	err := s.scheduler.Shutdown()
	if err != nil {
		s.logger.Error("Failed to stop scheduler", "error", err)
		return fmt.Errorf("failed to stop scheduler: %w", err)
	}

	s.logger.Info("Scheduler stopped")

	return nil
}

func (s *Scrapper) notifyBot(ctx context.Context, activities []*domain.Activity, link *domain.Link) error {
	s.logger.Info("Notifying bot for link", "url", link.URL)

	chatIDs, err := s.repository.GetChatIDsByLink(ctx, link)
	if err != nil {
		s.logger.Error("Error getting chat IDs", "error", err)
		return fmt.Errorf("error getting chat IDs: %w", err)
	}

	if len(chatIDs) == 0 {
		s.logger.Warn("No chat IDs found for link", "url", link.URL)
		return nil
	}

	for _, activity := range activities {
		activityType := activity.MapActivityTypeToBotAPI()
		if activityType == nil {
			s.logger.Error("Activity type is nil", "activity", activity)
			return fmt.Errorf("invalid activity type: %v", activity.Type)
		}

		userName := "Unknown"
		if activity.UserName != "" {
			userName = activity.UserName
		}

		description := "No description"
		if activity.Body != "" {
			description = activity.Body
		}

		update := bottypes.LinkUpdate{
			TgChatIds:   utils.SliceInt64Ptr(chatIDs),
			СreatedAt:   &activity.CreatedAt,
			Type:        activityType,
			Url:         aws.String(link.URL),
			UserName:    aws.String(userName),
			Description: aws.String(description),
		}

		if err := s.botClient.PostUpdates(ctx, update); err != nil {
			s.logger.Error("Error posting update to bot", "error", err)
			return fmt.Errorf("error posting update to bot: %w", err)
		}

		s.logger.Info(*update.Description)
	}

	return nil
}

func (s *Scrapper) getActivity(ctx context.Context, link *domain.Link) ([]*domain.Activity, error) {
	s.logger.Info("Checking link for update", "url", link.URL)

	switch link.Type {
	case domain.StackoverflowType:
		return s.getStackOverflowActivity(ctx, link)
	case domain.GithubType:
		return s.getGitHubActivity(ctx, link)
	default:
		s.logger.Error("Unknown link type", "type", link.Type)
		return nil, fmt.Errorf("unknown link type: %s", link.Type)
	}
}

func (s *Scrapper) getStackOverflowActivity(ctx context.Context, link *domain.Link) ([]*domain.Activity, error) {
	s.logger.Info("Checking StackOverflow link for update", "url", link.URL)

	question, err := s.stackOverflowClient.GetQuestion(ctx, link.URL)
	if err != nil {
		s.logger.Error("Failed to get question", "error", err)
		return nil, fmt.Errorf("failed to get question: %w", err)
	}

	var activities []*domain.Activity

	if question.LastActivityDate > link.LastCheck.Unix() {
		activity, err := s.stackOverflowClient.GetActivity(ctx, question, link.LastCheck)
		if err != nil {
			s.logger.Error("Failed to get activity", "error", err)
			return nil, fmt.Errorf("failed to get activity: %w", err)
		}

		for _, act := range activity {
			var activityType domain.ActivityType

			switch act.Type {
			case stackoverflow.ActivityTypeAnswer:
				activityType = domain.StackoverflowAnswer
			case stackoverflow.ActivityTypeQuestion:
				activityType = domain.StackoverflowQuestion
			case stackoverflow.ActivityTypeComment:
				activityType = domain.StackoverflowComment
			default:
				s.logger.Error("Unknown activity type", "type", act.Type)
				return nil, fmt.Errorf("unknown activity type: %s", act.Type)
			}

			activities = append(activities, domain.NewActivity(activityType, "", time.Unix(act.CreatedAt, 0), act.Body, act.UserName))
		}
	}

	return activities, nil
}

func (s *Scrapper) getGitHubActivity(ctx context.Context, link *domain.Link) ([]*domain.Activity, error) {
	s.logger.Info("Checking GitHub link for update", "url", link.URL)

	repo, err := s.gitHubClient.GetRepo(ctx, link.URL)
	if err != nil {
		s.logger.Error("Failed to get repository", "error", err)
		return nil, fmt.Errorf("failed to get repository: %w", err)
	}

	var activities []*domain.Activity

	if repo.UpdatedAt.After(link.LastCheck) {
		activity, err := s.gitHubClient.GetActivity(ctx, repo, link.LastCheck)
		if err != nil {
			s.logger.Error("Failed to get activity", "error", err)
			return nil, fmt.Errorf("failed to get activity: %w", err)
		}

		for _, act := range activity {
			var activityType domain.ActivityType

			switch act.Type {
			case github.ActivityTypeIssue:
				activityType = domain.GitHubIssue
			case github.ActivityTypePullRequest:
				activityType = domain.GitHubPullRequest
			case github.ActivityTypeRepository:
				activityType = domain.GitHubPullRequest
			default:
				s.logger.Error("Unknown activity type", "type", act.Type)
				return nil, fmt.Errorf("unknown activity type: %s", act.Type)
			}

			activities = append(activities, domain.NewActivity(activityType, act.Title, act.CreatedAt, act.Body, act.UserName))
		}
	}

	return activities, nil
}

func (s *Scrapper) scrappeLinksTask() {
	s.logger.Info("Starting scrappeLinksTask")

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	sema := make(chan struct{}, runtime.NumCPU()*4)

	var wg sync.WaitGroup

	offset := uint64(0)

	for {
		select {
		case <-ctx.Done():
			s.logger.Info("Context done, stopping scrappeLinksTask")
			wg.Wait()

			return
		default:
			links, err := s.repository.GetLinksPagination(ctx, offset, PaginationLimit)
			if err != nil {
				s.logger.Error("Error getting links for pagination", "error", err)
				wg.Wait()

				return
			}

			wg.Add(1)

			go func(links []*domain.Link) {
				defer wg.Done()

				select {
				case sema <- struct{}{}:
					defer func() { <-sema }()
				case <-ctx.Done():
					s.logger.Warn("Context done before processing link")
					return
				}

				for _, link := range links {
					if ctx.Err() != nil {
						s.logger.Warn("Context error", "error", ctx.Err())
						return
					}

					if err := s.processLink(ctx, link); err != nil {
						s.logger.Error("Error processing link", "url", link.URL, "error", err)
						return
					}
				}
			}(links)

			if len(links) < int(PaginationLimit) {
				s.logger.Info("No more links to process, stopping scrappeLinksTask")
				wg.Wait()

				return
			}

			offset += PaginationLimit
		}
	}
}

func (s *Scrapper) processLink(ctx context.Context, link *domain.Link) error {
	activities, err := s.getActivity(ctx, link)
	if err != nil {
		return err
	}

	if len(activities) == 0 {
		s.logger.Info("No new activities found for link", "url", link.URL)
		return nil
	}

	if err := s.notifyBot(ctx, activities, link); err != nil {
		s.logger.Error("Error notifying bot", "error", err)
		return err
	}

	if err := s.repository.UpdateLastCheck(ctx, link); err != nil {
		s.logger.Error("Error updating last check", "error", err)
		return err
	}

	s.logger.Info("Successfully processed link", "url", link.URL)

	return nil
}
