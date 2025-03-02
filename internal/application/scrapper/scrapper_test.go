package scrapper_test

import (
	"testing"
	"time"

	botapi "github.com/AFK068/bot/internal/api/openapi/bot/v1"
	"github.com/AFK068/bot/internal/application/scrapper"
	"github.com/AFK068/bot/internal/domain"
	repoMock "github.com/AFK068/bot/internal/domain/mocks"
	botMock "github.com/AFK068/bot/internal/infrastructure/clients/bot/mocks"
	"github.com/AFK068/bot/internal/infrastructure/logger"
	"github.com/AFK068/bot/pkg/client/github"
	githubMock "github.com/AFK068/bot/pkg/client/github/mocks"
	"github.com/AFK068/bot/pkg/client/stackoverflow"
	stackoverflowMock "github.com/AFK068/bot/pkg/client/stackoverflow/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestScrapper_GitHubLinkUpdate(t *testing.T) {
	repo := repoMock.NewChatLinkRepository(t)
	githubClient := githubMock.NewRepoFetcher(t)
	stackoverflowClient := stackoverflowMock.NewQuestionFetcher(t)
	botClient := botMock.NewService(t)

	testLink := &domain.Link{
		URL:       "https://github.com/test/repo",
		Type:      domain.GithubType,
		LastCheck: time.Now().Add(-1 * time.Hour),
	}

	repo.On("GetAllLinks").Return([]*domain.Link{testLink})
	repo.On("GetChatIDsByLink", testLink).Return([]int64{123})
	repo.On("UpdateLastCheck", testLink).Return(nil)

	githubClient.On("GetRepo", mock.Anything, testLink.URL).Return(&github.Repository{
		UpdatedAt: time.Now(),
	}, nil)

	botClient.On("PostUpdates", mock.Anything, mock.MatchedBy(func(update botapi.LinkUpdate) bool {
		return *update.Url == testLink.URL && (*update.TgChatIds)[0] == 123
	})).Return(nil)

	s, err := scrapper.NewScrapperScheduler(repo, stackoverflowClient, githubClient, botClient, logger.NewDiscardLogger())
	assert.NoError(t, err)

	s.Run(time.Second)
	time.Sleep(2 * time.Second)

	repo.AssertExpectations(t)
	githubClient.AssertExpectations(t)
	botClient.AssertExpectations(t)
}

func TestScrapper_StackOverflowLinkUpdate(t *testing.T) {
	repo := repoMock.NewChatLinkRepository(t)
	githubClient := githubMock.NewRepoFetcher(t)
	stackoverflowClient := stackoverflowMock.NewQuestionFetcher(t)
	botClient := botMock.NewService(t)

	testLink := &domain.Link{
		URL:       "https://stackoverflow.com/test/question",
		Type:      domain.StackoverflowType,
		LastCheck: time.Now().Add(-1 * time.Hour),
	}

	repo.On("GetAllLinks").Return([]*domain.Link{testLink})
	repo.On("GetChatIDsByLink", testLink).Return([]int64{123})
	repo.On("UpdateLastCheck", testLink).Return(nil)

	stackoverflowClient.On("GetQuestion", mock.Anything, testLink.URL).Return(&stackoverflow.Question{
		LastActivityDate: time.Now().Unix(),
	}, nil)

	botClient.On("PostUpdates", mock.Anything, mock.MatchedBy(func(update botapi.LinkUpdate) bool {
		return *update.Url == testLink.URL && (*update.TgChatIds)[0] == 123
	})).Return(nil)

	s, err := scrapper.NewScrapperScheduler(repo, stackoverflowClient, githubClient, botClient, logger.NewDiscardLogger())
	assert.NoError(t, err)

	s.Run(time.Second)
	time.Sleep(2 * time.Second)

	repo.AssertExpectations(t)
	stackoverflowClient.AssertExpectations(t)
	botClient.AssertExpectations(t)
}
