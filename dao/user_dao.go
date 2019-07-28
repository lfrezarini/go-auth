package dao

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/LucasFrezarini/go-auth-manager/models"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/bson"
)

// UserDao is a representation of a User DAO
type UserDao struct{}

// UserCollection defines the name of the collection that this DAO uses
const UserCollection = "users"

// GetAll fetch all the users registered on the database
func (d *UserDao) GetAll() []*models.User {
	collection := db.Collection(UserCollection)
	cursor, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		log.Panic(err)
	}

	defer cursor.Close(context.Background())

	var users []*models.User

	for cursor.Next(context.Background()) {
		user := models.User{}
		err := cursor.Decode(&user)

		if err != nil {
			log.Panic(err)
		}

		users = append(users, &user)
	}

	if err := cursor.Err(); err != nil {
		log.Panic(err)
	}

	return users
}

// CreateOne create an user in the collection on the database
func (d *UserDao) CreateOne(user models.User) (primitive.ObjectID, error) {
	collection := db.Collection(UserCollection)
	bson, err := bson.Marshal(user)

	if err != nil {
		log.Panic(err)
		return primitive.NilObjectID, errors.New("Error while trying to convert the input to BSON")
	}

	res, err := collection.InsertOne(context.Background(), bson)
	if err != nil {
		log.Panic(err)
		return primitive.NilObjectID, errors.New("Error while trying to insert the data into the collection User")
	}

	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		return oid, nil
	}

	log.Panic("Error while trying to parse the InsertedID")
	return primitive.NilObjectID, errors.New("Error while trying to parse the InsertedID")
}

// FindOne returns a result based on the fields passed on the user struct
func (d *UserDao) FindOne(user models.User) (*models.User, error) {
	collection := db.Collection(UserCollection)
	bson, err := bson.Marshal(user)

	if err != nil {
		log.Panic(err)
		return nil, errors.New("Error while trying to convert the input to BSON")
	}

	result := models.User{}
	err = collection.FindOne(context.Background(), bson).Decode(&result)

	if err != nil {
		return nil, fmt.Errorf("Error while trying to fetch user from the database: %v", err)
	}

	return &result, nil
}

// FindByID returns the user from the database with the respective _id
func (d *UserDao) FindByID(id primitive.ObjectID) (*models.User, error) {
	collection := db.Collection(UserCollection)

	result := models.User{}
	err := collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&result)

	if err != nil {
		return nil, fmt.Errorf("Error while trying to fetch user with id %s from the database: %v", id.String(), err)
	}

	return &result, nil
}
