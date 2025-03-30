package scrapper_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/AFK068/bot/internal/application/scrapper"
	"github.com/AFK068/bot/internal/domain"
	"github.com/AFK068/bot/internal/infrastructure/logger"
	"github.com/AFK068/bot/pkg/client/github"
	"github.com/AFK068/bot/pkg/client/stackoverflow"

	bottypes "github.com/AFK068/bot/internal/api/openapi/bot/v1"
	scrapperMock "github.com/AFK068/bot/internal/application/scrapper/mocks"
	repoMock "github.com/AFK068/bot/internal/domain/mocks"
	botMock "github.com/AFK068/bot/internal/infrastructure/clients/bot/mocks"
)

func Test_GitHubLink_Update_Success(t *testing.T) {
	repo := repoMock.NewChatLinkRepository(t)
	githubClient := scrapperMock.NewGitHubRepoFetcher(t)
	stackoverflowClient := scrapperMock.NewStackOverlowQuestionFetcher(t)
	botClient := botMock.NewService(t)

	testLink := &domain.Link{
		UserAddID: 123,
		URL:       "https://stackoverflow.com/test/question",
		Type:      domain.GithubType,
		LastCheck: time.Now().Add(-1 * time.Hour),
	}

	repo.On("GetLinksPagination", mock.Anything, uint64(0), scrapper.PaginationLimit).Return([]*domain.Link{testLink}, nil)

	githubRepo := &github.Repository{
		UpdatedAt: time.Now().Add(-1 * time.Hour),
	}

	githubClient.On("GetRepo", mock.Anything, testLink.URL).Return(githubRepo, nil)

	githubClient.On("GetActivity", mock.Anything, githubRepo, testLink.LastCheck).Return([]*github.Activity{
		{
			Type:      github.ActivityTypeIssue,
			Body:      "Test answer body",
			UserName:  "TestUser",
			CreatedAt: time.Now(),
		},
	}, nil)

	repo.On("GetChatIDsByLink", mock.Anything, testLink).Return([]int64{123}, nil)

	botClient.On("PostUpdates", mock.Anything, mock.MatchedBy(func(update bottypes.LinkUpdate) bool {
		return *update.Url == testLink.URL && (*update.TgChatIds)[0] == 123 &&
			*update.Description == "Test answer body" && *update.UserName == "TestUser"
	})).Return(nil)

	repo.On("UpdateLastCheck", mock.Anything, testLink).Return(nil)

	s, err := scrapper.NewScrapperScheduler(repo, stackoverflowClient, githubClient, botClient, logger.NewDiscardLogger())
	assert.NoError(t, err)

	s.Run(time.Second)
	time.Sleep(2 * time.Second)

	repo.AssertExpectations(t)
	stackoverflowClient.AssertExpectations(t)
	botClient.AssertExpectations(t)
}

func Test_GitHubLink_NoUpdate_Success(t *testing.T) {
	repo := repoMock.NewChatLinkRepository(t)
	githubClient := scrapperMock.NewGitHubRepoFetcher(t)
	stackoverflowClient := scrapperMock.NewStackOverlowQuestionFetcher(t)
	botClient := botMock.NewService(t)

	testLink := &domain.Link{
		UserAddID: 123,
		URL:       "https://stackoverflow.com/test/question",
		Type:      domain.GithubType,
		LastCheck: time.Now(),
	}

	repo.On("GetLinksPagination", mock.Anything, uint64(0), scrapper.PaginationLimit).Return([]*domain.Link{testLink}, nil)

	githubRepo := &github.Repository{
		UpdatedAt: time.Now().Add(-1 * time.Hour),
	}

	githubClient.On("GetRepo", mock.Anything, testLink.URL).Return(githubRepo, nil)

	s, err := scrapper.NewScrapperScheduler(repo, stackoverflowClient, githubClient, botClient, logger.NewDiscardLogger())
	assert.NoError(t, err)

	s.Run(time.Second)
	time.Sleep(2 * time.Second)

	repo.AssertExpectations(t)
	stackoverflowClient.AssertExpectations(t)
	botClient.AssertExpectations(t)
}

