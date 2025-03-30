package repository_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"

	"github.com/AFK068/bot/internal/config"
	"github.com/AFK068/bot/internal/domain"
	"github.com/AFK068/bot/internal/infrastructure/repository"
	"github.com/AFK068/bot/internal/infrastructure/repository/link/ormrepo"
	"github.com/AFK068/bot/internal/infrastructure/repository/link/sqlrepo"
)

type mockLC struct{}

func (m *mockLC) Append(_ fx.Hook) {}

func TestRepoCreation(t *testing.T) {
	testCases := []struct {
		cfg      *config.Config
		expected interface{}
	}{
		{
			cfg: &config.Config{
				Storage: config.Storage{
					Type: domain.ORMRepository,
				},
			},
			expected: &ormrepo.Repository{},
		},
		{
			cfg: &config.Config{
				Storage: config.Storage{
					Type: domain.DirectSQLRepository,
				},
			},
			expected: &sqlrepo.Repository{},
		},
	}

	for _, tc := range testCases {
		t.Run(string(tc.cfg.Storage.Type), func(t *testing.T) {
			repo, _, err := repository.NewPostgresRepo(tc.cfg, &mockLC{})
			assert.NoError(t, err)
			assert.IsType(t, tc.expected, repo)
		})
	}
}
