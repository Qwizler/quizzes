package service

import (
	"context"

	pb "quiz/api/quizzes/v1"
)

type QuizAnswersService struct {
	pb.UnimplementedQuizAnswersServer
}

func NewQuizAnswersService() *QuizAnswersService {
	return &QuizAnswersService{}
}

func (s *QuizAnswersService) ValidateQuestionAnswers(ctx context.Context, req *pb.ValidateQuestionAnswersRequest) (*pb.ValidateQuestionAnswersResponse, error) {
	return &pb.ValidateQuestionAnswersResponse{}, nil
}
func (s *QuizAnswersService) CreateChoice(ctx context.Context, req *pb.CreateChoiceRequest) (*pb.CreateChoiceResponse, error) {
	return &pb.CreateChoiceResponse{}, nil
}
func (s *QuizAnswersService) UpdateChoice(ctx context.Context, req *pb.UpdateChoiceRequest) (*pb.UpdateChoiceResponse, error) {
	return &pb.UpdateChoiceResponse{}, nil
}
func (s *QuizAnswersService) DeleteChoice(ctx context.Context, req *pb.DeleteChoiceRequest) (*pb.DeleteChoiceResponse, error) {
	return &pb.DeleteChoiceResponse{}, nil
}
func (s *QuizAnswersService) ListChoice(ctx context.Context, req *pb.ListChoiceRequest) (*pb.ListChoiceResponse, error) {
	return &pb.ListChoiceResponse{}, nil
}
