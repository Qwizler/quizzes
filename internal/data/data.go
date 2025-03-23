package data

import (
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/surrealdb/surrealdb.go"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"

	"quiz/internal/conf"
	"quiz/internal/dep"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var DataProviderSet = wire.NewSet(NewData, NewQuizRepo, NewQuestionsRepo)

// Data .
type Data struct {
	gorm    *gorm.DB
	mongo   *mongo.Database
	surreal *surrealdb.DB
	logger  log.Logger
}

// NewData .
func NewData(c *conf.Data, logger log.Logger, tp trace.TracerProvider) (*Data, func(), error) {
	lg := log.NewHelper(logger)
	var g *dep.Gorm
	var m *dep.Mongo
	var s *dep.Surreal
	var mongoClean func()
	var surrealClean func()
	var err error
	noDB := true

	if c.GetDatabase().GetSource() != "" && c.GetDatabase().GetDriver() != "" {
		g, _, err = dep.NewGorm(c, logger, tp)
		if err != nil {
			lg.Warn("failed to connect to PostgreSQL", err)
			return nil, nil, err
		}
		noDB = false
	}

	if c.GetMongo().GetUri() != "" && c.Mongo.Database != "" && c.Mongo.Username != "" && c.Mongo.Password != "" {
		m, mongoClean, err = dep.NewMongo(c, logger)
		if err != nil {
			lg.Warn("failed to connect to MongoDB", err)
			return nil, nil, err
		}
		noDB = false
	}

	if c.GetSurreal().GetAddr() != "" && c.Surreal.Database != "" && c.Surreal.Namespace != "" && c.Surreal.Username != "" && c.Surreal.Password != "" {
		s, surrealClean, err = dep.NewSurreal(c, logger)
		if err != nil {
			lg.Warn("failed to connect to SurrealDB", err)
			return nil, nil, err
		}
		noDB = false
	}

	cleanup := func() {
		if mongoClean != nil {
			mongoClean()
		}
		if surrealClean != nil {
			surrealClean()
		}
		log.NewHelper(logger).Info("closing the data resources")
	}

	if noDB {
		return nil, nil, errors.InternalServer("no database configured", "no database configured")
	}

	data := &Data{logger: logger}
	if g != nil {
		lg.Debug("Attaching PostgreSQL")
		data.gorm = g.DB
	}
	if m != nil {
		lg.Debug("Attaching MongoDB")
		data.mongo = m.DB
	}
	if s != nil {
		lg.Debug("Attaching SurrealDB")
		data.surreal = s.DB
	}

	return data, cleanup, nil
}
