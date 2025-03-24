package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.opentelemetry.io/otel/trace"
	pb "quiz/api/quizzes/v1"
)

type Answer struct {
	ID          string `json:"id"`
	Text        string `json:"text"`
	isCorrect   bool
	explanation string
}

func (a *Answer) IsCorrect() bool {
	return a.isCorrect
}

func (a *Answer) Explanation() string {
	return a.explanation
}

func (a *Answer) SetIsCorrect(isCorrect bool) *Answer {
	a.isCorrect = isCorrect
	return a
}
func (a *Answer) SetExplanation(explanation string) *Answer {
	a.explanation = explanation
	return a
}

type Question struct {
	ID         string `json:"id"`
	QuizID     string `json:"quiz_id"`
	Question   string `json:"question"`
	Difficulty uint64 `json:"difficulty"`
	Order      float64
	Answers    []Answer
	Hint       string `json:"hint"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
	DeletedAt  string `json:"deleted_at"`
	CreatedBy  string `json:"created_by"`
	UpdatedBy  string `json:"updated_by"`
	DeletedBy  string `json:"deleted_by"`
}

type ReorderPayload struct {
	MakeFirst bool
	MakeLast  bool
	AboveID   string
	BelowID   string
}

type QuestionsRepo interface {
	Save(ctx context.Context, q *Question) (*Question, error)
	GetByID(ctx context.Context, id string) (*Question, error)
	List(ctx context.Context, quizID string, pagination *Pagination) ([]*Question, error)
	Update(ctx context.Context, q *Question) (*Question, error)
	Delete(ctx context.Context, id string) (*Question, error)
}

type QuestionsUsecase struct {
	repo   QuestionsRepo
	log    *log.Helper
	tracer trace.Tracer
}

func NewQuestionUsecase(repo QuestionsRepo, logger log.Logger, tracer trace.Tracer) *QuestionsUsecase {
	return &QuestionsUsecase{
		repo:   repo,
		log:    log.NewHelper(logger),
		tracer: tracer,
	}
}
func (u *QuestionsUsecase) CreateQuestion(ctx context.Context, q *Question, _answers []*pb.AnswerCreation) (*Question, error) {
	ctx, span := u.tracer.Start(ctx, "biz.QuestionsUsecase.CreateQuestion")
	defer span.End()

	answers := make([]Answer, 0, len(_answers))
	for _, a := range _answers {
		answers = append(answers, Answer{
			ID:          uuid.New().String(),
			Text:        a.Text,
			isCorrect:   a.IsCorrect,
			explanation: *a.Explanation,
		})
	}

	q.Answers = answers

	res, err := u.repo.Save(ctx, q)
	if err != nil {
		u.log.Warn(err)
		return nil, err
	}
	return res, nil
}

func (u *QuestionsUsecase) GetQuestion(ctx context.Context, id string) (*Question, error) {
	ctx, span := u.tracer.Start(ctx, "biz.QuestionsUsecase.GetQuestion")
	defer span.End()

	res, err := u.repo.GetByID(ctx, id)
	if err != nil {
		u.log.Warn(err)
		return nil, err
	}
	return res, nil
}

func (u *QuestionsUsecase) ListQuestion(ctx context.Context, quizID string, pagination *Pagination) ([]*Question, error) {
	ctx, span := u.tracer.Start(ctx, "biz.QuestionsUsecase.ListQuestion")
	defer span.End()

	res, err := u.repo.List(ctx, quizID, pagination)
	if err != nil {
		u.log.Warn(err)
		return nil, err
	}
	return res, nil
}

func (u *QuestionsUsecase) UpdateQuestion(ctx context.Context, q *Question) (*Question, error) {
	ctx, span := u.tracer.Start(ctx, "biz.QuestionsUsecase.UpdateQuestion")
	defer span.End()

	res, err := u.repo.Update(ctx, q)
	if err != nil {
		u.log.Warn(err)
		return nil, err
	}
	return res, nil
}

func (u *QuestionsUsecase) DeleteQuestion(ctx context.Context, id string) (*Question, error) {
	ctx, span := u.tracer.Start(ctx, "biz.QuestionsUsecase.DeleteQuestion")
	defer span.End()

	res, err := u.repo.Delete(ctx, id)
	if err != nil {
		u.log.Warn(err)
		return nil, err
	}
	return res, nil
}

func (u *QuestionsUsecase) ValidateQuestionAnswers(ctx context.Context, questionID string, uAnswers []*pb.UserAnswer) (*pb.ValidateQuestionAnswersResponse, error) {
	ctx, span := u.tracer.Start(ctx, "biz.QuestionsUsecase.ValidateQuestionAnswers")
	defer span.End()

	q, err := u.repo.GetByID(ctx, questionID)
	if err != nil {
		u.log.Warn(err)
		return nil, err
	}
	answers := q.Answers
	if len(answers) == 0 {
		return nil, errors.BadRequest("Invalid question", "question has no answers")
	}
	results, score, err := validateAnswers(q.Answers, uAnswers)
	if err != nil {
		u.log.Warn(err)
		return nil, err
	}

	return &pb.ValidateQuestionAnswersResponse{
		QuestionId: questionID,
		Results:    results,
		Score:      score,
	}, nil
}

func (u *QuestionsUsecase) ReorderQuestion(ctx context.Context, questionID string, payload ReorderPayload) (*Question, error) {
	ctx, span := u.tracer.Start(ctx, "biz.QuestionsUsecase.ReorderQuestion")
	defer span.End()

	return nil, errors.InternalServer("not implemented", "not implemented")
}

func (u *QuestionsUsecase) AddAnswer(ctx context.Context, questionID string, answer pb.AnswerCreation) (*pb.Answer, error) {
	ctx, span := u.tracer.Start(ctx, "biz.QuestionsUsecase.AddAnswer")
	defer span.End()
	q, err := u.repo.GetByID(ctx, questionID)
	if err != nil {
		u.log.Warn(err)
		return nil, err
	}

	newAnswer := &Answer{
		ID:          bson.NewObjectID().Hex(),
		Text:        answer.Text,
		isCorrect:   answer.IsCorrect,
		explanation: *answer.Explanation,
	}
	newAnswers := append(q.Answers, *newAnswer)
	q.Answers = newAnswers

	res, err := u.repo.Update(ctx, q)
	if err != nil {
		u.log.Warn(err)
		return nil, err
	}
	var resAnswers []*pb.Answer
	for _, a := range res.Answers {
		resAnswers = append(resAnswers, &pb.Answer{
			Id:          a.ID,
			Text:        a.Text,
			IsCorrect:   a.isCorrect,
			Explanation: &a.explanation,
		})
	}

	return &pb.Answer{
		Id:          newAnswer.ID,
		Text:        newAnswer.Text,
		IsCorrect:   newAnswer.isCorrect,
		Explanation: &newAnswer.explanation,
	}, nil
}

func (u *QuestionsUsecase) DeleteAnswer(ctx context.Context, QuestionID string, AnswerID string) (*pb.DeleteAnswerResponse, error) {
	ctx, span := u.tracer.Start(ctx, "biz.QuestionsUsecase.DeleteAnswer")
	defer span.End()

	q, err := u.repo.GetByID(ctx, QuestionID)
	if err != nil {
		u.log.Warn(err)
		return nil, err
	}
	answers := q.Answers
	if len(answers) == 0 {
		return nil, errors.BadRequest("Invalid question", "question has no answers")
	}
	newAnswers := make([]Answer, 0, len(answers)-1)
	for _, a := range answers {
		if a.ID != AnswerID {
			newAnswers = append(newAnswers, a)
		}
	}
	q.Answers = newAnswers

	res, err := u.repo.Update(ctx, q)
	if err != nil {
		u.log.Warn(err)
		return nil, err
	}
	var resAnswers []*pb.Answer
	for _, a := range res.Answers {
		resAnswers = append(resAnswers, &pb.Answer{
			Id:          a.ID,
			Text:        a.Text,
			IsCorrect:   a.isCorrect,
			Explanation: &a.explanation,
		})
	}

	return &pb.DeleteAnswerResponse{
		QuestionId: res.ID,
		QuizId:     res.QuizID,
		AnswerId:   AnswerID,
	}, nil

}

func (u *QuestionsUsecase) OverrideAnswer(ctx context.Context, request *pb.OverrideAnswerRequest) (*pb.OverrideAnswerResponse, error) {
	ctx, span := u.tracer.Start(ctx, "biz.QuestionsUsecase.OverrideAnswer")
	defer span.End()

	q, err := u.repo.GetByID(ctx, request.GetQuestionId())
	if err != nil {
		u.log.Warn(err)
		return nil, err
	}
	answers := q.Answers
	if len(answers) == 0 {
		return nil, errors.BadRequest("Invalid question", "question has no answers")
	}
	if request.GetAnswerId() == "" {
		return nil, errors.BadRequest("Invalid answer ID", "answer ID is empty")
	}
	var target *Answer
	for _, a := range answers {
		if a.ID == request.GetAnswerId() {
			target = &a
		}
	}
	if target == nil {
		return nil, errors.BadRequest("Invalid target ID", "target ID is empty")
	}
	target.Text = request.GetAnswer().GetText()
	target.isCorrect = request.GetAnswer().GetIsCorrect()
	target.explanation = request.GetAnswer().GetExplanation()

	res, err := u.repo.Update(ctx, q)
	if err != nil {
		u.log.Warn(err)
		return nil, err
	}
	var resAnswers []*pb.Answer
	for _, a := range res.Answers {
		resAnswers = append(resAnswers, &pb.Answer{
			Id:          a.ID,
			Text:        a.Text,
			IsCorrect:   a.isCorrect,
			Explanation: &a.explanation,
		})
	}

	return &pb.OverrideAnswerResponse{
		QuestionId: res.ID,
		QuizId:     res.QuizID,
		Answer:     &pb.Answer{Id: request.GetAnswerId(), Text: target.Text, IsCorrect: target.isCorrect, Explanation: &target.explanation},
	}, nil
}

func (u *QuestionsUsecase) PutAnswers(ctx context.Context, QuestionID string, answers []*pb.AnswerCreation) ([]*pb.Answer, error) {
	ctx, span := u.tracer.Start(ctx, "biz.QuestionsUsecase.PutAnswers")
	defer span.End()

	q, err := u.repo.GetByID(ctx, QuestionID)
	if err != nil {
		u.log.Warn(err)
		return nil, err
	}
	newAnswers := make([]Answer, 0, len(answers))
	for _, a := range answers {
		newID := bson.NewObjectID().Hex()
		newAnswers = append(newAnswers, Answer{
			ID:          newID,
			Text:        a.Text,
			isCorrect:   a.IsCorrect,
			explanation: *a.Explanation,
		})
	}

	q.Answers = newAnswers

	res, err := u.repo.Update(ctx, q)
	if err != nil {
		u.log.Warn(err)
		return nil, err
	}
	var resAnswers []*pb.Answer
	for _, a := range res.Answers {
		resAnswers = append(resAnswers, &pb.Answer{
			Id:          a.ID,
			Text:        a.Text,
			IsCorrect:   a.isCorrect,
			Explanation: &a.explanation,
		})
	}

	return resAnswers, nil
}

func (u *QuestionsUsecase) ReorderAnswers(ctx context.Context, response *pb.ReorderAnswersRequest) (*pb.ReorderAnswersResponse, error) {
	ctx, span := u.tracer.Start(ctx, "biz.QuestionsUsecase.ReorderAnswers")
	defer span.End()

	q, err := u.repo.GetByID(ctx, response.GetQuestionId())
	if err != nil {
		u.log.Warn(err)
		return nil, err
	}
	answers := q.Answers
	if len(answers) != len(response.GetAnswerIds()) {
		return nil, errors.BadRequest("Invalid answer IDs", "answer IDs are not equal to number of answers")
	}

	newOrder := make([]Answer, 0, len(response.GetAnswerIds()))
	oldAnswersMap := make(map[string]Answer, len(answers))
	for _, a := range answers {
		oldAnswersMap[a.ID] = a
	}
	for _, aID := range response.GetAnswerIds() {
		if a, exists := oldAnswersMap[aID]; exists {
			newOrder = append(newOrder, a)
		}
	}
	q.Answers = newOrder
	res, err := u.repo.Update(ctx, q)
	if err != nil {
		u.log.Warn(err)
		return nil, err
	}
	pbAnswers := make([]*pb.Answer, 0, len(res.Answers))
	for _, a := range res.Answers {
		pbAnswers = append(pbAnswers, &pb.Answer{
			Id:          a.ID,
			Text:        a.Text,
			IsCorrect:   a.isCorrect,
			Explanation: &a.explanation,
		})
	}
	return &pb.ReorderAnswersResponse{
		QuestionId: response.GetQuestionId(),
		Answers:    pbAnswers,
	}, nil

}

// ReorderAnswers NOTE: DEPRECATED
//func (u *QuestionsUsecase) ReorderAnswers(ctx context.Context, QuestionID string, AnswerID string, payload ReorderPayload) (*pb.ReorderAnswersResponse, error) {
//	ctx, span := u.tracer.Start(ctx, "biz.QuestionsUsecase.ReorderAnswers")
//	defer span.End()
//	q, err := u.repo.GetByID(ctx, QuestionID)
//	if err != nil {
//		u.log.Warn(err)
//		return nil, err
//	}
//	answers := q.Answers
//	if len(answers) == 0 {
//		return nil, errors.BadRequest("Invalid question", "question has no answers")
//	}
//	if AnswerID == "" {
//		return nil, errors.BadRequest("Invalid answer ID", "answer ID is empty")
//	}
//
//	newAnswers, err := ReorderAnswers(answers, payload, AnswerID)
//	if err != nil {
//		u.log.Warn(err)
//		return nil, err
//	}
//	q.Answers = newAnswers
//
//	res, err := u.repo.Update(ctx, q)
//	if err != nil {
//		u.log.Warn(err)
//		return nil, err
//	}
//	var resAnswers []*pb.Answer
//	for _, a := range res.Answers {
//		resAnswers = append(resAnswers, &pb.Answer{
//			Id:          a.ID,
//			Text:        a.Text,
//			IsCorrect:   a.isCorrect,
//			Explanation: &a.explanation,
//		})
//	}
//
//	return &pb.ReorderAnswersResponse{
//		QuestionId: QuestionID,
//		Answers:    resAnswers,
//	}, nil
//}

func validateAnswers(answers []Answer, userAnswers []*pb.UserAnswer) ([]*pb.AnswerResult, float32, error) {
	if len(answers) != len(userAnswers) {
		return nil, 0, errors.BadRequest("invalid user answers", "number of user answers is not equal to number of answers")
	}

	// map for faster lookups
	userAnswerMap := make(map[string]*pb.UserAnswer, len(userAnswers))
	for _, ua := range userAnswers {
		userAnswerMap[ua.AnswerId] = ua
	}

	results := make([]*pb.AnswerResult, 0, len(answers))
	correctAnswers := 0
	totalAnswers := len(answers)

	for _, a := range answers {
		var isValid bool
		if ua, exists := userAnswerMap[a.ID]; exists {
			isValid = a.isCorrect == ua.GetChecked()
			if isValid {
				correctAnswers++
			}
			results = append(results, &pb.AnswerResult{
				AnswerId: a.ID,
				IsValid:  isValid,
			})
		} else {
			// Handle the case where a user answer ID doesn't match any answer
			results = append(results, &pb.AnswerResult{
				AnswerId: a.ID,
				IsValid:  false,
			})
		}
	}

	var score float32 = 0
	if totalAnswers > 0 {
		score = float32(correctAnswers) / float32(totalAnswers) * 100
	}

	return results, score, nil
}

// ReorderAnswers NOTE: DEPRECATED
func ReorderAnswers(answers []Answer, payload ReorderPayload, targetID string) ([]Answer, error) {
	if len(answers) == 0 {
		return nil, errors.BadRequest("Invalid question", "question has no answers")
	}
	if len(answers) == 1 {
		return answers, nil
	}
	var newAnswers []Answer
	var target *Answer
	for _, a := range answers {
		if a.ID == targetID {
			target = &a
		}
	}
	if target == nil {
		return nil, errors.BadRequest("Invalid target ID", "target ID is empty")
	}
	if payload.MakeFirst {
		newAnswers = append(newAnswers, *target)
		for _, a := range answers {
			if a.ID != targetID {
				newAnswers = append(newAnswers, a)
			}
		}
	} else if payload.MakeLast {
		for _, a := range answers {
			if a.ID != targetID {
				newAnswers = append(newAnswers, a)
			}
		}
		newAnswers = append(newAnswers, *target)
	} else if payload.AboveID != "" {
		for i, a := range answers {
			if a.ID == payload.AboveID {
				newAnswers = append(newAnswers, answers[i-1])
				newAnswers = append(newAnswers, *target)
				newAnswers = append(newAnswers, answers[i+1])
			}
		}
	} else if payload.BelowID != "" {
		for i, a := range answers {
			if a.ID == payload.BelowID {
				newAnswers = append(newAnswers, answers[i-1])
				newAnswers = append(newAnswers, *target)
				newAnswers = append(newAnswers, answers[i+1])
			}
		}
	} else {
		return nil, errors.BadRequest("Invalid payload", "payload is invalid")
	}

	return answers, nil
}
