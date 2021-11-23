package db

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/abdollahpour/almaniha-draft/internal/config"
	"github.com/abdollahpour/almaniha-draft/internal/model"
	"github.com/abdollahpour/almaniha-draft/internal/util"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func teardown() {
	util.DisconnectMongoDB()
}

func newLocalClient() *mongo.Client {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	return client
}

func TestCreate(t *testing.T) {
	client := newLocalClient()
	defer client.Disconnect(context.Background())

	business := model.Draft{}

	mongoDb := NewMongoDb(config.NewEnvConfiguration(), true)

	added, err := mongoDb.CreateOrUpdate(business)
	assert.Nil(t, err)
	assert.NotNil(t, added.ID)
}

func TestUpdate(t *testing.T) {
	client := newLocalClient()
	defer client.Disconnect(context.Background())

	_id := primitive.NewObjectID()
	data := map[string]interface{}{"key": "value"}
	business := model.Draft{
		ID:   _id,
		Data: data,
	}

	mongoDb := NewMongoDb(config.NewEnvConfiguration(), true)

	updated, err := mongoDb.CreateOrUpdate(business)
	assert.Nil(t, err)
	assert.Equal(t, data, updated.Data)
	assert.Equal(t, _id, updated.ID)
}

func TestRead(t *testing.T) {
	draft := model.Draft{}

	mongoDb := NewMongoDb(config.NewEnvConfiguration(), true)

	add, err := mongoDb.CreateOrUpdate(draft)
	assert.Nil(t, err)
	assert.NotNil(t, add.ID)

	read, err := mongoDb.Read(add.ID.Hex())
	assert.Nil(t, err)
	assert.Equal(t, add, read)
}

func TestDelete(t *testing.T) {
	draft := model.Draft{}

	mongoDb := NewMongoDb(config.NewEnvConfiguration(), true)

	added, err := mongoDb.CreateOrUpdate(draft)
	assert.Nil(t, err)
	assert.NotNil(t, added.ID)

	delete, err := mongoDb.Delete(added.ID.Hex())
	assert.Nil(t, err)
	assert.Equal(t, added.ID, delete.ID)

	readed, err := mongoDb.Read(added.ID.Hex())
	assert.Nil(t, readed)
	assert.Nil(t, err)
}

func TestPaginatedSimple(t *testing.T) {
	draft := model.Draft{
		Data: map[string]interface{}{"key": "value"},
	}

	mongoDb := NewMongoDb(config.NewEnvConfiguration(), true)

	mongoDb.CreateOrUpdate(draft)
	mongoDb.CreateOrUpdate(draft)
	mongoDb.CreateOrUpdate(draft)
	mongoDb.CreateOrUpdate(draft)

	page, err := mongoDb.Paginated(1, 10, model.DraftFilter{}, model.DraftSort{})
	assert.Nil(t, err)
	assert.NotNil(t, page)
	assert.Equal(t, 4, len(page.Data))
}

func TestPaginatedValidKeyword(t *testing.T) {
	mongoDb := NewMongoDb(config.NewEnvConfiguration(), true)

	mongoDb.CreateOrUpdate(model.Draft{
		Type: "Business",
		Data: map[string]interface{}{"key1": "value1"},
	})
	mongoDb.CreateOrUpdate(model.Draft{
		Type: "Business",
		Data: map[string]interface{}{"key2": "value2"},
	})
	mongoDb.CreateOrUpdate(model.Draft{
		Type: "Business",
		Data: map[string]interface{}{"key3": "value3"},
	})
	mongoDb.CreateOrUpdate(model.Draft{
		Type: "Business",
		Data: map[string]interface{}{"key4": "value4"},
	})

	page, err := mongoDb.Paginated(1, 10, model.DraftFilter{
		Type: "Business",
	}, model.DraftSort{})
	assert.Nil(t, err)
	assert.NotNil(t, page)
	assert.Equal(t, 4, len(page.Data))
}
