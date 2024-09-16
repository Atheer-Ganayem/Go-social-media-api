package models

import (
	"context"
	"time"

	db "github.com/Atheer-Ganayem/Go-social-media-api/DB"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Comment struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Author    primitive.ObjectID `json:"author" bson:"author,omitempty"`
	Post    primitive.ObjectID `json:"post" bson:"post,omitempty"`
	Content   string             `json:"content" bson:"content,omitempty" binding:"required,min=1,max=100"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt,omitempty"`
}

func (comment Comment) Save() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := db.DB.Collection("comments").InsertOne(ctx, comment)

	return err
}

func DeletePostComments(postId primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := db.DB.Collection("comments").DeleteMany(ctx, bson.M{"post": postId})

	return err
}

func DeleteComment(commentId primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var comment Comment
	err := db.DB.Collection("comments").FindOneAndDelete(ctx, bson.M{"_id": commentId}).Decode(&comment)
	if err != nil {
		return err
	}

	_, err = db.DB.Collection("posts").UpdateByID(ctx, comment.Post, bson.M{"$pull": bson.M{"comments": commentId}})

	return err
}