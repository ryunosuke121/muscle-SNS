package application

import (
	"context"
	"fmt"
	"testing"

	"github.com/ryunosuke121/muscle-SNS/src/domain"
	mock_domain "github.com/ryunosuke121/muscle-SNS/src/mocks/domain"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestRankingService_GetUserTotalWeightInMonth(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx    context.Context
		userId domain.UserID
		year   int
		month  int
	}
	tests := []struct {
		name    string
		args    args
		mockFn  func(ur *mock_domain.MockIUserRepository, pr *mock_domain.MockIPostRepository)
		want    *GetUserTotalWeightPublic
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				ctx:    context.Background(),
				userId: domain.UserID("test-user-id"),
				year:   2021,
				month:  1,
			},
			mockFn: func(ur *mock_domain.MockIUserRepository, pr *mock_domain.MockIPostRepository) {
				pr.EXPECT().GetUsersTotalWeightInMonth(gomock.Any(), gomock.Eq([]domain.UserID{domain.UserID("test-user-id")}), gomock.Eq(2021), gomock.Eq(1)).Return(
					map[domain.UserID]uint{domain.UserID("test-user-id"): 100}, nil)

				ur.EXPECT().GetUsersByIds(gomock.Any(), []domain.UserID{domain.UserID("test-user-id")}).Return(
					[]*domain.User{{
						ID:        domain.UserID("test-user-id"),
						Name:      domain.UserName("test-user-name"),
						Email:     "test-user-email",
						AvatarUrl: "test-user-avatar-url",
					}}, nil)
			},
			want: &GetUserTotalWeightPublic{
				ID:          "test-user-id",
				Name:        "test-user-name",
				Email:       "test-user-email",
				AvatarUrl:   "test-user-avatar-url",
				TotalWeight: 100,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ur := mock_domain.NewMockIUserRepository(ctrl)
			pr := mock_domain.NewMockIPostRepository(ctrl)
			tt.mockFn(ur, pr)

			rs := NewRankingService(ur, pr)

			got, err := rs.GetUserTotalWeightInMonth(tt.args.ctx, tt.args.userId, tt.args.year, tt.args.month)
			if (err != nil) != tt.wantErr {
				t.Errorf("RankingService.GetUserTotalWeightInMonth() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRankingService_GetMonthRankingInGroup(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx     context.Context
		groupId domain.UserGroupID
		year    int
		month   int
	}
	tests := []struct {
		name    string
		args    args
		mockFn  func(ur *mock_domain.MockIUserRepository, pr *mock_domain.MockIPostRepository)
		want    []*GetUserTotalWeightPublic
		wantErr bool
	}{
		{
			name: "success with getting two users",
			args: args{
				ctx:     context.Background(),
				groupId: domain.UserGroupID(3),
				year:    2021,
				month:   1,
			},
			mockFn: func(ur *mock_domain.MockIUserRepository, pr *mock_domain.MockIPostRepository) {
				ur.EXPECT().GetUsersInGroup(gomock.Any(), domain.UserGroupID(3)).Return(
					[]*domain.User{
						{
							ID:        domain.UserID("test-user-id"),
							Name:      domain.UserName("test-user-name"),
							Email:     "test-user-email",
							AvatarUrl: "test-user-avatar-url",
						},
						{
							ID:        domain.UserID("test-user-id2"),
							Name:      domain.UserName("test-user-name2"),
							Email:     "test-user-email2",
							AvatarUrl: "test-user-avatar-url2",
						},
					}, nil)
				pr.EXPECT().GetUsersTotalWeightInMonth(gomock.Any(), gomock.Eq([]domain.UserID{"test-user-id", "test-user-id2"}), gomock.Eq(2021), gomock.Eq(1)).Return(
					map[domain.UserID]uint{
						domain.UserID("test-user-id"):  200,
						domain.UserID("test-user-id2"): 100,
					}, nil)
			},
			want: []*GetUserTotalWeightPublic{
				{
					ID:          "test-user-id",
					Name:        "test-user-name",
					Email:       "test-user-email",
					AvatarUrl:   "test-user-avatar-url",
					TotalWeight: 200,
				},
				{
					ID:          "test-user-id2",
					Name:        "test-user-name2",
					Email:       "test-user-email2",
					AvatarUrl:   "test-user-avatar-url2",
					TotalWeight: 100,
				},
			},
		},
		{
			name: "is sorted test",
			args: args{
				ctx:     context.Background(),
				groupId: domain.UserGroupID(3),
				year:    2021,
				month:   1,
			},
			mockFn: func(ur *mock_domain.MockIUserRepository, pr *mock_domain.MockIPostRepository) {
				ur.EXPECT().GetUsersInGroup(gomock.Any(), domain.UserGroupID(3)).Return(
					[]*domain.User{
						{
							ID:        domain.UserID("test-user-id"),
							Name:      domain.UserName("test-user-name"),
							Email:     "test-user-email",
							AvatarUrl: "test-user-avatar-url",
						},
						{
							ID:        domain.UserID("test-user-id2"),
							Name:      domain.UserName("test-user-name2"),
							Email:     "test-user-email2",
							AvatarUrl: "test-user-avatar-url2",
						},
						{
							ID:        domain.UserID("test-user-id3"),
							Name:      domain.UserName("test-user-name3"),
							Email:     "test-user-email3",
							AvatarUrl: "test-user-avatar-url3",
						},
					}, nil)
				pr.EXPECT().GetUsersTotalWeightInMonth(gomock.Any(), gomock.Eq([]domain.UserID{"test-user-id", "test-user-id2", "test-user-id3"}), gomock.Eq(2021), gomock.Eq(1)).Return(
					map[domain.UserID]uint{
						domain.UserID("test-user-id"):  200,
						domain.UserID("test-user-id2"): 100,
						domain.UserID("test-user-id3"): 500,
					}, nil)
			},
			want: []*GetUserTotalWeightPublic{
				{
					ID:          "test-user-id3",
					Name:        "test-user-name3",
					Email:       "test-user-email3",
					AvatarUrl:   "test-user-avatar-url3",
					TotalWeight: 500,
				},
				{
					ID:          "test-user-id",
					Name:        "test-user-name",
					Email:       "test-user-email",
					AvatarUrl:   "test-user-avatar-url",
					TotalWeight: 200,
				},
				{
					ID:          "test-user-id2",
					Name:        "test-user-name2",
					Email:       "test-user-email2",
					AvatarUrl:   "test-user-avatar-url2",
					TotalWeight: 100,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ur := mock_domain.NewMockIUserRepository(ctrl)
			pr := mock_domain.NewMockIPostRepository(ctrl)
			tt.mockFn(ur, pr)

			rs := NewRankingService(ur, pr)
			got, err := rs.GetMonthRankingInGroup(tt.args.ctx, tt.args.groupId, tt.args.year, tt.args.month)
			fmt.Printf("%+v\n", got)
			if (err != nil) != tt.wantErr {
				t.Errorf("RankingService.GetMonthRankingInGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRankingService_GetMonthRankingInGroupByMenu(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx     context.Context
		groupId domain.UserGroupID
		menuId  domain.MenuID
		year    int
		month   int
	}
	tests := []struct {
		name    string
		args    args
		mockFn  func(ur *mock_domain.MockIUserRepository, pr *mock_domain.MockIPostRepository)
		want    []*GetUserTotalWeightPublic
		wantErr bool
	}{
		{
			name: "success with two users",
			args: args{
				ctx:     context.Background(),
				groupId: domain.UserGroupID(3),
				menuId:  domain.MenuID(1),
				year:    2021,
				month:   1,
			},
			mockFn: func(ur *mock_domain.MockIUserRepository, pr *mock_domain.MockIPostRepository) {
				ur.EXPECT().GetUsersInGroup(gomock.Any(), domain.UserGroupID(3)).Return(
					[]*domain.User{
						{
							ID:        domain.UserID("test-user-id"),
							Name:      domain.UserName("test-user-name"),
							Email:     "test-user-email",
							AvatarUrl: "test-user-avatar-url",
						},
						{
							ID:        domain.UserID("test-user-id2"),
							Name:      domain.UserName("test-user-name2"),
							Email:     "test-user-email2",
							AvatarUrl: "test-user-avatar-url2",
						},
					}, nil)
				pr.EXPECT().GetUsersTotalWeightByMenuInMonth(gomock.Any(), gomock.Eq([]domain.UserID{"test-user-id", "test-user-id2"}), gomock.Eq(domain.MenuID(1)), gomock.Eq(2021), gomock.Eq(1)).Return(
					map[domain.UserID]uint{
						domain.UserID("test-user-id"):  200,
						domain.UserID("test-user-id2"): 100,
					}, nil)
			},
			want: []*GetUserTotalWeightPublic{
				{
					ID:          "test-user-id",
					Name:        "test-user-name",
					Email:       "test-user-email",
					AvatarUrl:   "test-user-avatar-url",
					TotalWeight: 200,
				},
				{
					ID:          "test-user-id2",
					Name:        "test-user-name2",
					Email:       "test-user-email2",
					AvatarUrl:   "test-user-avatar-url2",
					TotalWeight: 100,
				},
			},
		},
		{
			name: "is sorted test",
			args: args{
				ctx:     context.Background(),
				groupId: domain.UserGroupID(3),
				menuId:  domain.MenuID(1),
				year:    2021,
				month:   1,
			},
			mockFn: func(ur *mock_domain.MockIUserRepository, pr *mock_domain.MockIPostRepository) {
				ur.EXPECT().GetUsersInGroup(gomock.Any(), domain.UserGroupID(3)).Return(
					[]*domain.User{
						{
							ID:        domain.UserID("test-user-id"),
							Name:      domain.UserName("test-user-name"),
							Email:     "test-user-email",
							AvatarUrl: "test-user-avatar-url",
						},
						{
							ID:        domain.UserID("test-user-id2"),
							Name:      domain.UserName("test-user-name2"),
							Email:     "test-user-email2",
							AvatarUrl: "test-user-avatar-url2",
						},
						{
							ID:        domain.UserID("test-user-id3"),
							Name:      domain.UserName("test-user-name3"),
							Email:     "test-user-email3",
							AvatarUrl: "test-user-avatar-url3",
						},
					}, nil)
				pr.EXPECT().GetUsersTotalWeightByMenuInMonth(gomock.Any(), gomock.Eq([]domain.UserID{"test-user-id", "test-user-id2", "test-user-id3"}), gomock.Eq(domain.MenuID(1)), gomock.Eq(2021), gomock.Eq(1)).Return(
					map[domain.UserID]uint{
						domain.UserID("test-user-id"):  200,
						domain.UserID("test-user-id2"): 100,
						domain.UserID("test-user-id3"): 500,
					}, nil)
			},
			want: []*GetUserTotalWeightPublic{
				{
					ID:          "test-user-id3",
					Name:        "test-user-name3",
					Email:       "test-user-email3",
					AvatarUrl:   "test-user-avatar-url3",
					TotalWeight: 500,
				},
				{
					ID:          "test-user-id",
					Name:        "test-user-name",
					Email:       "test-user-email",
					AvatarUrl:   "test-user-avatar-url",
					TotalWeight: 200,
				},
				{
					ID:          "test-user-id2",
					Name:        "test-user-name2",
					Email:       "test-user-email2",
					AvatarUrl:   "test-user-avatar-url2",
					TotalWeight: 100,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ur := mock_domain.NewMockIUserRepository(ctrl)
			pr := mock_domain.NewMockIPostRepository(ctrl)
			tt.mockFn(ur, pr)

			rs := NewRankingService(ur, pr)
			got, err := rs.GetMonthRankingInGroupByMenu(tt.args.ctx, tt.args.groupId, tt.args.menuId, tt.args.year, tt.args.month)
			if (err != nil) != tt.wantErr {
				t.Errorf("RankingService.GetMonthRankingInGroupByMenu() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_sortUserByTotalWeights(t *testing.T) {
	type args struct {
		users []*GetUserTotalWeightPublic
	}
	tests := []struct {
		name string
		args args
		want []*GetUserTotalWeightPublic
	}{
		{
			name: "is sorted test",
			args: args{
				users: []*GetUserTotalWeightPublic{
					{
						ID:          "test-user-id",
						Name:        "test-user-name",
						Email:       "test-user-email",
						AvatarUrl:   "test-user-avatar-url",
						TotalWeight: 200,
					},
					{
						ID:          "test-user-id2",
						Name:        "test-user-name2",
						Email:       "test-user-email2",
						AvatarUrl:   "test-user-avatar-url2",
						TotalWeight: 100,
					},
					{
						ID:          "test-user-id3",
						Name:        "test-user-name3",
						Email:       "test-user-email3",
						AvatarUrl:   "test-user-avatar-url3",
						TotalWeight: 500,
					},
					{
						ID:          "test-user-id4",
						Name:        "test-user-name4",
						Email:       "test-user-email4",
						AvatarUrl:   "test-user-avatar-url4",
						TotalWeight: 300,
					},
					{
						ID:          "test-user-id5",
						Name:        "test-user-name5",
						Email:       "test-user-email5",
						AvatarUrl:   "test-user-avatar-url5",
						TotalWeight: 500,
					},
				},
			},
			want: []*GetUserTotalWeightPublic{
				{
					ID:          "test-user-id3",
					Name:        "test-user-name3",
					Email:       "test-user-email3",
					AvatarUrl:   "test-user-avatar-url3",
					TotalWeight: 500,
				},
				{
					ID:          "test-user-id5",
					Name:        "test-user-name5",
					Email:       "test-user-email5",
					AvatarUrl:   "test-user-avatar-url5",
					TotalWeight: 500,
				},
				{
					ID:          "test-user-id4",
					Name:        "test-user-name4",
					Email:       "test-user-email4",
					AvatarUrl:   "test-user-avatar-url4",
					TotalWeight: 300,
				},
				{
					ID:          "test-user-id",
					Name:        "test-user-name",
					Email:       "test-user-email",
					AvatarUrl:   "test-user-avatar-url",
					TotalWeight: 200,
				},
				{
					ID:          "test-user-id2",
					Name:        "test-user-name2",
					Email:       "test-user-email2",
					AvatarUrl:   "test-user-avatar-url2",
					TotalWeight: 100,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sortUserByTotalWeights(tt.args.users)
			assert.Equal(t, tt.want, got)
		})
	}
}
