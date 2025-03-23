package biz

import pb "quiz/api/quizzes/v1"

func QuizToPb(q *Quiz) *pb.Quiz {
	var quiz pb.Quiz
	if q.ID != "" {
		quiz.Id = q.ID
	}
	if q.UserID != "" {
		quiz.UserId = q.UserID
	}
	if q.Title != "" {
		quiz.Title = q.Title
	}
	if q.Description != "" {
		quiz.Description = q.Description
	}
	if q.Duration != nil {
		quiz.Duration = q.Duration
	}
	if q.Thumbnail != nil {
		quiz.Thumbnail = q.Thumbnail
	}
	if q.Cover != nil {
		quiz.Cover = q.Cover
	}
	if q.Category != nil {
		quiz.Category = q.Category
	}
	if q.Tags != nil {
		quiz.Tags = q.Tags
	}
	if q.Metadata != nil {
		quiz.Metadata = q.Metadata
	}
	var audit pb.Audit
	if q.CreatedBy != "" {
		audit.CreatedBy = &q.CreatedBy
	}
	if q.UpdatedBy != "" {
		audit.UpdatedBy = &q.UpdatedBy
	}
	if q.DeletedBy != "" {
		audit.DeletedBy = &q.DeletedBy
	}
	if q.CreatedAt != "" {
		audit.CreatedAt = q.CreatedAt
	}
	if q.UpdatedAt != "" {
		audit.UpdatedAt = q.UpdatedAt
	}
	if q.DeletedAt != "" {
		audit.DeletedAt = q.DeletedAt
	}
	quiz.Audit = &audit

	return &quiz
}
