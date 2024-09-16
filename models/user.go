package models

import (
	"context"
	"time"

	db "github.com/Atheer-Ganayem/Go-social-media-api/DB"
	"github.com/Atheer-Ganayem/Go-social-media-api/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID primitive.ObjectID `json:"_id" bson:"_id"`
	Name string `json:"name" bson:"name" form:"name" binding:"required,min=3,max=30"`
	Email string `json:"email" bson:"email" form:"email" binding:"required,email"` 
	Avatar string `json:"avatar" bson:"avatar"`
	Password string `json:"password" bson:"password" form:"password" binding:"required,min=6"`
	Followers []primitive.ObjectID `json:"followers" bson:"followers"`
	Followings []primitive.ObjectID `json:"followings" bson:"followings"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
}

type LoginRequest struct {
	Email string `json:"email" bson:"email" binding:"required,email"` 
	Password string `json:"password" bson:"password" binding:"required"`
}

type UpdateRequest struct {
	Name string `json:"name" bson:"name" form:"name" binding:"omitempty,min=3,max=30"`
	Email string `json:"email" bson:"email" binding:"omitempty,email"`
	Password string `json:"password" bson:"password"`
	CurrentPassword string `json:"currentPassword" bson:"currentPassword"`
}

type SearchResults struct {
	ID primitive.ObjectID `json:"_id" bson:"_id"`
	Name string `json:"name" bson:"name" form:"name" binding:"required,min=3,max=30"`
	Avatar string `json:"avatar" bson:"avatar"`
}

func UpdateUser(filter, update bson.M) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := db.DB.Collection("users").UpdateOne(ctx, filter, update)

	return err
}

func (user User) Save() error {
	hashedPw, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}

	user.ID = primitive.NewObjectID()
	user.Password = hashedPw
	user.Followers = make([]primitive.ObjectID, 0);
	user.Followings = make([]primitive.ObjectID, 0);
	user.CreatedAt = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err = db.DB.Collection("users").InsertOne(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func CheckUserExists(email string) (bool) {
	var user User
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
	defer cancel()

	err := db.DB.Collection("users").FindOne(ctx, bson.M{"email": email}).Decode(&user)
	
	return err != mongo.ErrNoDocuments
}

func GetUser(filter bson.M) (User, error) {
	var user User
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
	defer cancel()

	err := db.DB.Collection("users").FindOne(ctx, filter).Decode(&user)
	
	return user, err
}

type Profile struct {
	ID primitive.ObjectID `json:"_id" bson:"_id"`
	Name string `json:"name" bson:"name"`
	Email string `json:"email" bson:"email"` 
	Avatar string `json:"avatar" bson:"avatar"`
	Followers []primitive.ObjectID `json:"followers" bson:"followers"`
	Followings []primitive.ObjectID `json:"followings" bson:"followings"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	Posts []Post `json:"posts" bson:"posts"`
}

func GetProfile(userId primitive.ObjectID) (Profile, error) {
	var profile Profile
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
	defer cancel()

	pipeline := mongo.Pipeline{
    bson.D{
        {Key: "$match", Value: bson.D{{Key: "_id", Value: userId}}},
    },
    bson.D{
        {Key: "$lookup", Value: bson.D{
            {Key: "from", Value: "posts"},
            {Key: "localField", Value: "_id"},
            {Key: "foreignField", Value: "author"},
            {Key: "as", Value: "posts"},
        }},
    },
	}

	cursor, err := db.DB.Collection("users").Aggregate(ctx, pipeline)
	if err != nil {return profile, err}

	if cursor.Next(ctx) {
		err := cursor.Decode(&profile)
		if err != nil {
			return profile, err
	}
	} else {
			return profile, err
	}

	return profile, err
}

func GetUsers(filter bson.M) ([]SearchResults, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
	defer cancel()

	cursor, err := db.DB.Collection("users").Find(ctx, filter)
	if err != nil {return nil, err}

	var result []SearchResults
	err = cursor.All(ctx, &result)

	return result, err
}