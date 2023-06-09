package repository

import (
	"context"
	"github.com/RafalSalwa/interview-app-srv/pkg/logger"
	"github.com/RafalSalwa/interview-app-srv/pkg/models"
	apiMongo "github.com/RafalSalwa/interview-app-srv/pkg/mongo"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

type Mongo struct {
	log *logger.Logger
	cfg apiMongo.Config
	db  *mongo.Client
}

func NewMongoRepository(db *mongo.Client, cfg apiMongo.Config, log *logger.Logger) *Mongo {
	return &Mongo{log: log, cfg: cfg, db: db}
}

func (r *Mongo) CreateUser(ctx context.Context, user *models.UserDBModel) error {
	ctx, span := otel.GetTracerProvider().Tracer("auth_service-mongo").Start(ctx, "MongoDB CreateUser")
	defer span.End()

	mongoUser := models.UserMongoModel{}
	err := mongoUser.FromDBModel(user)
	if err != nil {
		return err
	}
	collection := r.db.Database(r.cfg.Database).Collection("users")

	_, err = collection.InsertOne(ctx, mongoUser, &options.InsertOneOptions{})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return errors.Wrap(err, "InsertOne")
	}

	return nil
}

func (r *Mongo) UpdateUser(ctx context.Context, user *models.UserMongoModel) error {
	ctx, span := otel.GetTracerProvider().Tracer("auth_service-rpc").Start(ctx, "GRPC SignUpUser")
	defer span.End()

	collection := r.db.Database(r.cfg.Database).Collection("users")

	ops := options.FindOneAndUpdate()
	ops.SetReturnDocument(options.After)
	ops.SetUpsert(true)

	var updated models.UserMongoModel
	if err := collection.FindOneAndUpdate(ctx, bson.M{"_id": user.Id}, bson.M{"$set": user}, ops).Decode(&updated); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return errors.Wrap(err, "Decode")
	}

	return nil
}

func (r *Mongo) GetUser(ctx context.Context, id string) (*models.UserMongoModel, error) {
	ctx, span := otel.GetTracerProvider().Tracer("auth_service-rpc").Start(ctx, "GRPC SignUpUser")
	defer span.End()

	collection := r.db.Database(r.cfg.Database).Collection("users")

	var user models.UserMongoModel
	if err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, errors.Wrap(err, "Decode")
	}

	return &user, nil
}

func (r *Mongo) DeleteUser(ctx context.Context, id string) error {
	ctx, span := otel.GetTracerProvider().Tracer("auth_service-rpc").Start(ctx, "GRPC SignUpUser")
	defer span.End()

	collection := r.db.Database(r.cfg.Database).Collection("users")

	return collection.FindOneAndDelete(ctx, bson.M{"_id": id}).Err()
}
