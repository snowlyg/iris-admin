package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Client struct {
	mc *mongo.Client
}

// GetClient
func GetClient() (*Client, error) {
	mc, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(CONFIG.GetApplyURI()))
	if err != nil {
		return nil, err
	}
	client := &Client{mc: mc}
	return client, nil
}

// Ping
func (c *Client) Ping() error {
	return c.mc.Ping(context.TODO(), readpref.Primary())
}

// getCollection
func (c *Client) getCollection(name string) *mongo.Collection {
	return c.mc.Database(CONFIG.DB).Collection(name)
}

// Aggregate
func (c *Client) Aggregate(name string, groupStage mongo.Pipeline) ([]bson.M, error) {
	// pass the stage into a pipeline
	// pass the pipeline as the second paramter in the Aggregate() method
	cursor, err := c.getCollection(name).Aggregate(context.TODO(), groupStage)
	if err != nil {
		return nil, err
	}
	// display the results
	results := []bson.M{}
	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())
	return results, nil
}

// Find
func (c *Client) Find(name string, filters ...interface{}) ([]bson.M, error) {
	var cursor *mongo.Cursor
	var err error
	if len(filters) == 0 {
		cursor, err = c.getCollection(name).Find(context.TODO(), bson.D{})
		if err != nil {
			return nil, err
		}
	} else {
		cursor, err = c.getCollection(name).Find(context.TODO(), filters[0])
		if err != nil {
			return nil, err
		}
	}
	var results []bson.M
	for cursor.Next(context.TODO()) {
		b := bson.M{}
		err := cursor.Decode(b)
		if err != nil {
			return nil, err
		}
		results = append(results, b)
	}
	return results, nil
}

// FindOne
func (c *Client) FindOne(name string, filter interface{}) *mongo.SingleResult {
	return c.getCollection(name).FindOne(context.TODO(), filter)
}

// Disconnect
func (c *Client) Disconnect() error {
	return c.mc.Disconnect(context.TODO())
}

// InsertOne
func (c *Client) InsertOne(name string, filter interface{}) (interface{}, error) {
	cur, err := c.getCollection(name).InsertOne(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	return cur.InsertedID, nil
}

// DeleteOne
func (c *Client) DeleteOne(name string, filter interface{}) error {
	_, err := c.getCollection(name).DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	return nil
}

// UpdateByID
func (c *Client) UpdateOne(name string, filter, update interface{}) (*mongo.UpdateResult, error) {
	return c.getCollection(name).UpdateOne(context.TODO(), filter, update)
}
