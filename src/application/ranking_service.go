package application

import (
	"context"
	"sort"

	"github.com/ryunosuke121/muscle-SNS/src/domain"
)

type IRankingService interface {
	GetUserTotalWeightInMonth(ctx context.Context, userId domain.UserID, year int, month int) (*GetUserTotalWeightPublic, error)
	GetMonthRankingInGroup(ctx context.Context, groupId domain.UserGroupID, year int, month int) ([]*GetUserTotalWeightPublic, error)
	GetMonthRankingInGroupByMenu(ctx context.Context, groupId domain.UserGroupID, menuId domain.MenuID, year int, month int) ([]*GetUserTotalWeightPublic, error)
}

type RankingService struct {
	ur domain.IUserRepository
	pr domain.IPostRepository
}

func NewRankingService(ur domain.IUserRepository, pr domain.IPostRepository) IRankingService {
	return &RankingService{ur, pr}
}

// ユーザーの総重量を取得する
func (rs *RankingService) GetUserTotalWeightInMonth(ctx context.Context, userId domain.UserID, year int, month int) (*GetUserTotalWeightPublic, error) {
	userTotalWeight, err := rs.pr.GetUsersTotalWeightInMonth(ctx, []domain.UserID{userId}, year, month)
	if err != nil {
		return nil, err
	}

	user, err := rs.ur.GetUsersByIds(ctx, []domain.UserID{userId})
	if err != nil {
		return nil, err
	}

	resUser := &GetUserTotalWeightPublic{
		ID:          user[0].ID.String(),
		Name:        user[0].Name.String(),
		Email:       user[0].Email,
		AvatarUrl:   user[0].AvatarUrl,
		TotalWeight: userTotalWeight[userId],
	}

	return resUser, nil
}

// グループ内のユーザーの月間ランキングを取得する
func (rs *RankingService) GetMonthRankingInGroup(ctx context.Context, groupId domain.UserGroupID, year int, month int) ([]*GetUserTotalWeightPublic, error) {
	users, err := rs.ur.GetUsersInGroup(ctx, groupId)
	if err != nil {
		return nil, err
	}

	userIds := make([]domain.UserID, len(users))
	for i, user := range users {
		userIds[i] = user.ID
	}

	totalWeightMap, err := rs.pr.GetUsersTotalWeightInMonth(ctx, userIds, year, month)
	if err != nil {
		return nil, err
	}

	resUsers := make([]*GetUserTotalWeightPublic, len(users))
	for i, user := range users {
		resUser := &GetUserTotalWeightPublic{
			ID:          user.ID.String(),
			Name:        user.Name.String(),
			Email:       user.Email,
			AvatarUrl:   user.AvatarUrl,
			TotalWeight: totalWeightMap[user.ID],
		}
		resUsers[i] = resUser
	}

	resUsers = sortUserByTotalWeights(resUsers)

	return resUsers, nil
}

func (rs *RankingService) GetMonthRankingInGroupByMenu(ctx context.Context, groupId domain.UserGroupID, menuId domain.MenuID, year int, month int) ([]*GetUserTotalWeightPublic, error) {
	users, err := rs.ur.GetUsersInGroup(ctx, groupId)
	if err != nil {
		return nil, err
	}

	userIds := make([]domain.UserID, len(users))
	for i, user := range users {
		userIds[i] = user.ID
	}

	totalWeightMap, err := rs.pr.GetUsersTotalWeightByMenuInMonth(ctx, userIds, menuId, year, month)
	if err != nil {
		return nil, err
	}

	resUsers := make([]*GetUserTotalWeightPublic, len(users))
	for i, user := range users {
		resUser := &GetUserTotalWeightPublic{
			ID:          user.ID.String(),
			Name:        user.Name.String(),
			Email:       user.Email,
			AvatarUrl:   user.AvatarUrl,
			TotalWeight: totalWeightMap[user.ID],
		}
		resUsers[i] = resUser
	}

	resUsers = sortUserByTotalWeights(resUsers)

	return resUsers, nil
}

type ByTotalWeight []*GetUserTotalWeightPublic

func (a ByTotalWeight) Len() int           { return len(a) }
func (a ByTotalWeight) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTotalWeight) Less(i, j int) bool { return a[i].TotalWeight < a[j].TotalWeight }

func sortUserByTotalWeights(users []*GetUserTotalWeightPublic) []*GetUserTotalWeightPublic {
	sort.Sort(sort.Reverse(ByTotalWeight(users)))
	return users
}
