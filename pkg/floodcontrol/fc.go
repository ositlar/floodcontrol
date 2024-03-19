package floodcontrol

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FloodController struct {
	client     *mongo.Client
	collection *mongo.Collection
	l          int
	n          time.Duration
}

func NewFloodController(collectionName, dbName, conn string, l int, n time.Duration) (*FloodController, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(conn))
	if err != nil {
		return nil, err
	}
	collection := client.Database(dbName).Collection(collectionName)
	return &FloodController{
		client:     client,
		collection: collection,
		l:          l,
		n:          n,
	}, nil
}

func (fc *FloodController) Check(ctx context.Context, userID int64) (bool, error) {
	filter := make(map[string]interface{})
	filter["userID"] = userID
	var result struct {
		LastCallTime time.Time `json:"lastCallTime"`
	}
	err := fc.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return false, err
	}
	update := map[string]interface{}{
		"$set": map[string]interface{}{
			"userID":       userID,
			"lastCallTime": time.Now(),
		},
	}
	opts := options.Update().SetUpsert(true)
	_, err = fc.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return false, err
	}
	if time.Since(result.LastCallTime) >= fc.n {
		return true, nil
	}
	return false, nil
}

func (fc *FloodController) Disconnect(ctx context.Context) error {
	return fc.client.Disconnect(ctx)
}
