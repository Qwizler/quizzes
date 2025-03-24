package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/jinzhu/copier"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"quiz/internal/biz"
	"time"
)

type Answer struct {
	ID          string `bson:"id,omitempty"`
	Text        string `bson:"answer"`
	IsCorrect   bool   `bson:"is_correct"`
	Explanation string `bson:"explanation"`
}

type Question struct {
	ID         bson.ObjectID `bson:"_id,omitempty"`
	QuizID     string        `bson:"quiz_id"`
	Question   string        `bson:"question"`
	Difficulty *uint64       `bson:"difficulty"`
	Order      float64       `bson:"order"`
	Answers    []Answer      `bson:"answers"`
	Hint       string        `bson:"hint"`
	CreatedBy  string        `bson:"created_by"`
	UpdatedBy  string        `bson:"updated_by"`
	DeletedBy  string        `bson:"deleted_by"`
	CreatedAt  string        `bson:"created_at"`
	UpdatedAt  string        `bson:"updated_at"`
	DeletedAt  string        `bson:"deleted_at"`
}

type QuestionsRepo struct {
	coll   *mongo.Collection
	log    *log.Helper
	tracer trace.Tracer
	table  string
}

func NewQuestionsRepo(data *Data, logger log.Logger, tracer trace.Tracer) biz.QuestionsRepo {
	return &QuestionsRepo{
		coll:   data.mongo.Collection("questions"),
		log:    log.NewHelper(logger),
		tracer: tracer,
		table:  "questions",
	}
}

func (r *QuestionsRepo) Save(ctx context.Context, q *biz.Question) (*biz.Question, error) {
	ctx, span := r.tracer.Start(ctx, "data.QuestionsRepo.Save")
	defer span.End()
	createdAt := time.Now().String()

	var question Question
	err := copier.Copy(&question, q)
	if err != nil {
		r.log.Warn(err)
		return nil, err
	}
	question.CreatedAt = createdAt
	question.UpdatedAt = createdAt

	res, err := r.coll.InsertOne(ctx, question)
	if err != nil {
		r.log.Warn(err)
		return nil, err
	}

	if oid, ok := res.InsertedID.(bson.ObjectID); ok {
		hex := oid.Hex()
		r.log.Debugf("inserted id is object id: %s", hex)
		resQuestion := question.Biz()
		resQuestion.ID = hex
		r.log.Debugf("inserted question: %+v", resQuestion)
		return resQuestion, nil
	}
	return nil, errors.InternalServer("inserted id is not object id", "inserted id is not object id")
}

func (r QuestionsRepo) GetByID(ctx context.Context, id string) (*biz.Question, error) {
	ctx, span := r.tracer.Start(ctx, "data.QuestionsRepo.GetByID")
	defer span.End()
	span.SetAttributes(attribute.KeyValue{
		Key:   "id",
		Value: attribute.StringValue(id),
	})
	idObj, err := bson.ObjectIDFromHex(id)
	if err != nil {
		r.log.Warn(err)
		return nil, err
	}
	res := r.coll.FindOne(ctx, bson.M{"_id": idObj})
	if res.Err() != nil {
		r.log.Warn(res.Err())
		return nil, res.Err()
	}
	if res.Err() != nil {
		r.log.Warn(res.Err())
		return nil, res.Err()
	}
	var q Question
	err = res.Decode(&q)
	if err != nil {
		r.log.Warn(err)
		return nil, err
	}
	return q.Biz(), nil
}

func (r QuestionsRepo) List(ctx context.Context, quizID string, pagination *biz.Pagination) ([]*biz.Question, error) {
	ctx, span := r.tracer.Start(ctx, "data.QuestionsRepo.List")
	defer span.End()
	// error the provided hex string is not a valid ObjectID"

	oid, err := bson.ObjectIDFromHex(quizID)
	if err != nil {
		r.log.Warn(err)
		return nil, err
	}
	opts := options.Find().SetSkip(int64(pagination.Page * pagination.Size)).SetLimit(int64(pagination.Size))
	filter := bson.M{"quiz_id": oid}
	cur, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		r.log.Warn(err)
		return nil, err
	}
	var res []*biz.Question
	for cur.Next(ctx) {
		var q Question
		if err := cur.Decode(&q); err != nil {
			r.log.Warn(err)
			return nil, err
		}
		res = append(res, q.Biz())
	}
	return res, nil
}

func (r QuestionsRepo) Update(ctx context.Context, q *biz.Question) (*biz.Question, error) {
	ctx, span := r.tracer.Start(ctx, "data.QuestionsRepo.Update")
	defer span.End()

	question := &Question{
		Question:   q.Question,
		Difficulty: &q.Difficulty,
		Order:      q.Order,
	}
	idObj, err := bson.ObjectIDFromHex(q.ID)
	if err != nil {
		r.log.Warn(err)
		return nil, err
	}
	res, err := r.coll.UpdateOne(ctx, bson.M{"_id": idObj}, bson.M{"$set": question})
	if err != nil {
		r.log.Warn(err)
		return nil, err
	}
	if res.Acknowledged == false {
		return nil, errors.InternalServer("failed to update quiz", "failed to update quiz")
	}
	if res.ModifiedCount == 0 {
		return nil, errors.NotFound("quiz not found", "quiz not found")
	}
	return question.Biz(), nil
}

func (r QuestionsRepo) Delete(ctx context.Context, id string) (*biz.Question, error) {
	ctx, span := r.tracer.Start(ctx, "data.QuestionsRepo.Delete")
	defer span.End()

	idObj, err := bson.ObjectIDFromHex(id)
	if err != nil {
		r.log.Warn(err)
		return nil, err
	}
	res, err := r.coll.DeleteOne(ctx, bson.M{"_id": idObj})
	if err != nil {
		r.log.Warn(err)
		return nil, err
	}
	if res.DeletedCount == 0 {
		return nil, errors.NotFound("quiz not found", "quiz not found")
	}
	return &biz.Question{
		ID: id,
	}, nil
}
