package mongodb

import (
	"context"
	"os"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func TestGetClient(t *testing.T) {
	CONFIG.Addr = os.Getenv("mongoAddr")
	InitMongoDBConfig()
	defer Remove()
	ctx, cancel := context.WithTimeout(context.Background(), CONFIG.Timeout*time.Second)
	defer cancel()
	t.Run("test mongodb getClient", func(t *testing.T) {
		client, err := GetClient(ctx)
		if err != nil {
			t.Error(err.Error())
			return
		}
		if client == nil {
			t.Error("mongodb clinet is nil")
		}
		defer func() {
			if err = client.Disconnect(ctx); err != nil {
				t.Error(err)
			}
		}()
	})

}
func TestPing(t *testing.T) {
	CONFIG.Addr = os.Getenv("mongoAddr")
	InitMongoDBConfig()
	defer Remove()
	ctx, cancel := context.WithTimeout(context.Background(), CONFIG.Timeout*time.Second)
	defer cancel()
	t.Run("test mongodb ping", func(t *testing.T) {
		client, err := GetClient(ctx)
		if err != nil {
			t.Error(err.Error())
			return
		}
		if client == nil {
			t.Error("mongodb clinet is nil")
		}
		defer func() {
			if err = client.Disconnect(ctx); err != nil {
				t.Error(err)
			}
		}()
		err = client.Ping(ctx)
		if err != nil {
			t.Error(err.Error())
			return
		}
	})
}
func TestInsertOne(t *testing.T) {
	CONFIG.Addr = os.Getenv("mongoAddr")
	InitMongoDBConfig()
	defer Remove()
	ctx, cancel := context.WithTimeout(context.Background(), CONFIG.Timeout*time.Second)
	defer cancel()
	t.Run("test mongodb InsertOne", func(t *testing.T) {
		client, err := GetClient(ctx)
		if err != nil {
			t.Error(err.Error())
			return
		}
		if client == nil {
			t.Error("mongodb clinet is nil")
		}
		defer func() {
			if err = client.Disconnect(ctx); err != nil {
				t.Error(err)
			}
		}()
		res, err := client.InsertOne(ctx, "testing", bson.D{
			{"name", "pi"}, {"value", 3.14159},
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
	CONFIG.Addr = os.Getenv("mongoAddr")
	InitMongoDBConfig()
	defer Remove()
	ctx, cancel := context.WithTimeout(context.Background(), CONFIG.Timeout*time.Second)
	defer cancel()
	t.Run("test mongodb GetCollection", func(t *testing.T) {
		client, err := GetClient(ctx)
		if err != nil {
			t.Error(err.Error())
			return
		}
		if client == nil {
			t.Error("mongodb clinet is nil")
		}
		defer func() {
			if err = client.Disconnect(ctx); err != nil {
				t.Error(err)
			}
		}()
		res := client.getCollection("testing")
		if res == nil {
			t.Error("Collection return empty")
		}
	})
}
