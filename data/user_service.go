package data

import (
	"context"
	"errors"
	"task_manager/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection
var userCtx = context.TODO()

func init() {
	clientOptions := options.Client().ApplyURI(getMongoURI())
	client, err := mongo.Connect(userCtx, clientOptions)
	if err != nil {
		panic(err)
	}
	if err = client.Ping(userCtx, nil); err != nil {
		panic(err)
	}
	userCollection = client.Database("task_manager").Collection("users")
}

// getMongoURI returns the MongoDB connection string.
// For simplicity this returns the default local URI. Set this value
// back to environment lookup if you need environment configurability.
func getMongoURI() string {
	return "mongodb://localhost:27017"
}

// CreateUser creates a new user. If this is the first user, assign admin role.
func CreateUser(u models.User) (models.User, error) {
	// ensure username unique
	filter := bson.M{"username": u.Username}
	count, err := userCollection.CountDocuments(userCtx, filter)
	if err != nil {
		return models.User{}, err
	}
	if count > 0 {
		return models.User{}, errors.New("username already exists")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, err
	}
	u.Password = string(hashed)
	u.CreatedAt = time.Now()

	total, err := userCollection.CountDocuments(userCtx, bson.M{})
	if err == nil && total == 0 {
		u.Role = "admin"
	} else if u.Role == "" {
		u.Role = "user"
	}

	res, err := userCollection.InsertOne(userCtx, u)
	if err != nil {
		return models.User{}, err
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		u.ID = oid
	}
	u.Password = ""
	return u, nil
}

// AuthenticateUser verifies credentials and returns user without password.
func AuthenticateUser(username, password string) (models.User, error) {
	var u models.User
	err := userCollection.FindOne(userCtx, bson.M{"username": username}).Decode(&u)
	if err != nil {
		return models.User{}, errors.New("invalid credentials")
	}
	if err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return models.User{}, errors.New("invalid credentials")
	}
	u.Password = ""
	return u, nil
}

// PromoteUser sets role to admin. Returns updated user.
func PromoteUser(username string) (models.User, error) {
	filter := bson.M{"username": username}
	update := bson.M{"$set": bson.M{"role": "admin"}}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	res := userCollection.FindOneAndUpdate(userCtx, filter, update, opts)
	var u models.User
	if err := res.Decode(&u); err != nil {
		return models.User{}, errors.New("user not found")
	}
	u.Password = ""
	return u, nil
}

// GetUserByUsername returns a user struct (including password hash) — internal use.
func GetUserByUsername(username string) (models.User, error) {
	var u models.User
	err := userCollection.FindOne(userCtx, bson.M{"username": username}).Decode(&u)
	return u, err
}

// GetUserByID returns user by ObjectID.
func GetUserByID(id string) (models.User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.User{}, err
	}
	var u models.User
	err = userCollection.FindOne(userCtx, bson.M{"_id": oid}).Decode(&u)
	return u, err
}
