package application

import (
	"context"

	"github.com/ryunosuke121/muscle-SNS/src/domain"
)

type IPostService interface {
	GetPostsByIds(ctx context.Context, ids []domain.PostID) ([]*PostPublic, error)
	GetUserPosts(ctx context.Context, id domain.UserID) ([]*PostPublic, error)
	CreatePost(ctx context.Context, post *CreatePostRequest) (*PostPublic, error)
	DeletePost(ctx context.Context, userId domain.UserID, postId domain.PostID) error
	GetGroupPosts(ctx context.Context, id domain.UserGroupID) ([]*PostPublic, error)
}

type PostService struct {
	pr domain.IPostRepository
}

func NewPostService(pr domain.IPostRepository) IPostService {
	return &PostService{pr}
}

func (ps *PostService) GetPostById(ctx context.Context, id domain.PostID) (*PostPublic, error) {
	post, err := ps.pr.GetPostsByIds(ctx, []domain.PostID{id})
	if err != nil {
		return nil, err
	}

	resPost := ps.convertPublicPost(post[0])
	return resPost, nil
}

// 投稿を複数件取得する
func (ps *PostService) GetPostsByIds(ctx context.Context, ids []domain.PostID) ([]*PostPublic, error) {
	posts, err := ps.pr.GetPostsByIds(ctx, ids)
	if err != nil {
		return nil, err
	}

	var resPosts []*PostPublic
	for _, post := range posts {
		resPost := &PostPublic{
			ID:        post.ID.String(),
			UserID:    post.UserID.String(),
			Comment:   post.Comment,
			CreatedAt: post.CreatedAt.String(),
			ImageUrl:  post.ImageUrl,
		}

		if post.Training != nil {
			resPost.Training = &TrainingPublic{
				ID:     uint(post.Training.ID),
				UserID: post.Training.UserID.String(),
				Menu: &Menu{
					ID:   uint(post.Training.Menu.ID),
					Name: post.Training.Menu.Name,
				},
				Times:     post.Training.Times,
				Weight:    post.Training.Weight,
				Sets:      post.Training.Sets,
				CreatedAt: post.Training.CreatedAt.String(),
			}
		}

		resPosts = append(resPosts, resPost)
	}

	return resPosts, nil
}

func (ps *PostService) GetUserPosts(ctx context.Context, id domain.UserID) ([]*PostPublic, error) {
	posts, err := ps.pr.GetUserPosts(ctx, id)
	if err != nil {
		return nil, err
	}

	var resPosts []*PostPublic
	for _, post := range posts {
		resPost := ps.convertPublicPost(post)
		resPosts = append(resPosts, resPost)
	}

	return resPosts, nil
}

// 投稿を作成する
func (ps *PostService) CreatePost(ctx context.Context, post *CreatePostRequest) (*PostPublic, error) {
	// トランザクション開始
	fileName, err := ps.pr.SavePostImage(ctx, post.Image)
	if err != nil {
		return nil, err
	}

	newPost := &domain.Post{
		UserID:   post.UserID,
		Comment:  post.Comment,
		ImageUrl: fileName,
		Training: &domain.Training{
			UserID: post.UserID,
			Menu: &domain.Menu{
				ID: post.Training.MenuID,
			},
			Times:  post.Training.Times,
			Weight: post.Training.Weight,
			Sets:   post.Training.Sets,
		},
	}

	newPost, err = ps.pr.CreatePost(ctx, newPost)
	if err != nil {
		return nil, err
	}

	resPost := ps.convertPublicPost(newPost)
	return resPost, nil
}

// 投稿を削除する
func (ps *PostService) DeletePost(ctx context.Context, userId domain.UserID, postId domain.PostID) error {
	post, err := ps.GetPostById(ctx, postId)
	if err != nil {
		return err
	}

	if post.UserID != userId.String() {
		return domain.ErrForbidden
	}

	err = ps.pr.DeletePost(ctx, postId)
	if err != nil {
		return err
	}
	return nil
}

func (ps *PostService) GetGroupPosts(ctx context.Context, id domain.UserGroupID) ([]*PostPublic, error) {
	posts, err := ps.pr.GetGroupPosts(ctx, id)
	if err != nil {
		return nil, err
	}

	var resPosts []*PostPublic
	for _, post := range posts {
		resPost := ps.convertPublicPost(post)
		resPosts = append(resPosts, resPost)
	}

	return resPosts, nil
}

func (ps *PostService) convertPublicPost(post *domain.Post) *PostPublic {
	resPost := &PostPublic{
		ID:        post.ID.String(),
		UserID:    post.UserID.String(),
		Comment:   post.Comment,
		CreatedAt: post.CreatedAt.String(),
		ImageUrl:  post.ImageUrl,
	}

	if post.Training != nil {
		resPost.Training = &TrainingPublic{
			ID:     uint(post.Training.ID),
			UserID: post.Training.UserID.String(),
			Menu: &Menu{
				ID:   uint(post.Training.Menu.ID),
				Name: post.Training.Menu.Name,
			},
			Times:     post.Training.Times,
			Weight:    post.Training.Weight,
			Sets:      post.Training.Sets,
			CreatedAt: post.Training.CreatedAt.String(),
		}
	}

	return resPost
}
