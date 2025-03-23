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

type Quiz struct {
	ID          bson.ObjectID     `bson:"_id,omitempty"`
	UserID      string            `bson:"user_id"`
	Title       string            `bson:"title"`
	Description string            `bson:"description"`
	Duration    *uint64           `bson:"duration"`
	Thumbnail   *string           `bson:"thumbnail"`
	Cover       *string           `bson:"cover"`
	Category    *string           `bson:"category"`
	Tags        []string          `bson:"tags"`
	Metadata    map[string]string `bson:"metadata"`
	CreatedBy   string            `bson:"created_by"`
	UpdatedBy   string            `bson:"updated_by"`
	DeletedBy   string            `bson:"deleted_by"`
	CreatedAt   string            `bson:"created_at"`
	UpdatedAt   string            `bson:"updated_at"`
	DeletedAt   string            `bson:"deleted_at"`
}

type QuizRepo struct {
	coll   *mongo.Collection
	log    *log.Helper
	tracer trace.Tracer
	table  string
}

func NewQuizRepo(data *Data, logger log.Logger, tracer trace.Tracer) biz.QuizRepo {
	return &QuizRepo{
		coll:   data.mongo.Collection("quizzes"),
		log:    log.NewHelper(logger),
		tracer: tracer,
		table:  "quizzes",
	}
}

func (r *QuizRepo) Save(ctx context.Context, q *biz.Quiz) (*biz.Quiz, error) {
	ctx, span := r.tracer.Start(ctx, "data.QuizRepo.Save")
	defer span.End()
	createdAt := time.Now().String()

	var quiz Quiz
	err := copier.Copy(&quiz, q)
	if err != nil {
		r.log.Warn(err)
		return nil, err
	}
	quiz.CreatedBy = q.UserID
	quiz.CreatedAt = createdAt
	quiz.UpdatedBy = q.UserID
	quiz.UpdatedAt = createdAt

	res, err := r.coll.InsertOne(ctx, quiz)
	if err != nil {
		r.log.Warn(err)
		return nil, err
	}

	if oid, ok := res.InsertedID.(bson.ObjectID); ok {
		hex := oid.Hex()
		r.log.Debugf("inserted id is object id: %s", hex)
		resQuiz := quiz.QuizToBiz()
		return resQuiz, nil
	}
	return nil, errors.InternalServer("inserted id is not object id", "inserted id is not object id")
}

func (r *QuizRepo) GetByID(ctx context.Context, id string) (*biz.Quiz, error) {
	ctx, span := r.tracer.Start(ctx, "data.QuizRepo.GetByID", trace.WithAttributes(attribute.String("id", id)))
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
	var q Quiz
	err = res.Decode(&q)
	if err != nil {
		r.log.Warn(err)
		return nil, err
	}
	return q.QuizToBiz(), nil
}

func (r *QuizRepo) List(ctx context.Context, pagination *biz.Pagination) ([]*biz.Quiz, error) {
	ctx, span := r.tracer.Start(ctx, "data.QuizRepo.List")
	defer span.End()
	opts := options.Find().SetSkip(int64(pagination.Page * pagination.Size)).SetLimit(int64(pagination.Size))
	cur, err := r.coll.Find(ctx, bson.M{}, opts)
	if err != nil {
		r.log.Warn(err)
		return nil, err
	}
	var res []*biz.Quiz
	for cur.Next(ctx) {
		var q Quiz
		if err := cur.Decode(&q); err != nil {
			r.log.Warn(err)
			return nil, err
		}
		res = append(res, q.QuizToBiz())
	}
	return res, nil
}

func (r *QuizRepo) Update(ctx context.Context, q *biz.Quiz) (*biz.Quiz, error) {
	ctx, span := r.tracer.Start(ctx, "data.QuizRepo.Update")
	defer span.End()

	quiz := &Quiz{
		Title:       q.Title,
		Description: q.Description,
		Duration:    q.Duration,
		Thumbnail:   q.Thumbnail,
		Cover:       q.Cover,
		Category:    q.Category,
		Tags:        q.Tags,
		Metadata:    q.Metadata,
	}
	idObj, err := bson.ObjectIDFromHex(q.ID)
	if err != nil {
		r.log.Warn(err)
		return nil, err
	}
	quiz.UpdatedAt = time.Now().String()
	res, err := r.coll.UpdateOne(ctx, bson.M{"_id": idObj}, bson.M{"$set": quiz})
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
	return quiz.QuizToBiz(), nil
}

func (r *QuizRepo) Delete(ctx context.Context, id string) (*biz.Quiz, error) {
	ctx, span := r.tracer.Start(ctx, "data.QuizRepo.Delete")
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
	return &biz.Quiz{
		ID: id,
	}, nil
}
func (r *QuizRepo) Search(ctx context.Context, keyword string, pagination *biz.Pagination) ([]*biz.Quiz, error) {
	ctx, span := r.tracer.Start(ctx, "data.QuizRepo.Search")
	defer span.End()

	return nil, errors.InternalServer("not implemented", "not implemented")
}
