package repositories

import (
	"context"
	domain "task_manager/Domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	GetByID(id string) (domain.User, error)
	GetByUsername(username string) (domain.User, error)
	Create(user domain.User) (domain.User, error)
	Update(id string, user domain.User) (domain.User, error)
	Delete(id string) error
	PromoteToAdmin(username string) (domain.User, error)
}

// MongoUserRepository implements UserRepository using MongoDB
type MongoUserRepository struct {
	collection *mongo.Collection
}

func NewMongoUserRepository(client *mongo.Client, dbName, collName string) *MongoUserRepository {
	coll := client.Database(dbName).Collection(collName)
	return &MongoUserRepository{collection: coll}
}

func (r *MongoUserRepository) GetByID(id string) (domain.User, error) {
	ctx := context.Background()
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.User{}, err
	}
	var user domain.User
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&user)
	return user, err
}

func (r *MongoUserRepository) GetByUsername(username string) (domain.User, error) {
	ctx := context.Background()
	var user domain.User
	err := r.collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	return user, err
}

func (r *MongoUserRepository) Create(user domain.User) (domain.User, error) {
	ctx := context.Background()
	res, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return domain.User{}, err
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		user.ID = oid.Hex()
	}
	return user, nil
}

func (r *MongoUserRepository) Update(id string, user domain.User) (domain.User, error) {
	ctx := context.Background()
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.User{}, err
	}
	filter := bson.M{"_id": objID}
	update := bson.M{"$set": user}
	_, err = r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return domain.User{}, err
	}
	user.ID = objID.Hex()
	return user, nil
}

func (r *MongoUserRepository) Delete(id string) error {
	ctx := context.Background()
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}

func (r *MongoUserRepository) PromoteToAdmin(username string) (domain.User, error) {
	ctx := context.Background()
	filter := bson.M{"username": username}
	update := bson.M{"$set": bson.M{"role": "admin"}}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return domain.User{}, err
	}
	return r.GetByUsername(username)
}
