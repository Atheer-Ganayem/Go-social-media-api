package routes

import (
	"net/http"
	"time"

	"github.com/Atheer-Ganayem/Go-social-media-api/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func createComment(ctx *gin.Context) {
	hexPostId, _ := ctx.Params.Get("postId")
	postId, err := primitive.ObjectIDFromHex(hexPostId)
	if err != nil {
		ctx.JSONP(http.StatusNotFound, gin.H{"message": "Post not found."})
		return
	}

	userId, err := primitive.ObjectIDFromHex(ctx.GetString("userId"))
	if err != nil {
		ctx.JSONP(http.StatusUnauthorized, gin.H{"message": "Not authenticated. Invalid auth token."})
		return
	}

	post, err := models.GetPost(bson.M{"_id": postId})
	if err == mongo.ErrNoDocuments {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Post not found."})
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "An error occurred. Please try again later."})
		return
	}

	var comment models.Comment
	err = ctx.ShouldBindJSON(&comment)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Comment's content is required. 100 chars at max."})
		return
	}

	comment.ID = primitive.NewObjectID()
	comment.Author = userId
	comment.Post = postId	
	comment.CreatedAt = time.Now()

	err = comment.Save()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "An error occurred, couldn't create comment. Please try again later."})
		return
	}

	err = post.AddComment(comment.ID)
		if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "An error occurred. comment created, but couldn't add to to the post. Please try again later."})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Comment has been created successfully.", "commentId": comment.ID})
}

func deleteComment(ctx *gin.Context) {
	hexCommentId, _ := ctx.Params.Get("id")
	commentId, err := primitive.ObjectIDFromHex(hexCommentId)
	if err != nil {
		ctx.JSONP(http.StatusNotFound, gin.H{"message": "Comment not found."})
		return
	}

	err = models.DeleteComment(commentId)
	if err == mongo.ErrNoDocuments {
		ctx.JSONP(http.StatusNotFound, gin.H{"message": "Comment not found."})
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "An error occurred. Please try again later."})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{"message": "Comment has been deleted successfully."})
}