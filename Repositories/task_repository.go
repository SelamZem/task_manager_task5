package repositories

import (
	"context"
	domain "task_manager/Domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TaskRepository interface {
	GetAllTasks() ([]domain.Task, error)
	GetTaskByID(id string) (domain.Task, error)
	CreateTask(task domain.Task) (domain.Task, error)
	UpdateTask(id string, task domain.Task) (domain.Task, error)
	DeleteTask(id string) error
}

type MongoTaskRepository struct {
	collection *mongo.Collection
}

func NewMongoTaskRepository(client *mongo.Client, dbName, collName string) *MongoTaskRepository {
	coll := client.Database(dbName).Collection(collName)
	return &MongoTaskRepository{collection: coll}
}

func (r *MongoTaskRepository) GetAllTasks() ([]domain.Task, error) {
	ctx := context.Background()
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var tasks []domain.Task
	for cursor.Next(ctx) {
		var task domain.Task
		if err := cursor.Decode(&task); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *MongoTaskRepository) GetTaskByID(id string) (domain.Task, error) {
	ctx := context.Background()
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.Task{}, err
	}
	var task domain.Task
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&task)
	return task, err
}

func (r *MongoTaskRepository) CreateTask(task domain.Task) (domain.Task, error) {
	ctx := context.Background()
	res, err := r.collection.InsertOne(ctx, task)
	if err != nil {
		return domain.Task{}, err
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		task.ID = oid.Hex()
	}
	return task, nil
}

func (r *MongoTaskRepository) UpdateTask(id string, task domain.Task) (domain.Task, error) {
	ctx := context.Background()
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return domain.Task{}, err
	}
	filter := bson.M{"_id": objID}
	update := bson.M{"$set": task}
	_, err = r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return domain.Task{}, err
	}
	task.ID = objID.Hex()
	return task, nil
}

func (r *MongoTaskRepository) DeleteTask(id string) error {
	ctx := context.Background()
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}
