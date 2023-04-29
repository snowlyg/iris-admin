package mongodb

import (
	"testing"

	"github.com/snowlyg/iris-admin/g"
	"github.com/snowlyg/iris-admin/server/cache"
	_ "github.com/snowlyg/iris-admin/server/cache"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func TestGetClient(t *testing.T) {
	CONFIG.Addr = g.TestMongoAddr
	defer Remove()
	defer cache.Remove()
	t.Run("test mongodb getClient", func(t *testing.T) {
		client, err := GetClient()
		if err != nil {
			t.Error(err.Error())
			return
		}
		if client == nil {
			t.Error("mongodb clinet is nil")
		}
		defer func() {
			if err = client.Disconnect(); err != nil {
				t.Error(err)
			}
		}()
	})

}
func TestPing(t *testing.T) {
	CONFIG.Addr = g.TestMongoAddr
	defer Remove()
	defer cache.Remove()
	t.Run("test mongodb ping", func(t *testing.T) {
		client, err := GetClient()
		if err != nil {
			t.Error(err.Error())
			return
		}
		if client == nil {
			t.Error("mongodb clinet is nil")
		}
		defer func() {
			if err = client.Disconnect(); err != nil {
				t.Error(err)
			}
		}()
		err = client.Ping()
		if err != nil {
			t.Error(err.Error())
			return
		}
	})
}
func TestInsertOne(t *testing.T) {
	CONFIG.Addr = g.TestMongoAddr
	defer Remove()
	defer cache.Remove()
	t.Run("test mongodb InsertOne", func(t *testing.T) {
		client, err := GetClient()
		if err != nil {
			t.Error(err.Error())
			return
		}
		if client == nil {
			t.Error("mongodb clinet is nil")
		}
		defer func() {
			if err = client.Disconnect(); err != nil {
				t.Error(err)
			}
		}()
		res, err := client.InsertOne("testing", bson.D{
			{Key: "name", Value: "pi"}, {Key: "value", Value: 3.14159},
		})
		if err != nil {
			t.Error(err.Error())
			return
		}
		if res == "" {
			t.Error("inserted id is empty")
		}
	})
}
func TestGetCollection(t *testing.T) {
	CONFIG.Addr = g.TestMongoAddr
	defer Remove()
	defer cache.Remove()
	t.Run("test mongodb GetCollection", func(t *testing.T) {
		client, err := GetClient()
		if err != nil {
			t.Error(err.Error())
			return
		}
		if client == nil {
			t.Error("mongodb clinet is nil")
		}
		defer func() {
			if err = client.Disconnect(); err != nil {
				t.Error(err)
			}
		}()
		res := client.getCollection("testing")
		if res == nil {
			t.Error("Collection return empty")
		}
	})
}

func TestGetAggregate(t *testing.T) {
	CONFIG.Addr = g.TestMongoAddr
	defer Remove()
	defer cache.Remove()
	t.Run("test mongodb Aggregate", func(t *testing.T) {
		client, err := GetClient()
		if err != nil {
			t.Error(err.Error())
			return
		}
		if client == nil {
			t.Error("mongodb clinet is nil")
			return
		}
		defer func() {
			if err = client.Disconnect(); err != nil {
				t.Error(err)
			}
		}()
		pipeline := mongo.Pipeline{
			{
				{"$match", bson.D{
					{"items.fruit", "banana"},
				}},
			},
			{
				{"$sort", bson.D{
					{"date", 1},
				}},
			},
		}
		res, err := client.Aggregate("testing", pipeline)
		if err != nil {
			t.Error(err.Error())
			return
		}
		if res == nil {
			t.Error("Collection return empty")
		}
	})
}

func TestFind(t *testing.T) {
	CONFIG.Addr = g.TestMongoAddr
	defer Remove()
	defer cache.Remove()
	t.Run("test mongodb Find", func(t *testing.T) {
		client, err := GetClient()
		if err != nil {
			t.Error(err.Error())
			return
		}
		if client == nil {
			t.Error("mongodb clinet is nil")
			return
		}
		defer func() {
			if err = client.Disconnect(); err != nil {
				t.Error(err)
			}
		}()
		res, err := client.Find("testing", bson.D{{"end", nil}})
		if err != nil {
			t.Error(err.Error())
			return
		}
		if res == nil {
			t.Error("Collection return empty")
		}
	})
}

func TestFindOne(t *testing.T) {
	CONFIG.Addr = g.TestMongoAddr
	defer Remove()
	defer cache.Remove()
	t.Run("test mongodb FindOne", func(t *testing.T) {
		client, err := GetClient()
		if err != nil {
			t.Error(err.Error())
			return
		}
		if client == nil {
			t.Error("mongodb clinet is nil")
			return
		}
		defer func() {
			if err = client.Disconnect(); err != nil {
				t.Error(err)
			}
		}()
		res := client.FindOne("testing", bson.D{{"end", nil}})
		if res == nil {
			t.Error("Collection return empty")
		}
	})
}

func TestDeleteOne(t *testing.T) {
	CONFIG.Addr = g.TestMongoAddr
	defer Remove()
	defer cache.Remove()
	t.Run("test mongodb DeleteOne", func(t *testing.T) {
		client, err := GetClient()
		if err != nil {
			t.Error(err.Error())
			return
		}
		if client == nil {
			t.Error("mongodb clinet is nil")
			return
		}
		defer func() {
			if err = client.Disconnect(); err != nil {
				t.Error(err)
			}
		}()
		err = client.DeleteOne("testing", bson.D{{"end", nil}})
		if err != nil {
			t.Error(err.Error())
			return
		}
	})
}

func TestUpdateOne(t *testing.T) {
	CONFIG.Addr = g.TestMongoAddr
	defer Remove()
	defer cache.Remove()
	t.Run("test mongodb UpdateOne", func(t *testing.T) {
		client, err := GetClient()
		if err != nil {
			t.Error(err.Error())
			return
		}
		if client == nil {
			t.Error("mongodb clinet is nil")
			return
		}
		defer func() {
			if err = client.Disconnect(); err != nil {
				t.Error(err)
			}
		}()
		id, err := client.InsertOne("testing", bson.D{
			{Key: "name", Value: "pi"}, {Key: "value", Value: 3.14159},
		})
		if err != nil {
			t.Error(err.Error())
			return
		}
		b := bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "name", Value: "pi"},
				{Key: "value", Value: 3.1415926},
			}},
		}
		res, err := client.UpdateOne("testing", bson.D{{"_id", id}}, b)
		if err != nil {
			t.Error(err.Error())
			return
		}
		if res == nil {
			t.Error("Collection return empty")
		}
	})
}
