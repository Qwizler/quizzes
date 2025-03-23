package service

import (
	"context"

	pb "quiz/api/quizzes/v1"
)

type QuizQuestionsService struct {
	pb.UnimplementedQuizQuestionsServer
}

func NewQuizQuestionsService() *QuizQuestionsService {
	return &QuizQuestionsService{}
}

func (s *QuizQuestionsService) AddQuestion(ctx context.Context, req *pb.AddQuestionRequest) (*pb.AddQuestionResponse, error) {
	return &pb.AddQuestionResponse{}, nil
}
func (s *QuizQuestionsService) UpdateQuestion(ctx context.Context, req *pb.UpdateQuestionRequest) (*pb.UpdateQuestionResponse, error) {
	return &pb.UpdateQuestionResponse{}, nil
}
func (s *QuizQuestionsService) DeleteQuestion(ctx context.Context, req *pb.DeleteQuestionRequest) (*pb.DeleteQuestionResponse, error) {
	return &pb.DeleteQuestionResponse{}, nil
}
func (s *QuizQuestionsService) ListQuestion(ctx context.Context, req *pb.ListQuestionRequest) (*pb.ListQuestionResponse, error) {
	return &pb.ListQuestionResponse{}, nil
}
