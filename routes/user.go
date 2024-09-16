package routes

import (
	"log"
	"net/http"

	"github.com/Atheer-Ganayem/Go-social-media-api/models"
	"github.com/Atheer-Ganayem/Go-social-media-api/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// follow/unfollow
func follow(ctx *gin.Context) {
	targetStringUserId, _ := ctx.Params.Get("id")
	clientStringUserId := ctx.GetString("userId")
	if targetStringUserId == clientStringUserId {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "You cannot follow yourself."})
		return 
	}
	
	targetUserId, err := primitive.ObjectIDFromHex(targetStringUserId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "User not found."})
		return 
	}
	
	clientUserId, err := primitive.ObjectIDFromHex(clientStringUserId)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Not authenticated. Invalid auth token."})
		return 
	}

	targetUser, err1 := models.GetUser(bson.M{"_id": targetUserId})
	clientUser, err2 := models.GetUser(bson.M{"_id": clientUserId})
	if err1 != nil || err2 != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "An error occurred. Please try again later."})
		return
	}

	exists := utils.IncludesObjectId(targetUser.Followers, clientUser.ID)

	if !exists {
		err := models.UpdateUser(bson.M{"_id": targetUserId}, bson.M{"$addToSet": bson.M{"followers": clientUserId}})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Couldn't follow user"})
			return
		}
		err = models.UpdateUser(bson.M{"_id": clientUserId}, bson.M{"$addToSet": bson.M{"followings": targetUserId}})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Couldn't follow user"})
			return
		}
	} else {
		err := models.UpdateUser(bson.M{"_id": targetUserId}, bson.M{"$pull": bson.M{"followers": clientUserId}})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Couldn't follow user"})
			return
		}
		err = models.UpdateUser(bson.M{"_id": clientUserId}, bson.M{"$pull": bson.M{"followings": targetUserId}})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Couldn't follow user"})
			return
		}
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "User followed/unfollowed successfully."})
}

func getProfile(ctx *gin.Context) {
	hexUserId, _ := ctx.Params.Get("id")
	userId, err := primitive.ObjectIDFromHex(hexUserId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "User not found."})
		return 
	}

	profile, err := models.GetProfile(userId)
	if err == mongo.ErrNoDocuments {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "User not found."})
		return 
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "An error occurred. Please try again later."})
		return 
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User fetched successfully.","profile": profile})
}

func updateProfile(ctx *gin.Context) {
	userId, err := primitive.ObjectIDFromHex(ctx.GetString("userId"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "User not found."})
		return 
	}

	user, err := models.GetUser(bson.M{"_id": userId})
	if err == mongo.ErrNoDocuments {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Post not found."})
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "An error occurred. Please try again later."})
		return
	}

	var data models.UpdateRequest
	err = ctx.ShouldBind(&data)
	if err != nil {
		println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid input."})
		return
	}

	if data.Name != "" {
		user.Name = data.Name
	}
	if data.Email != "" {
		exists := models.CheckUserExists(data.Email)
		if exists {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Email already in use. Pleae choose another email."})
			return
		}
		user.Email = data.Email
	}
	if data.Password != "" {
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.CurrentPassword))
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Incorrect current password."})
			return
		}

		hashedPw, err := utils.HashPassword(data.Password)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "An error occurred. Please try again later."})
			return
		}

		user.Password = hashedPw
	}

	_, _, err = ctx.Request.FormFile("avatar")
	if err == nil {
		avatarPath, code, err := utils.HandleImage(ctx.Request, "avatar")
		if err != nil {
			ctx.JSON(code, gin.H{"message": err.Error()})
			return
		}
		oldAvatar := user.Avatar
		go func() {
			err := utils.RemoveFileFromS3(oldAvatar)
			if err != nil {
					log.Printf("Failed to remove file from S3: %s, error: %v", user.Avatar, err)
			}
    }()
		user.Avatar = avatarPath
	}

	err = models.UpdateUser(bson.M{"_id": user.ID}, bson.M{"$set": bson.M{"name": user.Name, "email": user.Email,
		"password": user.Password, "avatar": user.Avatar }})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "An error occurred. Please try again later."})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Details has been updated successfully."})
}