package db

import (
	"context"
	"math"
	"time"

	"github.com/abdollahpour/almaniha-draft/internal/config"
	"github.com/abdollahpour/almaniha-draft/internal/model"
	"github.com/abdollahpour/almaniha-draft/internal/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDb struct {
	client     mongo.Client
	collection mongo.Collection
}

func NewMongoDb(conf config.Configuration, drop bool) *MongoDb {
	db := util.ConnectMongoDB(conf)

	collection := db.Collection("drafts")

	if drop {
		collection.Drop(context.TODO())
	}

	indexes := collection.Indexes()

	// Full text search
	typeIndex := mongo.IndexModel{
		Keys: bson.D{
			{"type", 1},
		},
	}
	_, err := indexes.CreateOne(context.TODO(), typeIndex)
	if err != nil {
		panic(err)
	}

	return &MongoDb{
		collection: *collection,
	}
}

func (mp *MongoDb) Read(id string) (*model.Draft, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var draft model.Draft

	err = mp.collection.FindOne(context.TODO(), bson.M{"_id": _id}).Decode(&draft)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &draft, nil
}

func (mp *MongoDb) CreateOrUpdate(draft model.Draft) (*model.Draft, error) {
	_id := draft.ID
	// Create a new record
	if _id.IsZero() {
		_id = primitive.NewObjectID()
	}

	now := time.Now()
	update := bson.M{
		"$set":         bson.M{"data": draft.Data, "upatedAt": now},
		"$setOnInsert": bson.M{"createdAt": time.Now(), "type": draft.Type, "createdBy": draft.CreatedBy},
	}

	upsert := true
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}

	var updated model.Draft

	if err := mp.collection.FindOneAndUpdate(context.TODO(), bson.M{"_id": _id}, update, &opt).Decode(&updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

func (mp *MongoDb) Delete(id string) (*model.Draft, error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var draft model.Draft

	if err = mp.collection.FindOneAndDelete(context.TODO(), bson.M{"_id": _id}).Decode(&draft); err != nil {
		return nil, err
	}
	return &draft, nil
}

func (mp *MongoDb) Paginated(page int64, size int64, filter model.DraftFilter, sort model.DraftSort) (*model.PaginatedDraft, error) {
	findOptions := options.Find()
	findOptions.SetSkip((page - 1) * size)
	findOptions.SetLimit(size)

	findFilter := bson.M{}
	if len(filter.Type) > 0 {
		findFilter["type"] = filter.Type
	}

	// var sortOrder int
	// if sort.Order == model.DESC {
	// 	sortOrder = -1
	// } else {
	// 	sortOrder = 1
	// }
	// findOptions.SetSort(bson.E{Key: string(sort.Field), Value: sortOrder})
	var sortOrder int
	if sort.Order == model.DESC {
		sortOrder = -1
	} else {
		sortOrder = 1
	}
	if len(string(sort.Field)) > 0 {
		findOptions.SetSort(bson.D{{string(sort.Field), sortOrder}})
	}

	cursor, err := mp.collection.Find(context.TODO(), findFilter, findOptions)
	if err != nil {
		panic(err)
	}
	var drafts []model.Draft
	err = cursor.All(context.TODO(), &drafts)
	if err != nil {
		return nil, err
	}

	total, err := mp.collection.CountDocuments(context.TODO(), bson.D{{}})
	if err != nil {
		panic(err)
	}
	totalPages := int64(math.Ceil(float64(total) / float64(size)))

	return &model.PaginatedDraft{
		Paginated: model.Paginated{
			Page:       page,
			Size:       size,
			Total:      total,
			TotalPages: totalPages,
			HasNext:    page < totalPages,
			HasPrev:    page > 1,
		},
		Data:   drafts,
		Filter: filter,
		Sort:   sort,
	}, nil
}
