package mapper_test

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/AFK068/bot/internal/application/mapper"
	"github.com/AFK068/bot/internal/domain"
	"github.com/AFK068/bot/internal/domain/apperrors"

	api "github.com/AFK068/bot/internal/api/openapi/scrapper/v1"
)

func TestMapAddLinkRequestToDomain_Success(t *testing.T) {
	type args struct {
		userID  int64
		request *api.AddLinkRequest
	}

	tests := []struct {
		name     string
		args     args
		wantType string
	}{
		{
			name: "GitHub link success",
			args: args{
				userID: 1,
				request: &api.AddLinkRequest{
					Link:    aws.String("https://github.com/test"),
					Tags:    &[]string{"tag"},
					Filters: &[]string{"filter"},
				},
			},
			wantType: domain.GithubType,
		},
		{
			name: "StackOverflow link success",
			args: args{
				userID: 1,
				request: &api.AddLinkRequest{
					Link:    aws.String("https://stackoverflow.com/test"),
					Tags:    &[]string{"tag"},
					Filters: &[]string{"filter"},
				},
			},
			wantType: domain.StackoverflowType,
		},
		{
			name: "LastCheck set correctly",
			args: args{
				userID: 1,
				request: &api.AddLinkRequest{
					Link: aws.String("https://github.com/test"),
				},
			},
			wantType: domain.GithubType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			link, err := mapper.MapAddLinkRequestToDomain(tt.args.userID, tt.args.request)

			require.NoError(t, err)
			require.NotNil(t, link)

			assert.Equal(t, tt.args.userID, link.UserAddID)

			if tt.args.request.Tags != nil {
				assert.Equal(t, *tt.args.request.Tags, link.Tags)
			}

			if tt.args.request.Filters != nil {
				assert.Equal(t, *tt.args.request.Filters, link.Filters)
			}

			assert.Equal(t, *tt.args.request.Link, link.URL)
			assert.Equal(t, tt.wantType, link.Type)

			assert.WithinDuration(t, time.Now(), link.LastCheck, time.Second)
		})
	}
}

func TestMapAddLinkRequestToDomain_Failure(t *testing.T) {
	type args struct {
		userID  int64
		request *api.AddLinkRequest
	}

	tests := []struct {
		name      string
		args      args
		expectErr bool
		errType   error
	}{
		{
			name: "Empty link failure",
			args: args{
				userID:  1,
				request: &api.AddLinkRequest{},
			},
			expectErr: true,
			errType:   &apperrors.LinkValidateError{},
		},
		{
			name: "Unsupported link type failure",
			args: args{
				userID: 1,
				request: &api.AddLinkRequest{
					Link: aws.String("https://test.com"),
				},
			},
			expectErr: true,
			errType:   &apperrors.LinkTypeError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			link, err := mapper.MapAddLinkRequestToDomain(tt.args.userID, tt.args.request)

			require.Error(t, err)
			require.Nil(t, link)
			assert.IsType(t, err, tt.errType)
		})
	}
}
