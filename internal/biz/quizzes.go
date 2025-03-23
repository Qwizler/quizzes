package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"go.opentelemetry.io/otel/trace"
)

type Quiz struct {
	ID          string            `json:"id"`
	UserID      string            `json:"user_id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Duration    *uint64           `json:"duration"`
	Thumbnail   *string           `json:"thumbnail"`
	Cover       *string           `json:"cover"`
	Category    *string           `json:"category"`
	Tags        []string          `json:"tags"`
	Metadata    map[string]string `json:"metadata"`
	CreatedAt   string            `json:"created_at"`
	UpdatedAt   string            `json:"updated_at"`
	DeletedAt   string            `json:"deleted_at"`
	CreatedBy   string            `json:"created_by"`
	UpdatedBy   string            `json:"updated_by"`
	DeletedBy   string            `json:"deleted_by"`
}

type QuizRepo interface {
	Save(ctx context.Context, q *Quiz) (*Quiz, error)
	GetByID(ctx context.Context, id string) (*Quiz, error)
	List(ctx context.Context, pagination *Pagination) ([]*Quiz, error)
	Update(ctx context.Context, q *Quiz) (*Quiz, error)
	Delete(ctx context.Context, id string) (*Quiz, error)
	Search(ctx context.Context, keyword string, pagination *Pagination) ([]*Quiz, error)
}
type QuizUsecase struct {
	repo   QuizRepo
	log    *log.Helper
	tracer trace.Tracer
}

func NewQuizUsecase(repo QuizRepo, logger log.Logger, tracer trace.Tracer) *QuizUsecase {
	return &QuizUsecase{
		repo:   repo,
		log:    log.NewHelper(logger),
		tracer: tracer,
	}
}

func (u *QuizUsecase) CreateQuiz(ctx context.Context, q *Quiz) (*Quiz, error) {
	ctx, span := u.tracer.Start(ctx, "biz.QuizUsecase.CreateQuiz")
	defer span.End()

	res, err := u.repo.Save(ctx, q)
	if err != nil {
		u.log.Warn(err)
		return nil, err
	}
	return res, nil
}

func (u *QuizUsecase) GetQuiz(ctx context.Context, id string) (*Quiz, error) {
	ctx, span := u.tracer.Start(ctx, "biz.QuizUsecase.GetQuiz")
	defer span.End()

	res, err := u.repo.GetByID(ctx, id)
	if err != nil {
		u.log.Warn(err)
		return nil, err
	}
	return res, nil
}

func (u *QuizUsecase) ListQuiz(ctx context.Context, pagination *Pagination) ([]*Quiz, error) {
	ctx, span := u.tracer.Start(ctx, "biz.QuizUsecase.ListQuiz")
	defer span.End()

	res, err := u.repo.List(ctx, pagination)
	if err != nil {
		u.log.Warn(err)
		return nil, err
	}
	return res, nil
}

func (u *QuizUsecase) UpdateQuiz(ctx context.Context, q *Quiz) (*Quiz, error) {
	ctx, span := u.tracer.Start(ctx, "biz.QuizUsecase.UpdateQuiz")
	defer span.End()

	res, err := u.repo.Update(ctx, q)
	if err != nil {
		u.log.Warn(err)
		return nil, err
	}
	return res, nil
}

func (u *QuizUsecase) DeleteQuiz(ctx context.Context, id string) (*Quiz, error) {
	ctx, span := u.tracer.Start(ctx, "biz.QuizUsecase.DeleteQuiz")
	defer span.End()

	res, err := u.repo.Delete(ctx, id)
	if err != nil {
		u.log.Warn(err)
		return nil, err
	}
	return res, nil
}

func (u *QuizUsecase) SearchQuiz(ctx context.Context, keyword string, pagination *Pagination) ([]*Quiz, error) {
	ctx, span := u.tracer.Start(ctx, "biz.QuizUsecase.SearchQuiz")
	defer span.End()

	res, err := u.repo.Search(ctx, keyword, pagination)
	if err != nil {
		u.log.Warn(err)
		return nil, err
	}
	return res, nil
}
