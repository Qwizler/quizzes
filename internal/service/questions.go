package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"go.opentelemetry.io/otel/trace"
	"quiz/internal/biz"

	pb "quiz/api/quizzes/v1"
)

type QuestionsService struct {
	pb.UnimplementedQuestionsServer
	uc     *biz.QuestionsUsecase
	log    *log.Helper
	tracer trace.Tracer
}

func NewQuestionsService(uc *biz.QuestionsUsecase, logger log.Logger, tracer trace.Tracer) *QuestionsService {
	return &QuestionsService{
		uc:     uc,
		log:    log.NewHelper(logger),
		tracer: tracer,
	}
}

func (s *QuestionsService) CreateQuestion(ctx context.Context, req *pb.CreateQuestionRequest) (*pb.CreateQuestionResponse, error) {
	ctx, span := s.tracer.Start(ctx, "service.QuestionsService.CreateQuestion")
	defer span.End()

	question := &biz.Question{
		QuizID:     req.GetQuizId(),
		Question:   req.GetQuestion(),
		Difficulty: uint64(req.GetDifficulty()),
		Order:      float64(req.GetOrder()),
		Hint:       req.GetHint(),
	}

	answers := req.GetAnswers()

	res, err := s.uc.CreateQuestion(ctx, question, answers)
	if err != nil {
		s.log.Warn(err)
		return nil, err
	}
	return &pb.CreateQuestionResponse{
		Id: res.ID,
	}, nil
}
func (s *QuestionsService) GetQuestion(ctx context.Context, req *pb.GetQuestionRequest) (*pb.GetQuestionResponse, error) {
	ctx, span := s.tracer.Start(ctx, "service.QuestionsService.GetQuestion")
	defer span.End()

	question, err := s.uc.GetQuestion(ctx, req.GetQuestionId())
	if err != nil {
		s.log.Warn(err)
		return nil, err
	}
	return &pb.GetQuestionResponse{
		Question: biz.QuestionToPb(question),
	}, nil
}
func (s *QuestionsService) ListQuestion(ctx context.Context, req *pb.ListQuestionRequest) (*pb.ListQuestionResponse, error) {
	ctx, span := s.tracer.Start(ctx, "service.QuestionsService.ListQuestion")
	defer span.End()

	pagination := &biz.Pagination{
		Page: req.GetPagination().GetPage(),
		Size: req.GetPagination().GetPageSize(),
	}

	questions, err := s.uc.ListQuestion(ctx, req.GetQuizId(), pagination)
	if err != nil {
		s.log.Warn(err)
		return nil, err
	}
	var resQuestions []*pb.Question
	for _, q := range questions {
		resQuestions = append(resQuestions, biz.QuestionToPb(q))
	}
	return &pb.ListQuestionResponse{
		Questions: resQuestions,
		Pagination: &pb.Pagination{
			Page:     &pagination.Page,
			PageSize: &pagination.Size,
		},
	}, nil
}
func (s *QuestionsService) UpdateQuestion(ctx context.Context, req *pb.UpdateQuestionRequest) (*pb.UpdateQuestionResponse, error) {
	ctx, span := s.tracer.Start(ctx, "service.QuestionsService.UpdateQuestion")
	defer span.End()

	question := &biz.Question{
		ID:         req.GetQuestionId(),
		QuizID:     req.GetQuizId(),
		Question:   req.GetQuestion(),
		Difficulty: uint64(req.GetDifficulty()),
		Hint:       req.GetHint(),
	}

	res, err := s.uc.UpdateQuestion(ctx, question)
	if err != nil {
		s.log.Warn(err)
		return nil, err
	}
	return &pb.UpdateQuestionResponse{
		Question: biz.QuestionToPb(res),
	}, nil
}
func (s *QuestionsService) DeleteQuestion(ctx context.Context, req *pb.DeleteQuestionRequest) (*pb.DeleteQuestionResponse, error) {
	ctx, span := s.tracer.Start(ctx, "service.QuestionsService.DeleteQuestion")
	defer span.End()

	question, err := s.uc.DeleteQuestion(ctx, req.GetQuestionId())
	if err != nil {
		s.log.Warn(err)
		return nil, err
	}
	return &pb.DeleteQuestionResponse{
		QuestionId: question.ID,
		QuizId:     question.QuizID,
	}, nil
}
func (s *QuestionsService) ReorderQuestion(ctx context.Context, req *pb.ReorderQuestionRequest) (*pb.ReorderQuestionResponse, error) {
	ctx, span := s.tracer.Start(ctx, "service.QuestionsService.ReorderQuestion")
	defer span.End()

	return nil, errors.InternalServer("not implemented", "not implemented")
}
func (s *QuestionsService) ValidateQuestionAnswers(ctx context.Context, req *pb.ValidateQuestionAnswersRequest) (*pb.ValidateQuestionAnswersResponse, error) {
	ctx, span := s.tracer.Start(ctx, "service.QuestionsService.ValidateQuestionAnswers")
	defer span.End()

	res, err := s.uc.ValidateQuestionAnswers(ctx, req.GetQuestionId(), req.GetAnswers())
	if err != nil {
		s.log.Warn(err)
		return nil, err
	}
	return res, nil
}
func (s *QuestionsService) AddAnswer(ctx context.Context, req *pb.AddAnswerRequest) (*pb.AddAnswerResponse, error) {
	ctx, span := s.tracer.Start(ctx, "service.QuestionsService.AddAnswer")
	defer span.End()

	answer, err := s.uc.AddAnswer(ctx, req.GetQuestionId(), *req.GetAnswer())
	if err != nil {
		s.log.Warn(err)
		return nil, err
	}
	return &pb.AddAnswerResponse{
		QuestionId: req.GetQuestionId(),
		QuizId:     req.GetQuizId(),
		Answer:     answer,
	}, nil
}
func (s *QuestionsService) DeleteAnswer(ctx context.Context, req *pb.DeleteAnswerRequest) (*pb.DeleteAnswerResponse, error) {
	ctx, span := s.tracer.Start(ctx, "service.QuestionsService.DeleteAnswer")
	defer span.End()

	res, err := s.uc.DeleteAnswer(ctx, req.GetQuestionId(), req.GetAnswerId())
	if err != nil {
		s.log.Warn(err)
		return nil, err
	}
	return res, nil
}
func (s *QuestionsService) OverrideAnswer(ctx context.Context, req *pb.OverrideAnswerRequest) (*pb.OverrideAnswerResponse, error) {
	ctx, span := s.tracer.Start(ctx, "service.QuestionsService.OverrideAnswer")
	defer span.End()

	res, err := s.uc.OverrideAnswer(ctx, req)
	if err != nil {
		s.log.Warn(err)
		return nil, err
	}
	return res, nil
}
func (s *QuestionsService) PutAnswers(ctx context.Context, req *pb.PutAnswersRequest) (*pb.PutAnswersResponse, error) {
	ctx, span := s.tracer.Start(ctx, "service.QuestionsService.PutAnswers")
	defer span.End()

	answer, err := s.uc.PutAnswers(ctx, req.GetQuestionId(), req.GetAnswers())
	if err != nil {
		s.log.Warn(err)
		return nil, err
	}
	return &pb.PutAnswersResponse{
		QuestionId: req.GetQuestionId(),
		QuizId:     req.GetQuizId(),
		Answers:    answer,
	}, nil
}
func (s *QuestionsService) ReorderAnswers(ctx context.Context, req *pb.ReorderAnswersRequest) (*pb.ReorderAnswersResponse, error) {
	ctx, span := s.tracer.Start(ctx, "service.QuestionsService.ReorderAnswers")
	defer span.End()
	return s.uc.ReorderAnswers(ctx, req)

}
