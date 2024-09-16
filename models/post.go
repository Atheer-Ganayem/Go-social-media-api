package models

import (
	"context"
	"time"

	db "github.com/Atheer-Ganayem/Go-social-media-api/DB"
	"github.com/Atheer-Ganayem/Go-social-media-api/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Post struct {
	ID        primitive.ObjectID   `json:"_id" bson:"_id"`
	Author  primitive.ObjectID   `json:"author" bson:"author"`
	Content   string               `json:"content" bson:"content" form:"content" binding:"max=250"`
	Image  string               `json:"image" bson:"image,omitempty"`
	Likes     []primitive.ObjectID `json:"likes" bson:"likes"`
	Comments []primitive.ObjectID `json:"comments" bson:"comments"`
	CreatedAt time.Time            `json:"createdAt" bson:"createdAt"`
}

func (post Post) Save() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := db.DB.Collection("posts").InsertOne(ctx, post)

	return err
}

func GetPost(filter bson.M) (Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var post Post
	err := db.DB.Collection("posts").FindOne(ctx, filter).Decode(&post)

	return post, err
}

func (post Post) Delete(filter bson.M) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := db.DB.Collection("posts").DeleteOne(ctx, filter)

	return err
}

func (post Post) Like(likerId primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var err error
	exists := utils.IncludesObjectId(post.Likes, likerId)
	if !exists {
		_, err = db.DB.Collection("posts").UpdateByID(ctx, post.ID, bson.M{"$addToSet": bson.M{"likes": likerId}})
	} else {
		_, err = db.DB.Collection("posts").UpdateByID(ctx, post.ID, bson.M{"$pull": bson.M{"likes": likerId}})
	}

	return err
}

func (post Post) AddComment(commentId primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := db.DB.Collection("posts").UpdateByID(ctx, post.ID, bson.M{"$addToSet": bson.M{"comments": commentId}})

	return err
}

func GetPosts(filter bson.M) ([]Post, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()


	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "createdAt", Value: -1}})
	cursor, err := db.DB.Collection("posts").Find(ctx, filter, findOptions)

	if err != nil {
		return nil, err
	}

	var posts []Post
	err = cursor.All(ctx, &posts)

	return posts, err
}