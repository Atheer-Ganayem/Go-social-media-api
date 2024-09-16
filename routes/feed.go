package routes

import (
	"net/http"

	"github.com/Atheer-Ganayem/Go-social-media-api/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func getFeed(ctx *gin.Context) {
	userId, err := primitive.ObjectIDFromHex(ctx.GetString("userId"))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Not authenticated. Invalid auth token."})
		return
	}

	user, err := models.GetUser(bson.M{"_id": userId})
	if err == mongo.ErrNoDocuments {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Not authenticated. Invalid auth token."})
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "An error occurred. Please try again later."})
		return
	}

	posts, err := models.GetPosts(bson.M{"author": bson.M{"$in": user.Followings}})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "An error occurred. Please try again later."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Feed fetched.", "posts": posts})
}

func search(ctx *gin.Context) {
	query := ctx.Query("search")
	if query == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter is required!"})
		return
	}

	postFilter := bson.M{
		"$or": bson.A{
			bson.M{"content": bson.M{"$regex": query, "$options": "i"}},
			bson.M{"title": bson.M{"$regex": query, "$options": "i"}},
		},
	}

	posts, err := models.GetPosts(postFilter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "An error occurred. Please try again later."})
		return
	}

	userFilter := bson.M{
		"$or": bson.A{
			bson.M{"name": bson.M{"$regex": query, "$options": "i"}},
		},
	}

	users, err := models.GetUsers(userFilter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "An error occurred. Please try again later."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Search done successfully.", "posts": posts, "users": users})
}