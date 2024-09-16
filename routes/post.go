package routes

import (
	"log"
	"net/http"
	"time"

	"github.com/Atheer-Ganayem/Go-social-media-api/models"
	"github.com/Atheer-Ganayem/Go-social-media-api/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func createPost(ctx *gin.Context) {
	stringUserId := ctx.GetString("userId")
	userId, err := primitive.ObjectIDFromHex(stringUserId)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Not authenticated. Invalid auth token."})
		return
	}

	var post models.Post
	err = ctx.ShouldBind(&post)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Post content length can be 250 at max."})
		return
	}

	post.ID = primitive.NewObjectID()
	post.Author = userId
	post.CreatedAt = time.Now()
	post.Likes = make([]primitive.ObjectID, 0)
	post.Comments = make([]primitive.ObjectID, 0)

	filePath, _, err := utils.HandleImage(ctx.Request, "image")
	if err == nil {
		post.Image = filePath
	}

	err = post.Save()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Couldn't create post."})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Post has been created successfully", "postId": post.ID})
}

func deletePost(ctx *gin.Context) {
	stringPostId, _ := ctx.Params.Get("id")

	postId, err := primitive.ObjectIDFromHex(stringPostId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Post not found."})
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

	if post.Author.Hex() != ctx.GetString("userId") {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "You cannot delete a post that is not your's."})
		return
	}

	err = post.Delete(bson.M{"_id": postId})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "An error occurred. Please try again later."})
		return
	}

	if post.Image != "" {
		go func() {
			err := utils.RemoveFileFromS3(post.Image)
			if err != nil {
					log.Printf("Failed to remove file from S3: %s, error: %v", post.Image, err)
			}
    }()
	}

	go models.DeletePostComments(postId)

	ctx.JSON(http.StatusOK, gin.H{"message": "Post has been delete successfully."})
}

func likePost(ctx *gin.Context) {
	hexPostId, _ := ctx.Params.Get("id")
	postId, err := primitive.ObjectIDFromHex(hexPostId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Post not found."})
		return
	}

	userId, err := primitive.ObjectIDFromHex(ctx.GetString("userId"))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Not authenticated. Invalid auth token."})
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


	err = post.Like(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "An error occurred. Please try again later."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Post liked/unliked successfully."})
}