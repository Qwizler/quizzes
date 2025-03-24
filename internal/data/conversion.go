package data

import (
	"github.com/go-kratos/kratos/v2/errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"quiz/internal/biz"
)

func (q *Quiz) QuizToBiz() *biz.Quiz {
	var bizQuiz biz.Quiz

	if !q.ID.IsZero() {
		bizQuiz.ID = q.ID.Hex()
	}
	if q.UserID != "" {
		bizQuiz.UserID = q.UserID
	}
	if q.Title != "" {
		bizQuiz.Title = q.Title
	}
	if q.Description != "" {
		bizQuiz.Description = q.Description
	}
	if q.Duration != nil {
		bizQuiz.Duration = q.Duration
	}
	if q.Thumbnail != nil {
		bizQuiz.Thumbnail = q.Thumbnail
	}
	if q.Cover != nil {
		bizQuiz.Cover = q.Cover
	}
	if q.Category != nil {
		bizQuiz.Category = q.Category
	}
	if q.Tags != nil {
		bizQuiz.Tags = q.Tags
	}
	if q.Metadata != nil {
		bizQuiz.Metadata = q.Metadata
	}
	if q.CreatedBy != "" {
		bizQuiz.CreatedBy = q.CreatedBy
	}
	if q.UpdatedBy != "" {
		bizQuiz.UpdatedBy = q.UpdatedBy
	}
	if q.DeletedBy != "" {
		bizQuiz.DeletedBy = q.DeletedBy
	}
	if q.CreatedAt != "" {
		bizQuiz.CreatedAt = q.CreatedAt
	}
	if q.UpdatedAt != "" {
		bizQuiz.UpdatedAt = q.UpdatedAt
	}
	if q.DeletedAt != "" {
		bizQuiz.DeletedAt = q.DeletedAt
	}
	return &bizQuiz
}

func QuizToData(q *biz.Quiz) (*Quiz, error) {
	var dataQuiz Quiz
	if q.ID != "" {
		oid, err := bson.ObjectIDFromHex(q.ID)
		if err != nil {
			return nil, errors.InternalServer("biz.Quiz.ID was defined as an invalid ObjectID", err.Error())
		}
		dataQuiz.ID = oid
	}
	if q.UserID != "" {
		dataQuiz.UserID = q.UserID
	}
	if q.Title != "" {
		dataQuiz.Title = q.Title
	}
	if q.Description != "" {
		dataQuiz.Description = q.Description
	}
	if q.Duration != nil {
		dataQuiz.Duration = q.Duration
	}
	if q.Thumbnail != nil {
		dataQuiz.Thumbnail = q.Thumbnail
	}
	if q.Cover != nil {
		dataQuiz.Cover = q.Cover
	}
	if q.Category != nil {
		dataQuiz.Category = q.Category
	}
	if q.Tags != nil {
		dataQuiz.Tags = q.Tags
	}
	if q.Metadata != nil {
		dataQuiz.Metadata = q.Metadata
	}
	if q.CreatedBy != "" {
		dataQuiz.CreatedBy = q.CreatedBy
	}
	if q.UpdatedBy != "" {
		dataQuiz.UpdatedBy = q.UpdatedBy
	}
	if q.DeletedBy != "" {
		dataQuiz.DeletedBy = q.DeletedBy
	}
	if q.CreatedAt != "" {
		dataQuiz.CreatedAt = q.CreatedAt
	}
	if q.UpdatedAt != "" {
		dataQuiz.UpdatedAt = q.UpdatedAt
	}
	if q.DeletedAt != "" {
		dataQuiz.DeletedAt = q.DeletedAt
	}
	return &dataQuiz, nil
}

func (q *Question) Biz() *biz.Question {
	var bizQuestion biz.Question

	if !q.ID.IsZero() {
		bizQuestion.ID = q.ID.Hex()
	}
	if q.QuizID != "" {
		bizQuestion.QuizID = q.QuizID
	}
	if q.Question != "" {
		bizQuestion.Question = q.Question
	}
	if q.Hint != "" {
		bizQuestion.Hint = q.Hint
	}
	if q.Difficulty != nil {
		bizQuestion.Difficulty = *q.Difficulty
	}
	if q.Answers != nil {
		var answers []biz.Answer
		for _, a := range q.Answers {
			answer := &biz.Answer{
				Text: a.Text,
			}
			answer.SetIsCorrect(a.IsCorrect)
			if a.Explanation != "" {
				answer.SetExplanation(a.Explanation)
			}
			answers = append(answers, *answer)
		}
		bizQuestion.Answers = answers
	}
	if q.Order != 0 {
		bizQuestion.Order = q.Order
	}
	if q.CreatedBy != "" {
		bizQuestion.CreatedBy = q.CreatedBy
	}
	if q.UpdatedBy != "" {
		bizQuestion.UpdatedBy = q.UpdatedBy
	}
	if q.DeletedBy != "" {
		bizQuestion.DeletedBy = q.DeletedBy
	}
	if q.CreatedAt != "" {
		bizQuestion.CreatedAt = q.CreatedAt
	}
	if q.UpdatedAt != "" {
		bizQuestion.UpdatedAt = q.UpdatedAt
	}
	if q.DeletedAt != "" {
		bizQuestion.DeletedAt = q.DeletedAt
	}
	return &bizQuestion
}