func Test_StackOverflowLink_Update_Success(t *testing.T) {
	repo := repoMock.NewChatLinkRepository(t)
	githubClient := scrapperMock.NewGitHubRepoFetcher(t)
	stackoverflowClient := scrapperMock.NewStackOverlowQuestionFetcher(t)
	botClient := botMock.NewService(t)

	testLink := &domain.Link{
		UserAddID: 123,
		URL:       "https://stackoverflow.com/test/question",
		Type:      domain.StackoverflowType,
		LastCheck: time.Now().Add(-1 * time.Hour),
	}

	repo.On("GetLinksPagination", mock.Anything, uint64(0), scrapper.PaginationLimit).Return([]*domain.Link{testLink}, nil)

	question := &stackoverflow.Question{
		LastActivityDate: time.Now().Unix(),
	}

	stackoverflowClient.On("GetQuestion", mock.Anything, testLink.URL).Return(question, nil)

	stackoverflowClient.On("GetActivity", mock.Anything, question, testLink.LastCheck).Return([]*stackoverflow.Activity{
		{
			Type:      stackoverflow.ActivityTypeAnswer,
			Body:      "Test answer body",
			UserName:  "TestUser",
			CreatedAt: time.Now().Unix(),
			Tags:      []string{"test", "tags"},
		},
	}, nil)

	repo.On("GetChatIDsByLink", mock.Anything, testLink).Return([]int64{123}, nil)

	botClient.On("PostUpdates", mock.Anything, mock.MatchedBy(func(update bottypes.LinkUpdate) bool {
		return *update.Url == testLink.URL && (*update.TgChatIds)[0] == 123 &&
			*update.Description == "Test answer body" && *update.UserName == "TestUser"
	})).Return(nil)

	repo.On("UpdateLastCheck", mock.Anything, testLink).Return(nil)

	s, err := scrapper.NewScrapperScheduler(repo, stackoverflowClient, githubClient, botClient, logger.NewDiscardLogger())
	assert.NoError(t, err)

	s.Run(time.Second)
	time.Sleep(2 * time.Second)

	repo.AssertExpectations(t)
	stackoverflowClient.AssertExpectations(t)
	botClient.AssertExpectations(t)
}

func Test_StackOverflowLink_NoUpdate_Success(t *testing.T) {
	repo := repoMock.NewChatLinkRepository(t)
	githubClient := scrapperMock.NewGitHubRepoFetcher(t)
	stackoverflowClient := scrapperMock.NewStackOverlowQuestionFetcher(t)
	botClient := botMock.NewService(t)

	testLink := &domain.Link{
		UserAddID: 123,
		URL:       "https://stackoverflow.com/test/question",
		Type:      domain.StackoverflowType,
		LastCheck: time.Now(),
	}

	repo.On("GetLinksPagination", mock.Anything, uint64(0), scrapper.PaginationLimit).Return([]*domain.Link{testLink}, nil)

	question := &stackoverflow.Question{
		LastActivityDate: time.Now().Add(-1 * time.Hour).Unix(),
	}

	stackoverflowClient.On("GetQuestion", mock.Anything, testLink.URL).Return(question, nil)

	s, err := scrapper.NewScrapperScheduler(repo, stackoverflowClient, githubClient, botClient, logger.NewDiscardLogger())
	assert.NoError(t, err)

	s.Run(time.Second)
	time.Sleep(2 * time.Second)

	repo.AssertExpectations(t)
	stackoverflowClient.AssertExpectations(t)
	botClient.AssertExpectations(t)
}

func Test_Pagination_Success(t *testing.T) {
	repo := repoMock.NewChatLinkRepository(t)
	githubClient := scrapperMock.NewGitHubRepoFetcher(t)
	stackoverflowClient := scrapperMock.NewStackOverlowQuestionFetcher(t)
	botClient := botMock.NewService(t)

	batch1 := make([]*domain.Link, 50)
	batch2 := make([]*domain.Link, 50)
	batch3 := make([]*domain.Link, 40)

	for i := range batch1 {
		batch1[i] = &domain.Link{
			URL:       fmt.Sprintf("https://example.com/%d", i),
			Type:      domain.GithubType,
			LastCheck: time.Now(),
		}

		batch2[i] = &domain.Link{
			URL:       fmt.Sprintf("https://example.com/%d", 50+i),
			Type:      domain.GithubType,
			LastCheck: time.Now(),
		}
	}

	for i := range batch3 {
		batch3[i] = &domain.Link{
			URL:       fmt.Sprintf("https://example.com/%d", 100+i),
			Type:      domain.GithubType,
			LastCheck: time.Now(),
		}
	}

	repo.On("GetLinksPagination", mock.Anything, uint64(0), scrapper.PaginationLimit).Return(batch1, nil).Once()
	repo.On("GetLinksPagination", mock.Anything, uint64(50), scrapper.PaginationLimit).Return(batch2, nil).Once()
	repo.On("GetLinksPagination", mock.Anything, uint64(100), scrapper.PaginationLimit).Return(batch3, nil).Once()

	githubClient.On("GetRepo", mock.Anything, mock.Anything).Return(&github.Repository{
		UpdatedAt: time.Now().Add(-1 * time.Hour),
	}, nil).Times(140)

	s, err := scrapper.NewScrapperScheduler(
		repo,
		stackoverflowClient,
		githubClient,
		botClient,
		logger.NewDiscardLogger(),
	)
	assert.NoError(t, err)

	s.Run(500 * time.Millisecond)
	time.Sleep(1 * time.Second)

	err = s.Stop()
	assert.NoError(t, err)

	repo.AssertNumberOfCalls(t, "GetLinksPagination", 3)
	repo.AssertExpectations(t)
}
